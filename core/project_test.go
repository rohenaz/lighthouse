package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	t.Run("valid project creation", func(t *testing.T) {
		title := "Test Project"
		description := "This is a test project"
		goalAmount := uint64(100000000) // 1 BSV
		address := "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q"

		project, err := NewProject(title, description, goalAmount, address)
		require.NoError(t, err)
		assert.NotNil(t, project)

		// Check properties
		assert.Equal(t, title, project.Title())
		assert.Equal(t, description, project.Description())
		assert.Equal(t, goalAmount, project.GoalAmount())
		assert.Equal(t, uint64(10000), project.MinPledgeAmount()) // Default
		assert.NotEmpty(t, project.ID())
		assert.False(t, project.IsExpired())
	})

	t.Run("invalid address", func(t *testing.T) {
		project, err := NewProject("Test", "Description", 100000000, "invalid-address")
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "invalid address")
	})

	t.Run("zero goal amount", func(t *testing.T) {
		project, err := NewProject("Test", "Description", 0, "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q")
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "goal amount must be greater than 0")
	})

	t.Run("empty title", func(t *testing.T) {
		project, err := NewProject("", "Description", 100000000, "1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q")
		assert.Error(t, err)
		assert.Nil(t, project)
		assert.Contains(t, err.Error(), "title and description are required")
	})
}

func TestProjectSerialization(t *testing.T) {
	// Create a project
	project, err := NewProject(
		"Serialization Test",
		"Testing serialization",
		200000000, // 2 BSV
		"1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
	)
	require.NoError(t, err)

	// Serialize
	data, err := project.Serialize()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Load from serialized data
	loaded, err := LoadProject(data)
	require.NoError(t, err)
	assert.NotNil(t, loaded)

	// Compare properties
	assert.Equal(t, project.ID(), loaded.ID())
	assert.Equal(t, project.Title(), loaded.Title())
	assert.Equal(t, project.Description(), loaded.Description())
	assert.Equal(t, project.GoalAmount(), loaded.GoalAmount())
	assert.Equal(t, project.MinPledgeAmount(), loaded.MinPledgeAmount())
}

func TestProjectOutputs(t *testing.T) {
	project, err := NewProject(
		"Output Test",
		"Testing outputs",
		150000000, // 1.5 BSV
		"1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
	)
	require.NoError(t, err)

	outputs, err := project.Outputs()
	require.NoError(t, err)
	assert.Len(t, outputs, 1)

	// Check the output
	output := outputs[0]
	assert.Equal(t, uint64(150000000), output.Satoshis)
	assert.NotNil(t, output.LockingScript)
	assert.Greater(t, len(*output.LockingScript), 0)
}

func TestProjectCoverImage(t *testing.T) {
	project, err := NewProject(
		"Image Test",
		"Testing cover image",
		100000000,
		"1NKNazRR5jKgGqELVHDK47JAZrqtAWWy5q",
	)
	require.NoError(t, err)

	// Test JPEG header
	jpegData := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10} // Valid JPEG header
	err = project.SetCoverImage(jpegData)
	assert.NoError(t, err)

	// Test PNG header
	pngData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // Valid PNG header
	err = project.SetCoverImage(pngData)
	assert.NoError(t, err)

	// Test invalid image data
	invalidData := []byte{0x00, 0x01, 0x02, 0x03}
	err = project.SetCoverImage(invalidData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "image must be JPEG or PNG format")

	// Test too small data
	err = project.SetCoverImage([]byte{0xFF})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid image data")
}