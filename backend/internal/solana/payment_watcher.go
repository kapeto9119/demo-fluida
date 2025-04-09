package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/ncapetillo/demo-fluida/internal/db"
	"github.com/ncapetillo/demo-fluida/internal/models"
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
	rpcClient *rpc.Client
	wsClient  *ws.Client
	database  *gorm.DB
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPaymentWatcher creates a new payment watcher for the Solana blockchain
func NewPaymentWatcher() (*PaymentWatcher, error) {
	// We're using devnet for development
	endpoint := rpc.DevNet_RPC
	wsEndpoint := rpc.DevNet_WS
	
	// Create RPC client
	rpcClient := rpc.New(endpoint)
	
	// Create WebSocket client
	wsClient, err := ws.Connect(context.Background(), wsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Solana WebSocket: %v", err)
	}
	
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	return &PaymentWatcher{
		rpcClient: rpcClient,
		wsClient:  wsClient,
		database:  db.DB,
		ctx:       ctx,
		cancel:    cancel,
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
	if pw.wsClient != nil {
		pw.wsClient.Close()
	}
	log.Println("Payment watcher stopped")
}

// checkPendingInvoices looks for pending invoices and checks for payments
func (pw *PaymentWatcher) checkPendingInvoices() error {
	var pendingInvoices []models.Invoice
	
	// Find all pending invoices
	if err := pw.database.Where("status = ?", models.StatusPending).Find(&pendingInvoices).Error; err != nil {
		return fmt.Errorf("failed to fetch pending invoices: %v", err)
	}
	
	// Check each invoice for payments
	for _, invoice := range pendingInvoices {
		paid, err := pw.checkForPayment(invoice)
		if err != nil {
			log.Printf("Error checking payment for invoice %s: %v", invoice.InvoiceNumber, err)
			continue
		}
		
		if paid {
			// Update invoice status to PAID
			if err := models.UpdateInvoiceStatus(pw.database, invoice.ID, models.StatusPaid); err != nil {
				log.Printf("Failed to update invoice %s to PAID: %v", invoice.InvoiceNumber, err)
			} else {
				log.Printf("Invoice %s marked as PAID", invoice.InvoiceNumber)
			}
		}
	}
	
	return nil
}

// checkForPayment checks if a specific invoice has been paid
func (pw *PaymentWatcher) checkForPayment(invoice models.Invoice) (bool, error) {
	// Parse receiver address
	receiverPubkey, err := solana.PublicKeyFromBase58(invoice.ReceiverAddr)
	if err != nil {
		return false, fmt.Errorf("invalid receiver address: %v", err)
	}
	
	// Parse USDC token mint
	usdcMint, err := solana.PublicKeyFromBase58(USDCDevnetMint)
	if err != nil {
		return false, fmt.Errorf("invalid USDC mint address: %v", err)
	}
	
	// Get recent transactions for the receiver address
	txSignatures, err := pw.rpcClient.GetSignaturesForAddress(
		context.Background(),
		receiverPubkey,
		&rpc.GetSignaturesForAddressOpts{
			Limit: 10, // Limit to recent transactions
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction signatures: %v", err)
	}
	
	// Check each transaction
	for _, sigInfo := range txSignatures {
		// Only check confirmed transactions
		if sigInfo.Err != nil {
			continue
		}
		
		// Parse transaction signature
		signature, err := solana.SignatureFromBase58(sigInfo.Signature)
		if err != nil {
			log.Printf("Invalid signature format: %v", err)
			continue
		}
		
		// Get transaction details
		tx, err := pw.rpcClient.GetTransaction(
			context.Background(),
			signature,
			&rpc.GetTransactionOpts{
				Encoding: solana.EncodingJSON,
			},
		)
		if err != nil {
			log.Printf("Failed to get transaction details: %v", err)
			continue
		}
		
		// Check if this is a token transfer
		if isUSDCPayment(tx, usdcMint, receiverPubkey, invoice.Amount) {
			return true, nil
		}
	}
	
	return false, nil
}

// isUSDCPayment checks if a transaction is a valid USDC payment for the invoice
func isUSDCPayment(tx *rpc.GetTransactionResult, usdcMint, receiver solana.PublicKey, expectedAmount float64) bool {
	// Log the transaction we're checking
	log.Printf("Checking transaction for USDC payment: %s", tx.Transaction.Signatures[0])
	
	// Ensure we have transaction data
	if tx == nil || tx.Meta == nil {
		return false
	}
	
	// USDC uses 6 decimals
	const usdcDecimals = 6
	
	// Get post token balances
	for _, balance := range tx.Meta.PostTokenBalances {
		// Check if this is for the USDC token
		mintAddress := balance.Mint
		if mintAddress != usdcMint.String() {
			continue
		}
		
		// Check if the destination is our receiver
		ownerAddress := balance.Owner
		if ownerAddress != receiver.String() {
			continue
		}
		
		// Find the pre-balance for the same account
		var preBalance *rpc.TokenBalance
		for _, pb := range tx.Meta.PreTokenBalances {
			if pb.AccountIndex == balance.AccountIndex && pb.Mint == balance.Mint {
				preBalance = &pb
				break
			}
		}
		
		// If we found a pre-balance, calculate the difference
		if preBalance != nil {
			// Get balances as uint64
			preAmount, ok := new(big.Int).SetString(preBalance.UITokenAmount.Amount, 10)
			if !ok {
				log.Printf("Failed to parse pre-balance amount: %s", preBalance.UITokenAmount.Amount)
				continue
			}
			
			postAmount, ok := new(big.Int).SetString(balance.UITokenAmount.Amount, 10)
			if !ok {
				log.Printf("Failed to parse post-balance amount: %s", balance.UITokenAmount.Amount)
				continue
			}
			
			// Calculate difference (postAmount - preAmount)
			diff := new(big.Int).Sub(postAmount, preAmount)
			
			// Convert to float with 6 decimals (USDC)
			amountReceived := convertTokenAmount(diff, usdcDecimals)
			
			// Check if the transferred amount matches the expected amount
			// We allow a small difference to account for potential calculation errors
			const tolerance = 0.001 // small tolerance for floating point comparison
			if amountReceived >= expectedAmount-tolerance && amountReceived <= expectedAmount+tolerance {
				log.Printf("Found matching USDC payment: expected %f, received %f", expectedAmount, amountReceived)
				return true
			}
			
			log.Printf("Found USDC transfer but amount doesn't match: expected %f, received %f", 
				expectedAmount, amountReceived)
		}
	}
	
	return false
}

// convertTokenAmount converts a token amount based on decimals
func convertTokenAmount(amount *big.Int, decimals uint8) float64 {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	quotient := new(big.Float).Quo(
		new(big.Float).SetInt(amount),
		new(big.Float).SetInt(divisor),
	)
	
	result, _ := quotient.Float64()
	return result
} 