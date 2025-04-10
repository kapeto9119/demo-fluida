package solana

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/models"
	"github.com/ncapetillo/demo-fluida/internal/repository"
	"gorm.io/gorm"
)

const (
	// USDC Token mint address on Solana devnet
	// This is the official USDC mint address for Solana devnet
	USDCDevnetMint = "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	
	// Polling interval for checking transactions
	pollInterval = 15 * time.Second
)

// PaymentWatcher monitors Solana blockchain for USDC payments to specific addresses
type PaymentWatcher struct {
	rpcClient  *rpc.Client
	repository repository.InvoiceRepository
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewPaymentWatcher creates a new payment watcher for the Solana blockchain
func NewPaymentWatcher() (*PaymentWatcher, error) {
	// We're using devnet for development
	endpoint := rpc.DevNet_RPC
	
	// Create RPC client
	rpcClient := rpc.New(endpoint)
	
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	// Initialize repository
	invoiceRepo := repository.NewInvoiceRepository(db.DB)
	
	return &PaymentWatcher{
		rpcClient:  rpcClient,
		repository: invoiceRepo,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Start begins the payment watching process using polling approach
func (pw *PaymentWatcher) Start() {
	log.Println("Starting Solana payment watcher")
	
	// Run the polling in a goroutine
	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-pw.ctx.Done():
				log.Println("Payment watcher shutting down")
				return
			case <-ticker.C:
				if err := pw.checkPendingInvoices(); err != nil {
					log.Printf("Error checking pending invoices: %v", err)
				}
			}
		}
	}()
}

// Stop halts the payment watcher
func (pw *PaymentWatcher) Stop() {
	pw.cancel()
	log.Println("Payment watcher stopped")
}

// WatchForPayments starts watching for payments (alias for Start method)
func (pw *PaymentWatcher) WatchForPayments() {
	pw.Start()
}

// checkPendingInvoices looks for pending invoices and checks for payments
func (pw *PaymentWatcher) checkPendingInvoices() error {
	// Create a timeout context for this operation
	ctx, cancel := context.WithTimeout(pw.ctx, 30*time.Second)
	defer cancel()
	
	// Find all pending invoices
	pendingInvoices, err := pw.repository.FindPendingInvoices(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch pending invoices: %v", err)
	}
	
	// Check each invoice for payments
	for _, invoice := range pendingInvoices {
		// Skip checking if context is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue processing
		}
		
		paid, err := pw.checkForPayment(ctx, invoice)
		if err != nil {
			log.Printf("Error checking payment for invoice %s: %v", invoice.InvoiceNumber, err)
			continue
		}
		
		if paid {
			// Use database transaction to update the invoice status
			err := db.DB.Transaction(func(tx *gorm.DB) error {
				txCtx, txCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer txCancel()
				
				txRepo := repository.NewInvoiceRepository(tx)
				// Update invoice status to PAID
				return txRepo.UpdateStatus(txCtx, invoice.ID, models.StatusPaid)
			})
			
			if err != nil {
				log.Printf("Failed to update invoice %s to PAID: %v", invoice.InvoiceNumber, err)
			} else {
				log.Printf("Invoice %s marked as PAID", invoice.InvoiceNumber)
			}
		}
	}
	
	return nil
}

// checkForPayment checks if a specific invoice has been paid
func (pw *PaymentWatcher) checkForPayment(ctx context.Context, invoice models.Invoice) (bool, error) {
	// Parse receiver address
	receiverPubkey, err := solana.PublicKeyFromBase58(invoice.ReceiverAddr)
	if err != nil {
		return false, fmt.Errorf("invalid receiver address: %v", err)
	}
	
	// Get recent signatures for the account
	signatures, err := pw.rpcClient.GetSignaturesForAddress(ctx, receiverPubkey)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction signatures: %v", err)
	}
	
	// Check each transaction
	for _, sig := range signatures {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			// Continue processing
		}
		
		// Skip failed transactions
		if sig.Err != nil {
			continue
		}
		
		// Get transaction details
		txSig, err := solana.SignatureFromBase58(sig.Signature.String())
		if err != nil {
			log.Printf("Invalid signature format: %v", err)
			continue
		}
		
		tx, err := pw.rpcClient.GetTransaction(
			ctx, 
			txSig,
			&rpc.GetTransactionOpts{
				Encoding: solana.EncodingJSONParsed,
			},
		)
		if err != nil {
			log.Printf("Failed to get transaction details: %v", err)
			continue
		}
		
		// Check if this is a USDC payment to the receiver
		if isPaymentForInvoice(tx, invoice, receiverPubkey) {
			return true, nil
		}
	}
	
	return false, nil
}

// isPaymentForInvoice determines if a transaction represents payment for an invoice
func isPaymentForInvoice(tx *rpc.GetTransactionResult, invoice models.Invoice, receiverPubkey solana.PublicKey) bool {
	// Ensure we have transaction data
	if tx == nil || tx.Meta == nil {
		return false
	}
	
	// USDC uses 6 decimals
	const usdcDecimals = 6
	
	// Check for token transfers in the transaction
	for _, postBalance := range tx.Meta.PostTokenBalances {
		// Check if this is for the USDC token mint
		// Use string comparison instead of direct type comparison
		if postBalance.Mint.String() != USDCDevnetMint {
			continue
		}
		
		// Check if the owner is our receiver
		if postBalance.Owner == nil {
			continue
		}
		
		if postBalance.Owner.String() != receiverPubkey.String() {
			continue
		}
		
		// Find the pre-balance for comparison
		var preBalance *rpc.TokenBalance
		for _, pre := range tx.Meta.PreTokenBalances {
			if pre.AccountIndex == postBalance.AccountIndex {
				preBalance = &pre
				break
			}
		}
		
		if preBalance == nil {
			continue
		}
		
		// Calculate the amount received (post - pre)
		preAmount := parseUiAmount(preBalance.UiTokenAmount.Amount)
		postAmount := parseUiAmount(postBalance.UiTokenAmount.Amount)
		
		if preAmount == nil || postAmount == nil {
			continue
		}
		
		diff := new(big.Float).Sub(postAmount, preAmount)
		
		// Convert to float64 for comparison
		amountReceived, _ := diff.Float64()
		
		// Allow a small tolerance for floating point comparison
		const tolerance = 0.001
		if amountReceived >= invoice.Amount-tolerance && amountReceived <= invoice.Amount+tolerance {
			log.Printf("Found matching USDC payment: expected %f, received %f", invoice.Amount, amountReceived)
			return true
		}
	}
	
	return false
}

// parseUiAmount converts a UI amount string to a big.Float
func parseUiAmount(amount string) *big.Float {
	if amount == "" {
		return nil
	}
	
	result, success := new(big.Float).SetString(amount)
	if !success {
		return nil
	}
	
	return result
} 