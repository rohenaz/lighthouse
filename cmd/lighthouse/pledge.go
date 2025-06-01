package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/template/p2pkh"
	"github.com/spf13/cobra"
	"github.com/yourusername/lighthouse/core"
)

// pledgeCreateCmd creates a new pledge
func pledgeCreateCmd() *cobra.Command {
	var (
		amount    float64
		message   string
		name      string
		email     string
		refund    string
		wif       string
		utxos     []string
		output    string
	)

	cmd := &cobra.Command{
		Use:   "create [project-file]",
		Short: "Create a pledge to a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectFile := args[0]
			
			// Read the project file
			data, err := ioutil.ReadFile(projectFile)
			if err != nil {
				return fmt.Errorf("failed to read project file: %w", err)
			}
			
			// Load the project
			project, err := core.LoadProject(data)
			if err != nil {
				return fmt.Errorf("failed to load project: %w", err)
			}
			
			// Convert BSV to satoshis
			amountSatoshis := uint64(amount * 100000000)
			
			// Parse WIF private key
			if wif == "" {
				return fmt.Errorf("private key (--wif) is required")
			}
			
			privKey, err := ec.PrivateKeyFromWif(wif)
			if err != nil {
				return fmt.Errorf("invalid WIF private key: %w", err)
			}
			
			// Parse UTXOs
			if len(utxos) == 0 {
				return fmt.Errorf("at least one UTXO is required (--utxo)")
			}
			
			var txUTXOs []*transaction.UTXO
			for _, utxoStr := range utxos {
				// Expected format: txid:vout:satoshis
				parts := strings.Split(utxoStr, ":")
				if len(parts) != 3 {
					return fmt.Errorf("invalid UTXO format: %s (expected txid:vout:satoshis)", utxoStr)
				}
				
				txid := parts[0]
				vout := 0
				if _, err := fmt.Sscanf(parts[1], "%d", &vout); err != nil {
					return fmt.Errorf("invalid vout in UTXO: %s", parts[1])
				}
				satoshis := uint64(0)
				if _, err := fmt.Sscanf(parts[2], "%d", &satoshis); err != nil {
					return fmt.Errorf("invalid satoshis in UTXO: %s", parts[2])
				}
				
				// Get the locking script for our address
				pubKey := privKey.PubKey()
				address, err := script.NewAddressFromPublicKey(pubKey, true) // mainnet
				if err != nil {
					return fmt.Errorf("failed to create address: %w", err)
				}
				lockingScriptHex := createP2PKHLockingScriptHex(address.AddressString)
				
				utxo, err := transaction.NewUTXO(txid, uint32(vout), lockingScriptHex, satoshis)
				if err != nil {
					return fmt.Errorf("failed to create UTXO: %w", err)
				}
				
				txUTXOs = append(txUTXOs, utxo)
			}
			
			// Create the pledge
			pledge, err := core.NewPledge(project, amountSatoshis, txUTXOs)
			if err != nil {
				return fmt.Errorf("failed to create pledge: %w", err)
			}
			
			// Set optional fields
			if message != "" {
				pledge.SetMemo(message)
			}
			if refund != "" {
				pledge.SetRefundAddress(refund)
			}
			if name != "" || email != "" {
				pledge.SetContactInfo(name, email)
			}
			
			// Sign the pledge
			if err := pledge.Sign([]*ec.PrivateKey{privKey}); err != nil {
				return fmt.Errorf("failed to sign pledge: %w", err)
			}
			
			// Serialize the pledge
			pledgeData, err := pledge.Serialize()
			if err != nil {
				return fmt.Errorf("failed to serialize pledge: %w", err)
			}
			
			// Determine output filename
			if output == "" {
				baseName := strings.TrimSuffix(filepath.Base(projectFile), filepath.Ext(projectFile))
				output = fmt.Sprintf("%s-%s.pledge", baseName, pledge.ID()[:8])
			}
			
			// Write to file
			if err := ioutil.WriteFile(output, pledgeData, 0644); err != nil {
				return fmt.Errorf("failed to write pledge file: %w", err)
			}
			
			fmt.Printf("Pledge created successfully!\n")
			fmt.Printf("File: %s\n", output)
			fmt.Printf("ID: %s\n", pledge.ID())
			fmt.Printf("Amount: %.8f BSV (%d satoshis)\n", amount, amountSatoshis)
			fmt.Printf("Project: %s\n", project.Title())
			
			return nil
		},
	}

	cmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Pledge amount in BSV (required)")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Optional message to project creator")
	cmd.Flags().StringVar(&name, "name", "", "Your name (optional)")
	cmd.Flags().StringVar(&email, "email", "", "Your email (optional)")
	cmd.Flags().StringVar(&refund, "refund", "", "Refund address if project fails")
	cmd.Flags().StringVarP(&wif, "wif", "w", "", "Private key in WIF format (required)")
	cmd.Flags().StringSliceVarP(&utxos, "utxo", "u", []string{}, "UTXOs to use (format: txid:vout:satoshis)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output filename")

	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("wif")
	cmd.MarkFlagRequired("utxo")

	return cmd
}

// pledgeViewCmd displays pledge details
func pledgeViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [pledge-file]",
		Short: "View pledge details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pledgeFile := args[0]
			
			// Read the pledge file
			data, err := ioutil.ReadFile(pledgeFile)
			if err != nil {
				return fmt.Errorf("failed to read pledge file: %w", err)
			}
			
			// Load the pledge
			pledge, err := core.LoadPledge(data)
			if err != nil {
				return fmt.Errorf("failed to load pledge: %w", err)
			}
			
			// Display pledge details
			fmt.Printf("Pledge ID: %s\n", pledge.ID())
			fmt.Printf("Project ID: %s\n", pledge.ProjectID())
			fmt.Printf("Amount: %.8f BSV (%d satoshis)\n", 
				float64(pledge.Amount())/100000000, pledge.Amount())
			
			// Display transaction details
			if tx := pledge.Transaction(); tx != nil {
				fmt.Printf("Transaction: %s\n", tx.TxID())
				fmt.Printf("Inputs: %d\n", len(tx.Inputs))
				for i, input := range tx.Inputs {
					fmt.Printf("  Input %d: %s:%d\n", i, 
						hex.EncodeToString(input.SourceTXID[:]), 
						input.SourceTxOutIndex)
				}
			}
			
			return nil
		},
	}
}

// pledgeRevokeCmd revokes a pledge
func pledgeRevokeCmd() *cobra.Command {
	var (
		broadcast bool
		wif       string
		output    string
	)

	cmd := &cobra.Command{
		Use:   "revoke [pledge-file]",
		Short: "Revoke a pledge and reclaim funds",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pledgeFile := args[0]
			
			// Read the pledge file
			data, err := ioutil.ReadFile(pledgeFile)
			if err != nil {
				return fmt.Errorf("failed to read pledge file: %w", err)
			}
			
			// Load the pledge
			pledge, err := core.LoadPledge(data)
			if err != nil {
				return fmt.Errorf("failed to load pledge: %w", err)
			}
			
			// Parse WIF private key
			if wif == "" {
				return fmt.Errorf("private key (--wif) is required")
			}
			
			privKey, err := ec.PrivateKeyFromWif(wif)
			if err != nil {
				return fmt.Errorf("invalid WIF private key: %w", err)
			}
			
			// Create revocation transaction
			// This spends the pledge inputs back to the original address
			revokeTx := transaction.NewTransaction()
			
			// Add inputs from the pledge
			pledgeTx := pledge.Transaction()
			if pledgeTx == nil {
				return fmt.Errorf("pledge has no transaction")
			}
			
			for _, input := range pledgeTx.Inputs {
				revokeTx.AddInput(input)
			}
			
			// Add output back to our address
			pubKey := privKey.PubKey()
			address, err := script.NewAddressFromPublicKey(pubKey, true) // mainnet
			if err != nil {
				return fmt.Errorf("failed to create address: %w", err)
			}
			addr, err := script.NewAddressFromString(address.AddressString)
			if err != nil {
				return fmt.Errorf("failed to create address: %w", err)
			}
			
			// Create P2PKH locking script
			lockingScript, err := p2pkh.Lock(addr)
			if err != nil {
				return fmt.Errorf("failed to create locking script: %w", err)
			}
			
			// Calculate total input (would need to look up UTXOs in real implementation)
			totalAmount := pledge.Amount()
			fee := uint64(1000) // Simple fixed fee for now
			
			if totalAmount > fee {
				revokeTx.AddOutput(&transaction.TransactionOutput{
					Satoshis:      totalAmount - fee,
					LockingScript: lockingScript,
				})
			}
			
			// Sign the revocation transaction
			// In a real implementation, we'd need to sign properly
			fmt.Printf("Revocation transaction created (signing not implemented)\n")
			
			// Save the transaction
			txHex := revokeTx.String()
			if output == "" {
				output = fmt.Sprintf("%s-revoke.tx", pledgeFile)
			}
			
			if err := ioutil.WriteFile(output, []byte(txHex), 0644); err != nil {
				return fmt.Errorf("failed to write transaction: %w", err)
			}
			
			fmt.Printf("Revocation transaction created!\n")
			fmt.Printf("File: %s\n", output)
			fmt.Printf("Note: Transaction signing not yet implemented\n")
			
			if broadcast {
				fmt.Printf("\nBroadcasting not yet implemented\n")
			}
			
			return nil
		},
	}

	cmd.Flags().BoolVarP(&broadcast, "broadcast", "b", false, "Broadcast the revocation transaction")
	cmd.Flags().StringVarP(&wif, "wif", "w", "", "Private key in WIF format (required)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output transaction file")

	cmd.MarkFlagRequired("wif")

	return cmd
}

// createP2PKHLockingScriptHex creates a P2PKH locking script for an address
func createP2PKHLockingScriptHex(address string) string {
	// This is a simplified version - in production, use proper script building
	// P2PKH script: OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
	
	// For now, return a dummy script
	// In real implementation, extract the pubkey hash from the address
	return "76a914" + strings.Repeat("00", 20) + "88ac"
}