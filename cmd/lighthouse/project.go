package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourusername/lighthouse/core"
)

// projectCreateCmd creates a new project
func projectCreateCmd() *cobra.Command {
	var (
		goal        float64
		address     string
		description string
		minPledge   float64
		expiry      int
		output      string
	)

	cmd := &cobra.Command{
		Use:   "create [title]",
		Short: "Create a new crowdfunding project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			
			// Convert BSV to satoshis
			goalSatoshis := uint64(goal * 100000000)
			minPledgeSatoshis := uint64(minPledge * 100000000)
			
			// Create the project
			project, err := core.NewProject(title, description, goalSatoshis, address)
			if err != nil {
				return fmt.Errorf("failed to create project: %w", err)
			}
			
			// Set minimum pledge if different from default
			if minPledgeSatoshis > 0 && minPledgeSatoshis != project.MinPledgeAmount() {
				// We'd need to add a method to update this in the project
				// For now, the default is used
			}
			
			// Serialize the project
			data, err := project.Serialize()
			if err != nil {
				return fmt.Errorf("failed to serialize project: %w", err)
			}
			
			// Determine output filename
			if output == "" {
				output = fmt.Sprintf("%s.lighthouse", sanitizeFilename(title))
			}
			
			// Write to file
			if err := ioutil.WriteFile(output, data, 0644); err != nil {
				return fmt.Errorf("failed to write project file: %w", err)
			}
			
			fmt.Printf("Project created successfully!\n")
			fmt.Printf("File: %s\n", output)
			fmt.Printf("ID: %s\n", project.ID())
			fmt.Printf("Goal: %.8f BSV (%d satoshis)\n", goal, goalSatoshis)
			fmt.Printf("Address: %s\n", address)
			fmt.Printf("Minimum pledge: %.8f BSV\n", minPledge)
			
			return nil
		},
	}

	cmd.Flags().Float64VarP(&goal, "goal", "g", 0, "Funding goal in BSV (required)")
	cmd.Flags().StringVarP(&address, "address", "a", "", "BSV address to receive funds (required)")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Project description")
	cmd.Flags().Float64VarP(&minPledge, "min-pledge", "m", 0.0001, "Minimum pledge amount in BSV")
	cmd.Flags().IntVarP(&expiry, "expiry", "e", 0, "Days until project expires (0 = no expiry)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output filename (default: title.lighthouse)")

	cmd.MarkFlagRequired("goal")
	cmd.MarkFlagRequired("address")

	return cmd
}

// projectViewCmd displays project details
func projectViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [project-file]",
		Short: "View project details",
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
			
			// Display project details
			fmt.Printf("Project: %s\n", project.Title())
			fmt.Printf("ID: %s\n", project.ID())
			fmt.Printf("Description: %s\n", project.Description())
			fmt.Printf("Goal: %.8f BSV (%d satoshis)\n", 
				float64(project.GoalAmount())/100000000, project.GoalAmount())
			fmt.Printf("Minimum pledge: %.8f BSV\n", 
				float64(project.MinPledgeAmount())/100000000)
			
			if project.IsExpired() {
				fmt.Printf("Status: EXPIRED\n")
			} else {
				fmt.Printf("Status: Active\n")
			}
			
			return nil
		},
	}
}

// projectStatusCmd shows project funding status
func projectStatusCmd() *cobra.Command {
	var pledgeDir string
	
	cmd := &cobra.Command{
		Use:   "status [project-file]",
		Short: "Check project funding status",
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
			
			// Create a contract
			contract := core.NewContract(project)
			
			// Load pledges from directory
			if pledgeDir == "" {
				pledgeDir = filepath.Dir(projectFile)
			}
			
			pledgeFiles, err := filepath.Glob(filepath.Join(pledgeDir, "*.pledge"))
			if err != nil {
				return fmt.Errorf("failed to list pledge files: %w", err)
			}
			
			// Load each pledge
			for _, pledgeFile := range pledgeFiles {
				pledgeData, err := ioutil.ReadFile(pledgeFile)
				if err != nil {
					fmt.Printf("Warning: failed to read pledge file %s: %v\n", pledgeFile, err)
					continue
				}
				
				pledge, err := core.LoadPledge(pledgeData)
				if err != nil {
					fmt.Printf("Warning: failed to load pledge from %s: %v\n", pledgeFile, err)
					continue
				}
				
				if err := contract.AddPledge(pledge); err != nil {
					fmt.Printf("Warning: failed to add pledge from %s: %v\n", pledgeFile, err)
					continue
				}
			}
			
			// Display status
			status := contract.GetStatus()
			fmt.Printf("Project: %s\n", project.Title())
			fmt.Printf("Goal: %.8f BSV\n", float64(status.GoalAmount)/100000000)
			fmt.Printf("Pledged: %.8f BSV (%.1f%%)\n", 
				float64(status.TotalPledged)/100000000, status.Progress)
			fmt.Printf("Pledges: %d\n", status.PledgeCount)
			
			if status.CanClaim {
				fmt.Printf("Status: READY TO CLAIM! 🎉\n")
			} else if status.IsExpired {
				fmt.Printf("Status: EXPIRED\n")
			} else {
				fmt.Printf("Status: Active (%.1f%% funded)\n", status.Progress)
			}
			
			return nil
		},
	}
	
	cmd.Flags().StringVarP(&pledgeDir, "pledge-dir", "p", "", "Directory containing pledge files (default: same as project)")
	
	return cmd
}

// projectClaimCmd claims funds when goal is reached
func projectClaimCmd() *cobra.Command {
	var (
		broadcast bool
		pledgeDir string
		output    string
	)

	cmd := &cobra.Command{
		Use:   "claim [project-file]",
		Short: "Claim funds when funding goal is reached",
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
			
			// Create a contract
			contract := core.NewContract(project)
			
			// Load pledges from directory
			if pledgeDir == "" {
				pledgeDir = filepath.Dir(projectFile)
			}
			
			pledgeFiles, err := filepath.Glob(filepath.Join(pledgeDir, "*.pledge"))
			if err != nil {
				return fmt.Errorf("failed to list pledge files: %w", err)
			}
			
			if len(pledgeFiles) == 0 {
				return fmt.Errorf("no pledge files found in %s", pledgeDir)
			}
			
			// Load each pledge
			fmt.Printf("Loading %d pledges...\n", len(pledgeFiles))
			for _, pledgeFile := range pledgeFiles {
				pledgeData, err := ioutil.ReadFile(pledgeFile)
				if err != nil {
					fmt.Printf("Warning: failed to read pledge file %s: %v\n", pledgeFile, err)
					continue
				}
				
				pledge, err := core.LoadPledge(pledgeData)
				if err != nil {
					fmt.Printf("Warning: failed to load pledge from %s: %v\n", pledgeFile, err)
					continue
				}
				
				if err := contract.AddPledge(pledge); err != nil {
					fmt.Printf("Warning: failed to add pledge from %s: %v\n", pledgeFile, err)
					continue
				}
			}
			
			// Check if we can claim
			if !contract.CanClaim() {
				status := contract.GetStatus()
				return fmt.Errorf("cannot claim: only %.1f%% funded (%.8f/%.8f BSV)", 
					status.Progress,
					float64(status.TotalPledged)/100000000,
					float64(status.GoalAmount)/100000000)
			}
			
			// Combine the transaction
			tx, err := contract.Combine()
			if err != nil {
				return fmt.Errorf("failed to combine transaction: %w", err)
			}
			
			// Save the transaction
			txHex := tx.String()
			if output == "" {
				output = fmt.Sprintf("%s-claim.tx", projectFile)
			}
			
			if err := ioutil.WriteFile(output, []byte(txHex), 0644); err != nil {
				return fmt.Errorf("failed to write transaction: %w", err)
			}
			
			fmt.Printf("Claim transaction created!\n")
			fmt.Printf("File: %s\n", output)
			fmt.Printf("Transaction ID: %s\n", tx.TxID())
			fmt.Printf("Total amount: %.8f BSV\n", float64(contract.TotalPledged())/100000000)
			
			if broadcast {
				fmt.Printf("\nBroadcasting transaction...\n")
				// TODO: Implement actual broadcasting
				fmt.Printf("Broadcasting not yet implemented. Use a BSV node or service to broadcast the transaction.\n")
			} else {
				fmt.Printf("\nTo broadcast, use: lighthouse broadcast %s\n", output)
			}
			
			return nil
		},
	}

	cmd.Flags().BoolVarP(&broadcast, "broadcast", "b", false, "Broadcast the claim transaction")
	cmd.Flags().StringVarP(&pledgeDir, "pledge-dir", "p", "", "Directory containing pledge files (default: same as project)")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output transaction file (default: project-claim.tx)")

	return cmd
}

// sanitizeFilename removes invalid characters from filenames
func sanitizeFilename(name string) string {
	// Replace spaces with underscores
	name = filepath.Clean(name)
	
	// Remove or replace invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for range invalid {
		name = filepath.Base(name) // Get base name to remove path separators
	}
	
	// Replace spaces with underscores
	for i := 0; i < len(name); i++ {
		if name[i] == ' ' {
			name = name[:i] + "_" + name[i+1:]
		}
	}
	
	return name
}