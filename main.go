	package main

	import (
	  "fmt"
	  "os"

	  "github.com/fatih/color"
	  "github.com/spf13/cobra"
	  "github.com/hassanzouhar/gitstatus/internal/config"
	  "github.com/hassanzouhar/gitstatus/pkg/checker"
	  "github.com/hassanzouhar/gitstatus/pkg/interactive"
	)

	var (
	  quietMode     bool
	  interactiveMode bool
	  Version        = "0.1.0"
	)

	// Exit codes
	const (
	  ExitSuccess          = 0
	  ExitError           = 1
	  ExitUncommitted     = 2
	  ExitUnpushed        = 3
	  ExitProtectedBranch = 4
	  ExitNeedsPull       = 5
	  ExitDiverged        = 6
	)

	func main() {
			    code, err := execute()
			    if err != nil {
			      color.Red("Error: %v", err)
			    }
			    os.Exit(code)
	}

	func execute() (int, error) {
	  rootCmd := &cobra.Command{
	    Use:   "gitstatus",
	    Short: "Git repository status checker",
	    Long: `A Git status checker that helps maintain repository hygiene by checking:
	      - Uncommitted changes
	      - Unpushed commits
	      - Branch protection
	      - Need for updates`,
	    RunE: runCheck,
	    Version: Version,
	  }

	  rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Suppress output except for errors")
	  rootCmd.PersistentFlags().BoolVarP(&interactiveMode, "interactive", "i", false, "Enable interactive mode")

			    err := rootCmd.Execute()
			    if err != nil {
			      return ExitError, err
			    }
			    return ExitSuccess, nil
			}

	func runCheck(cmd *cobra.Command, args []string) error {
	  cfg := config.New(quietMode, interactiveMode)
	  
	  status, err := checker.Check(cfg)
	  if err != nil {
	    return fmt.Errorf("failed to check status: %w", err)
	  }

	  if !quietMode {
	    status.Print()
	  }

			if interactiveMode && status.HasIssues() {
					    repo := status.Repository()
			  if err != nil {
			    return fmt.Errorf("failed to get repository: %w", err)
			  }
					    handler := interactive.NewHandler(repo, quietMode)
	    if err := handler.ProcessStatus(status); err != nil {
	      return fmt.Errorf("failed to process interactive actions: %w", err)
	    }
	    
	    // Recheck status after interactive operations
	    status, err = checker.Check(cfg)
	    if err != nil {
	      return fmt.Errorf("failed to recheck status: %w", err)
	    }
	    
	    if !quietMode {
	      color.Green("✓ All issues resolved")
	    }
	    
	    return nil
	  }

	  // Return status error only if not in interactive mode or if interactive mode didn't resolve all issues
	  return status.Error()
	}

