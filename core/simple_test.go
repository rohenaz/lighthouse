package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Simple tests that don't involve complex UTXO creation
func TestProjectBasics(t *testing.T) {
	t.Run("create and serialize project", func(t *testing.T) {
		project, err := NewProject(
			"Simple Test",
			"A simple test project",
			100000000, // 1 BSV
			"1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
		)
		require.NoError(t, err)
		assert.NotNil(t, project)

		// Test basic properties
		assert.Equal(t, "Simple Test", project.Title())
		assert.Equal(t, "A simple test project", project.Description())
		assert.Equal(t, uint64(100000000), project.GoalAmount())
		assert.NotEmpty(t, project.ID())

		// Test serialization roundtrip
		data, err := project.Serialize()
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		loaded, err := LoadProject(data)
		require.NoError(t, err)
		assert.Equal(t, project.ID(), loaded.ID())
		assert.Equal(t, project.Title(), loaded.Title())
		assert.Equal(t, project.Description(), loaded.Description())
		assert.Equal(t, project.GoalAmount(), loaded.GoalAmount())
	})
}

func TestContractBasics(t *testing.T) {
	t.Run("create empty contract", func(t *testing.T) {
		project, err := NewProject(
			"Contract Test",
			"Test contract functionality",
			200000000, // 2 BSV
			"1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
		)
		require.NoError(t, err)

		contract := NewContract(project)
		assert.NotNil(t, contract)

		// Test initial state
		assert.Equal(t, uint64(0), contract.TotalPledged())
		assert.Equal(t, 0.0, contract.Progress())
		assert.False(t, contract.CanClaim())
		assert.Len(t, contract.Pledges(), 0)

		// Test status
		status := contract.GetStatus()
		assert.Equal(t, project.ID(), status.ProjectID)
		assert.Equal(t, uint64(200000000), status.GoalAmount)
		assert.Equal(t, uint64(0), status.TotalPledged)
		assert.Equal(t, 0, status.PledgeCount)
		assert.Equal(t, 0.0, status.Progress)
		assert.False(t, status.CanClaim)
		assert.False(t, status.IsExpired)
	})
}

func TestProjectValidation(t *testing.T) {
	testCases := []struct {
		name        string
		title       string
		description string
		goal        uint64
		address     string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid project",
			title:       "Valid Project",
			description: "A valid project description",
			goal:        100000000,
			address:     "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
			shouldError: false,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Description",
			goal:        100000000,
			address:     "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
			shouldError: true,
			errorMsg:    "title and description are required",
		},
		{
			name:        "empty description",
			title:       "Title",
			description: "",
			goal:        100000000,
			address:     "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
			shouldError: true,
			errorMsg:    "title and description are required",
		},
		{
			name:        "zero goal",
			title:       "Title",
			description: "Description",
			goal:        0,
			address:     "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
			shouldError: true,
			errorMsg:    "goal amount must be greater than 0",
		},
		{
			name:        "invalid address",
			title:       "Title",
			description: "Description",
			goal:        100000000,
			address:     "invalid-address",
			shouldError: true,
			errorMsg:    "invalid address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			project, err := NewProject(tc.title, tc.description, tc.goal, tc.address)
			
			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, project)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, project)
			}
		})
	}
}