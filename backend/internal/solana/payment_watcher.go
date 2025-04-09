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
	// Note: This is just an example - you'll need to replace this with the actual devnet USDC mint address
	USDCDevnetMint = "Gh9ZwEmdLJ8DscKNTkTqPbNwLNNBjuSzaG9Vp2KGtKJr"
	
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
	// This is a simplified example. In a real implementation, you would:
	// 1. Parse the transaction data to find SPL token transfers
	// 2. Check if the token is USDC (matches usdcMint)
	// 3. Verify the destination is the receiver
	// 4. Confirm the amount matches expectedAmount
	
	// For this prototype, we'll use a mock implementation that just logs
	// In a real implementation, you'd parse the transaction instruction data
	
	log.Printf("Checking transaction for USDC payment: %s", tx.Transaction.Signatures[0])
	
	// Simplistic approach - in real code you'd properly parse the transaction
	// and its token transfer instructions
	if tx != nil && tx.Meta != nil {
		// For demonstration purposes only - this won't actually work as is
		// You would need to properly decode the transaction instructions
		
		// Let's pretend we found a match - would be implemented with proper
		// transaction instruction parsing in a real app
		
		// Just for logging - not real detection logic
		jsonData, _ := json.MarshalIndent(tx, "", "  ")
		log.Printf("Transaction data: %s", string(jsonData))
		
		// In a real implementation, return true ONLY if:
		// 1. The transaction contains an SPL token transfer
		// 2. The token is USDC (matches usdcMint)
		// 3. The destination address matches the receiver
		// 4. The amount (after adjusting for decimals) matches expectedAmount
		
		// For now, always return false since we're not actually parsing
		return false
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