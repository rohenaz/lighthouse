package core

import (
	"errors"
	"fmt"

	"github.com/bsv-blockchain/go-sdk/transaction"
)

// Contract represents an assurance contract that combines pledges
type Contract struct {
	project  *Project
	pledges  []*Pledge
	combined *transaction.Transaction
}

// NewContract creates a new assurance contract for a project
func NewContract(project *Project) *Contract {
	return &Contract{
		project: project,
		pledges: make([]*Pledge, 0),
	}
}

// AddPledge adds a pledge to the contract
func (c *Contract) AddPledge(pledge *Pledge) error {
	// Verify pledge is for this project
	if pledge.ProjectID() != c.project.ID() {
		return errors.New("pledge is for different project")
	}

	// Validate the pledge
	if err := pledge.Validate(); err != nil {
		return fmt.Errorf("invalid pledge: %w", err)
	}

	// Check for duplicate pledges (same inputs)
	for _, existing := range c.pledges {
		if c.hasDuplicateInputs(existing, pledge) {
			return errors.New("pledge uses same inputs as existing pledge")
		}
	}

	c.pledges = append(c.pledges, pledge)
	return nil
}

// TotalPledged returns the total amount pledged so far
func (c *Contract) TotalPledged() uint64 {
	total := uint64(0)
	for _, pledge := range c.pledges {
		total += pledge.Amount()
	}
	return total
}

// Progress returns the funding progress as a percentage
func (c *Contract) Progress() float64 {
	return float64(c.TotalPledged()) / float64(c.project.GoalAmount()) * 100
}

// CanClaim checks if the contract can be claimed (goal reached)
func (c *Contract) CanClaim() bool {
	return c.TotalPledged() >= c.project.GoalAmount()
}

// Combine creates the final transaction from all pledges
func (c *Contract) Combine() (*transaction.Transaction, error) {
	if !c.CanClaim() {
		return nil, fmt.Errorf("funding goal not reached: %d/%d", c.TotalPledged(), c.project.GoalAmount())
	}

	// Create a new transaction
	tx := transaction.NewTransaction()

	// Add all inputs from all pledges
	inputValue := uint64(0)
	for _, pledge := range c.pledges {
		for _, input := range pledge.Transaction().Inputs {
			tx.Inputs = append(tx.Inputs, input)
			// In a real implementation, we'd look up the input value
			// For now, we'll trust the pledge amount
		}
		inputValue += pledge.Amount()
	}

	// Add the project outputs
	outputs, err := c.project.Outputs()
	if err != nil {
		return nil, fmt.Errorf("failed to get project outputs: %w", err)
	}

	outputValue := uint64(0)
	for _, out := range outputs {
		tx.AddOutput(out)
		outputValue += out.Satoshis
	}

	// Calculate and add change if necessary
	if inputValue > outputValue {
		change := inputValue - outputValue
		// In a real implementation, we'd need to determine where change goes
		// For now, we'll add it to the first output
		if len(tx.Outputs) > 0 {
			tx.Outputs[0].Satoshis += change
		}
	}

	c.combined = tx
	return tx, nil
}

// Transaction returns the combined transaction if available
func (c *Contract) Transaction() *transaction.Transaction {
	return c.combined
}

// RemovePledge removes a pledge from the contract
func (c *Contract) RemovePledge(pledgeID string) error {
	for i, pledge := range c.pledges {
		if pledge.ID() == pledgeID {
			c.pledges = append(c.pledges[:i], c.pledges[i+1:]...)
			return nil
		}
	}
	return errors.New("pledge not found")
}

// Pledges returns all pledges in the contract
func (c *Contract) Pledges() []*Pledge {
	return c.pledges
}

// hasDuplicateInputs checks if two pledges share any inputs
func (c *Contract) hasDuplicateInputs(p1, p2 *Pledge) bool {
	inputs1 := make(map[string]bool)
	
	// Build a map of all inputs in first pledge
	for _, input := range p1.Transaction().Inputs {
		key := fmt.Sprintf("%x:%d", input.SourceTXID, input.SourceTxOutIndex)
		inputs1[key] = true
	}

	// Check if any inputs in second pledge exist in the map
	for _, input := range p2.Transaction().Inputs {
		key := fmt.Sprintf("%x:%d", input.SourceTXID, input.SourceTxOutIndex)
		if inputs1[key] {
			return true
		}
	}

	return false
}

// ValidatePledges verifies all pledges are still valid (unspent)
func (c *Contract) ValidatePledges() error {
	// In a real implementation, this would check the blockchain
	// to ensure all pledge inputs are still unspent
	for i, pledge := range c.pledges {
		if err := pledge.Validate(); err != nil {
			return fmt.Errorf("pledge %d invalid: %w", i, err)
		}
	}
	return nil
}

// Status returns the current status of the contract
type ContractStatus struct {
	ProjectID      string
	GoalAmount     uint64
	TotalPledged   uint64
	PledgeCount    int
	Progress       float64
	CanClaim       bool
	IsExpired      bool
}

// GetStatus returns the current contract status
func (c *Contract) GetStatus() ContractStatus {
	return ContractStatus{
		ProjectID:    c.project.ID(),
		GoalAmount:   c.project.GoalAmount(),
		TotalPledged: c.TotalPledged(),
		PledgeCount:  len(c.pledges),
		Progress:     c.Progress(),
		CanClaim:     c.CanClaim(),
		IsExpired:    c.project.IsExpired(),
	}
}