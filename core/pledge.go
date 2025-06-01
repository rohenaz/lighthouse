package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/bsv-blockchain/go-sdk/chainhash"
	ec "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	sighash "github.com/bsv-blockchain/go-sdk/transaction/sighash"
	"github.com/bsv-blockchain/go-sdk/transaction/template/p2pkh"
	pb "github.com/yourusername/lighthouse/core/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Pledge represents a contribution to a project
type Pledge struct {
	pb        *pb.Pledge
	id        string
	amount    uint64
	tx        *transaction.Transaction
}

// NewPledge creates a new pledge for a project
func NewPledge(project *Project, amount uint64, utxos []*transaction.UTXO) (*Pledge, error) {
	if amount < project.MinPledgeAmount() {
		return nil, fmt.Errorf("pledge amount %d is less than minimum %d", amount, project.MinPledgeAmount())
	}

	// Create a transaction with SIGHASH_ANYONECANPAY inputs
	tx := transaction.NewTransaction()
	
	// Add inputs from UTXOs
	totalInput := uint64(0)
	if err := tx.AddInputsFromUTXOs(utxos...); err != nil {
		return nil, fmt.Errorf("failed to add inputs: %w", err)
	}
	
	for _, utxo := range utxos {
		totalInput += utxo.Satoshis
		if totalInput >= amount {
			break
		}
	}

	if totalInput < amount {
		return nil, fmt.Errorf("insufficient funds: have %d, need %d", totalInput, amount)
	}

	// Add project outputs
	outputs, err := project.Outputs()
	if err != nil {
		return nil, fmt.Errorf("failed to get project outputs: %w", err)
	}

	// For a pledge, we create outputs proportional to the pledge amount
	// relative to the project goal
	pledgeRatio := float64(amount) / float64(project.GoalAmount())
	
	for _, out := range outputs {
		pledgeOutput := &transaction.TransactionOutput{
			Satoshis:      uint64(float64(out.Satoshis) * pledgeRatio),
			LockingScript: out.LockingScript,
		}
		tx.AddOutput(pledgeOutput)
	}

	// Create the pledge protobuf
	pledge := &pb.Pledge{
		ProjectId: []byte(project.ID()),
		Time:      timestamppb.Now(),
	}

	// Store input information
	for _, input := range tx.Inputs {
		pbInput := &pb.Input{
			TxHash:      input.SourceTXID[:],
			OutputIndex: input.SourceTxOutIndex,
			Sequence:    input.SequenceNumber,
		}
		
		// We'll add the unlock script after signing
		pledge.Inputs = append(pledge.Inputs, pbInput)
	}

	p := &Pledge{
		pb:     pledge,
		amount: amount,
		tx:     tx,
	}
	p.id = p.calculateID()

	return p, nil
}

// Sign signs the pledge with SIGHASH_ANYONECANPAY flag
func (p *Pledge) Sign(privateKeys []*ec.PrivateKey) error {
	if p.tx == nil {
		return errors.New("no transaction to sign")
	}

	// Sign each input with SIGHASH_ANYONECANPAY
	for i := range p.tx.Inputs {
		if i >= len(privateKeys) {
			return fmt.Errorf("no private key for input %d", i)
		}
		
		// Create P2PKH unlocker with ANYONECANPAY flag
		anyoneCanPayFlag := sighash.AllForkID | sighash.AnyOneCanPay
		unlocker, err := p2pkh.Unlock(privateKeys[i], &anyoneCanPayFlag)
		if err != nil {
			return fmt.Errorf("failed to create unlocker for input %d: %w", i, err)
		}

		// Sign the input
		unlockingScript, err := unlocker.Sign(p.tx, uint32(i))
		if err != nil {
			return fmt.Errorf("failed to sign input %d: %w", i, err)
		}
		
		p.tx.Inputs[i].UnlockingScript = unlockingScript
		
		// Update the protobuf with the unlock script
		p.pb.Inputs[i].UnlockScript = unlockingScript.Bytes()
	}

	return nil
}

// LoadPledge loads a pledge from serialized data
func LoadPledge(data []byte) (*Pledge, error) {
	var pledge pb.Pledge
	if err := proto.Unmarshal(data, &pledge); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pledge: %w", err)
	}

	// Reconstruct the transaction from the pledge data
	tx := transaction.NewTransaction()
	amount := uint64(0)

	// Add inputs
	for _, input := range pledge.Inputs {
		// Create a chainhash from the transaction ID bytes
		txid, err := chainhash.NewHash(input.TxHash)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction hash: %w", err)
		}
		
		unlockScript := script.Script(input.UnlockScript)
		txInput := &transaction.TransactionInput{
			SourceTXID:       txid,
			SourceTxOutIndex: input.OutputIndex,
			UnlockingScript:  &unlockScript,
			SequenceNumber:   input.Sequence,
		}
		tx.Inputs = append(tx.Inputs, txInput)
	}

	p := &Pledge{
		pb:     &pledge,
		amount: amount,
		tx:     tx,
	}
	p.id = p.calculateID()

	return p, nil
}

// Serialize returns the pledge as protobuf bytes
func (p *Pledge) Serialize() ([]byte, error) {
	return proto.Marshal(p.pb)
}

// ID returns the unique pledge ID
func (p *Pledge) ID() string {
	return p.id
}

// calculateID generates a unique ID from pledge data
func (p *Pledge) calculateID() string {
	data, _ := p.Serialize()
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Amount returns the pledged amount in satoshis
func (p *Pledge) Amount() uint64 {
	return p.amount
}

// ProjectID returns the ID of the project this pledge is for
func (p *Pledge) ProjectID() string {
	return string(p.pb.ProjectId)
}

// SetMemo sets a message from the pledger
func (p *Pledge) SetMemo(memo string) {
	p.pb.Memo = memo
	p.id = p.calculateID()
}

// SetRefundAddress sets where to refund if project fails
func (p *Pledge) SetRefundAddress(address string) {
	p.pb.RefundAddress = address
	p.id = p.calculateID()
}

// SetContactInfo sets optional contact information
func (p *Pledge) SetContactInfo(name, email string) {
	p.pb.Contact = &pb.ContactInfo{
		Name:  name,
		Email: email,
	}
	p.id = p.calculateID()
}

// Transaction returns the underlying transaction
func (p *Pledge) Transaction() *transaction.Transaction {
	return p.tx
}

// Validate checks if the pledge is valid
func (p *Pledge) Validate() error {
	if p.tx == nil {
		return errors.New("no transaction")
	}

	if len(p.tx.Inputs) == 0 {
		return errors.New("no inputs")
	}

	if len(p.tx.Outputs) == 0 {
		return errors.New("no outputs")
	}

	// Check that all inputs have SIGHASH_ANYONECANPAY signatures
	for i, input := range p.tx.Inputs {
		if input.UnlockingScript == nil || len(*input.UnlockingScript) == 0 {
			return fmt.Errorf("input %d is not signed", i)
		}
		// TODO: Actually verify the signature has ANYONECANPAY flag
	}

	return nil
}