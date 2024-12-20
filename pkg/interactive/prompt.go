	package interactive

	import (
	    "bufio"
	    "fmt"
	    "os"
	    "strings"

	    "github.com/fatih/color"
	    "github.com/hassanzouhar/gitstatus/pkg/checker"
	    "github.com/hassanzouhar/gitstatus/pkg/git"
	)

	// Handler handles interactive git operations
	type Handler struct {
	    repo git.Repository
	    quiet bool
	    reader *bufio.Reader
	}

	// NewHandler creates a new interactive handler
	func NewHandler(repo git.Repository, quiet bool) *Handler {
	    return &Handler{
	        repo:   repo,
	        quiet:  quiet,
	        reader: bufio.NewReader(os.Stdin),
	    }
	}

	// ProcessStatus handles all git status issues interactively
	func (h *Handler) ProcessStatus(status *checker.Status) error {
					      branch, err := h.repo.CurrentBranch()
					      if err != nil {
					          return fmt.Errorf("failed to get current branch: %w", err)
					      }
					      if status.IsProtected {
					          fmt.Printf("⚠ Warning: You are on protected branch '%s'\n", branch)
					      }

					      if status.HasUncommitted {
					          if err := h.handleUncommittedChanges(); err != nil {
					              return fmt.Errorf("failed to handle uncommitted changes: %w", err)
					          }
					      }

					      if status.HasUnpushed {
					          if err := h.handleUnpushedCommits(); err != nil {
					              return fmt.Errorf("failed to handle unpushed commits: %w", err)
					          }
					      }

	    if status.NeedsPull {
	        if err := h.handleNeedsPull(); err != nil {
	            return fmt.Errorf("failed to handle pull: %w", err)
	        }
	    }

	    return nil
	}

	func (h *Handler) handleUncommittedChanges() error {
	    if h.quiet {
	        return nil
	    }

	    fmt.Println("\n=== Recommended Actions ===")
	    fmt.Println("1. Commit your changes:")
	    fmt.Println("   git add .")
	    fmt.Println("   git commit -m \"your commit message\"")
	    
	    if !h.confirm("Would you like to commit all changes?") {
	        return nil
	    }

	    msg := h.prompt("Enter commit message")
	    if err := h.repo.CommitAll(msg); err != nil {
	        return fmt.Errorf("failed to commit: %w", err)
	    }

	    color.Green("✓ Changes committed successfully")
	    return nil
	}

	func (h *Handler) handleUnpushedCommits() error {
	    if h.quiet {
	        return nil
	    }

	    fmt.Println("2. Push your commits:")
					      branch, err := h.repo.CurrentBranch()
					      if err != nil {
					          return fmt.Errorf("failed to get current branch: %w", err)
					      }
					      fmt.Printf("   git push origin %s\n", branch)

	    if !h.confirm("Would you like to push your commits?") {
	        return nil
	    }

	    if err := h.repo.Push(); err != nil {
	        return fmt.Errorf("failed to push: %w", err)
	    }

	    color.Green("✓ Commits pushed successfully")
	    return nil
	}

	func (h *Handler) handleNeedsPull() error {
	    if h.quiet {
	        return nil
	    }

	    fmt.Println("3. Pull latest changes:")
					      branch, err := h.repo.CurrentBranch()
					      if err != nil {
					          return fmt.Errorf("failed to get current branch: %w", err)
					      }
					      fmt.Printf("   git pull origin %s\n", branch)

	    if !h.confirm("Would you like to pull the latest changes?") {
	        return nil
	    }

	    if err := h.repo.Pull(); err != nil {
	        return fmt.Errorf("failed to pull: %w", err)
	    }

	    color.Green("✓ Changes pulled successfully")
	    return nil
	}

	func (h *Handler) confirm(msg string) bool {
	    for {
	        fmt.Printf("%s (y/n): ", msg)
	        response, _ := h.reader.ReadString('\n')
	        response = strings.ToLower(strings.TrimSpace(response))
	        
	        if response == "y" || response == "yes" {
	            return true
	        }
	        if response == "n" || response == "no" {
	            return false
	        }
	    }
	}

	func (h *Handler) prompt(msg string) string {
	    fmt.Printf("%s: ", msg)
	    response, _ := h.reader.ReadString('\n')
	    return strings.TrimSpace(response)
	}

