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
	    quietMode       bool
	    interactiveMode bool
	    Version        = "0.1.0"
	    rootCmd        *cobra.Command
	)

	// GitStatusError represents repository status-related errors
	type GitStatusError struct {
	    code int
	    msg  string
	}

	func (e *GitStatusError) Error() string {
	    return e.msg
	}

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
	        if statusErr, isStatusErr := err.(*GitStatusError); isStatusErr {
	            color.Red("Error: %v", statusErr)
	        } else {
	            fmt.Println(err)
	            rootCmd.Help()
	        }
	    }
	    os.Exit(code)
	}

	func execute() (int, error) {
	    rootCmd = &cobra.Command{
	        Use:     "gitstatus",
	        Short:   "Git repository status checker",
	        Long: `A Git status checker that helps maintain repository hygiene by checking:
	- Uncommitted changes
	- Unpushed commits
	- Branch protection
	- Need for updates`,
	        RunE:    runCheck,
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

	    hadIssues := status.HasIssues()
	    
	    if !quietMode {
	        status.Print()
	        if !hadIssues {
	            color.Green("\n✓ Repository is clean")
	            return nil
	        }
	    }

	    if interactiveMode && hadIssues {
	        repo := status.Repository()
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
	            if status.HasIssues() {
	                color.Yellow("\n! Some issues remain unresolved")
	            } else {
	                color.Green("\n✓ All issues were resolved successfully")
	            }
	        }
	        
	        if !status.HasIssues() {
	            return nil
	        }
	    }

	    // Return status error only if not in interactive mode or if interactive mode didn't resolve all issues
	    if err := status.Error(); err != nil {
	        if status.HasUncommitted {
	            return &GitStatusError{code: ExitUncommitted, msg: "repository has uncommitted changes"}
	        }
	        if status.HasUnpushed {
	            return &GitStatusError{code: ExitUnpushed, msg: "repository has unpushed commits"}
	        }
	        if status.NeedsPull {
	            return &GitStatusError{code: ExitNeedsPull, msg: "repository needs to be pulled"}
	        }
	        if status.IsProtected {
	            return &GitStatusError{code: ExitProtectedBranch, msg: fmt.Sprintf("branch '%s' is protected", status.Branch)}
	        }
	        return &GitStatusError{code: ExitError, msg: err.Error()}
	    }
	    return nil
	}

