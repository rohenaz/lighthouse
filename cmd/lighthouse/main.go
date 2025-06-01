package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "lighthouse",
		Short: "Decentralized crowdfunding for BSV",
		Long: `Lighthouse is a decentralized crowdfunding platform for BSV using assurance contracts.
		
Create projects, make pledges, and claim funds trustlessly when funding goals are met.`,
		Version: version,
	}

	// Add commands
	rootCmd.AddCommand(
		projectCmd(),
		pledgeCmd(),
		serverCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// projectCmd handles project-related operations
func projectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage crowdfunding projects",
		Long:  "Create, view, and manage crowdfunding projects",
	}

	cmd.AddCommand(
		projectCreateCmd(),
		projectViewCmd(),
		projectStatusCmd(),
		projectClaimCmd(),
	)

	return cmd
}

// Command implementations are in project.go

// pledgeCmd handles pledge-related operations
func pledgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pledge",
		Short: "Manage pledges to projects",
		Long:  "Create, view, and revoke pledges to crowdfunding projects",
	}

	cmd.AddCommand(
		pledgeCreateCmd(),
		pledgeViewCmd(),
		pledgeRevokeCmd(),
	)

	return cmd
}

// Command implementations are in pledge.go

// serverCmd is now implemented in server.go
