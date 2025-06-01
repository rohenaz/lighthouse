package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/template/p2pkh"
	pb "github.com/yourusername/lighthouse/core/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Project represents a crowdfunding project
type Project struct {
	pb       *pb.Project
	id       string
	goalAmount uint64
}

// NewProject creates a new crowdfunding project
func NewProject(title, description string, goalAmount uint64, address string) (*Project, error) {
	if title == "" || description == "" {
		return nil, errors.New("title and description are required")
	}
	if goalAmount == 0 {
		return nil, errors.New("goal amount must be greater than 0")
	}

	// Parse address
	addr, err := script.NewAddressFromString(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	// Create P2PKH locking script from address
	lockingScript, err := p2pkh.Lock(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create locking script: %w", err)
	}

	// Create the project protobuf
	proj := &pb.Project{
		Version: 1,
		Details: &pb.ProjectDetails{
			Network: "mainnet",
			Outputs: []*pb.Output{{
				Amount: goalAmount,
				Script: lockingScript.Bytes(),
			}},
			Time: timestamppb.Now(),
			Memo: description,
		},
		Extra: &pb.ProjectExtraDetails{
			Title:           title,
			MinPledgeAmount: 10000, // 0.0001 BSV default minimum
		},
	}

	p := &Project{
		pb:         proj,
		goalAmount: goalAmount,
	}
	p.id = p.calculateID()

	return p, nil
}

// LoadProject loads a project from serialized data
func LoadProject(data []byte) (*Project, error) {
	var proj pb.Project
	if err := proto.Unmarshal(data, &proj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project: %w", err)
	}

	p := &Project{pb: &proj}
	
	// Calculate total goal amount from outputs
	for _, output := range proj.Details.Outputs {
		p.goalAmount += output.Amount
	}
	
	p.id = p.calculateID()
	return p, nil
}

// Serialize returns the project as protobuf bytes
func (p *Project) Serialize() ([]byte, error) {
	return proto.Marshal(p.pb)
}

// ID returns the unique project ID (hash of serialized data)
func (p *Project) ID() string {
	return p.id
}

// calculateID generates a unique ID from project data
func (p *Project) calculateID() string {
	data, _ := p.Serialize()
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Title returns the project title
func (p *Project) Title() string {
	if p.pb.Extra != nil {
		return p.pb.Extra.Title
	}
	return ""
}

// Description returns the project description
func (p *Project) Description() string {
	if p.pb.Details != nil {
		return p.pb.Details.Memo
	}
	return ""
}

// GoalAmount returns the funding goal in satoshis
func (p *Project) GoalAmount() uint64 {
	return p.goalAmount
}

// MinPledgeAmount returns the minimum pledge in satoshis
func (p *Project) MinPledgeAmount() uint64 {
	if p.pb.Extra != nil && p.pb.Extra.MinPledgeAmount > 0 {
		return p.pb.Extra.MinPledgeAmount
	}
	return 10000 // Default 0.0001 BSV
}

// IsExpired checks if the project has expired
func (p *Project) IsExpired() bool {
	if p.pb.Details == nil || p.pb.Details.Expires == nil {
		return false
	}
	return p.pb.Details.Expires.AsTime().Before(time.Now())
}

// Outputs returns the transaction outputs for this project
func (p *Project) Outputs() ([]*transaction.TransactionOutput, error) {
	if p.pb.Details == nil {
		return nil, errors.New("project has no details")
	}

	var outputs []*transaction.TransactionOutput
	for _, out := range p.pb.Details.Outputs {
		lockScript := script.Script(out.Script)
		txOut := &transaction.TransactionOutput{
			Satoshis:      out.Amount,
			LockingScript: &lockScript,
		}
		outputs = append(outputs, txOut)
	}

	return outputs, nil
}

// SetAuthKey sets the authentication key for project ownership
func (p *Project) SetAuthKey(pubKey []byte) {
	if p.pb.Extra == nil {
		p.pb.Extra = &pb.ProjectExtraDetails{}
	}
	p.pb.Extra.AuthKey = pubKey
	p.id = p.calculateID() // Recalculate ID
}

// SetCoverImage sets the project cover image
func (p *Project) SetCoverImage(imageData []byte) error {
	// Basic validation - check for JPEG or PNG header
	if len(imageData) < 4 {
		return errors.New("invalid image data")
	}
	
	// Check for JPEG (FF D8 FF) or PNG (89 50 4E 47) headers
	isJPEG := imageData[0] == 0xFF && imageData[1] == 0xD8 && imageData[2] == 0xFF
	isPNG := imageData[0] == 0x89 && imageData[1] == 0x50 && imageData[2] == 0x4E && imageData[3] == 0x47
	
	if !isJPEG && !isPNG {
		return errors.New("image must be JPEG or PNG format")
	}

	if p.pb.Extra == nil {
		p.pb.Extra = &pb.ProjectExtraDetails{}
	}
	p.pb.Extra.CoverImage = imageData
	p.id = p.calculateID() // Recalculate ID
	
	return nil
}