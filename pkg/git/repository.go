	package git

													import (
													    "fmt"
													    "os/exec"
													    "github.com/go-git/go-git/v5"
													    "github.com/go-git/go-git/v5/plumbing"
													    "github.com/go-git/go-git/v5/plumbing/object"
													)

	type Repository interface {
	  CurrentBranch() (string, error)
	  HasUncommittedChanges() (bool, error)
	  HasUnpushedCommits() (bool, error)
	  IsProtectedBranch(branch string) bool
	  Push() error
	  Pull() error
	  CommitAll(message string) error
	  CreateAndCheckoutBranch(branchName string) error
	}

	type GitRepo struct {
	  repo *git.Repository
	  protectedBranches []string
	}

	func OpenRepository(path string, protectedBranches []string) (Repository, error) {
	  repo, err := git.PlainOpen(path)
	  if err != nil {
	    return nil, fmt.Errorf("failed to open repository: %w", err)
	  }

	  return &GitRepo{
	    repo: repo,
	    protectedBranches: protectedBranches,
	  }, nil
	}

	func (r *GitRepo) CurrentBranch() (string, error) {
	  head, err := r.repo.Head()
	  if err != nil {
	    return "", fmt.Errorf("failed to get HEAD: %w", err)
	  }
	  return head.Name().Short(), nil
	}

	func (r *GitRepo) HasUncommittedChanges() (bool, error) {
	  worktree, err := r.repo.Worktree()
	  if err != nil {
	    return false, fmt.Errorf("failed to get worktree: %w", err)
	  }

	  status, err := worktree.Status()
	  if err != nil {
	    return false, fmt.Errorf("failed to get status: %w", err)
	  }

	  return !status.IsClean(), nil
	  }

	  func (r *GitRepo) HasUnpushedCommits() (bool, error) {
	    head, err := r.repo.Head()
	    if err != nil {
	      return false, fmt.Errorf("failed to get HEAD: %w", err)
	    }

					    // Get the remote tracking branch
					    remoteRef, err := r.repo.Reference(plumbing.NewRemoteReferenceName("origin", head.Name().Short()), true)
	    if err != nil {
	      if err == plumbing.ErrReferenceNotFound {
	        // If remote branch doesn't exist, we have unpushed commits
	        return true, nil
	      }
	      return false, fmt.Errorf("failed to get remote reference: %w", err)
	    }

	    // Compare local and remote HEADs
	    return head.Hash() != remoteRef.Hash(), nil
	  }

			    func execGitCommand(args ...string) error {
			      cmd := exec.Command("git", args...)
			      output, err := cmd.CombinedOutput()
			      if err != nil {
			        return fmt.Errorf("git command failed: %s: %w", string(output), err)
			      }
			      return nil
			    }

			    func (r *GitRepo) Push() error {
			      branch, err := r.CurrentBranch()
			      if err != nil {
			        return fmt.Errorf("failed to get current branch: %w", err)
			      }
			      return execGitCommand("push", "origin", branch)
							}

			    func (r *GitRepo) Pull() error {
			      branch, err := r.CurrentBranch()
			      if err != nil {
			        return fmt.Errorf("failed to get current branch: %w", err)
			      }
			      
			      err = execGitCommand("pull", "origin", branch)
			      if err != nil {
			        return fmt.Errorf("failed to pull: %w", err)
			      }
			      return nil
			    }

	  func (r *GitRepo) IsProtectedBranch(branch string) bool {
	    for _, protected := range r.protectedBranches {
	      if protected == branch {
	        return true
	      }
	    }
	    return false
							}

			    func (r *GitRepo) CommitAll(message string) error {
			      w, err := r.repo.Worktree()
			      if err != nil {
			        return fmt.Errorf("failed to get worktree: %w", err)
			      }

			      // Stage all changes
																			if _, err := w.Add("."); err != nil {
																			    return fmt.Errorf("failed to stage changes: %w", err)
																			}

			      // Create commit
			      commit, err := w.Commit(message, &git.CommitOptions{
			        Author: &object.Signature{
			          Name:  "Git Status Checker",
			          Email: "gitstatus@local",
			        },
			      })
			      if err != nil {
			        return fmt.Errorf("failed to commit changes: %w", err)
			      }

			      _, err = r.repo.CommitObject(commit)
			      if err != nil {
			        return fmt.Errorf("failed to get commit object: %w", err)
			      }

			      return nil
							}

							func (r *GitRepo) CreateAndCheckoutBranch(branchName string) error {
							  // Get the worktree
							  w, err := r.repo.Worktree()
							  if err != nil {
							    return fmt.Errorf("failed to get worktree: %w", err)
							  }

							  // Get HEAD reference
							  head, err := r.repo.Head()
							  if err != nil {
							    return fmt.Errorf("failed to get HEAD reference: %w", err)
							  }

							  // Create new reference for the branch
							  branchRef := plumbing.NewHashReference(
							    plumbing.NewBranchReferenceName(branchName),
							    head.Hash(),
							  )

							  // Create the branch in the repository
							  err = r.repo.Storer.SetReference(branchRef)
							  if err != nil {
							    return fmt.Errorf("failed to create branch %s: %w", branchName, err)
							  }

							  // Checkout the new branch
							  err = w.Checkout(&git.CheckoutOptions{
							    Branch: branchRef.Name(),
							    Create: false,
							    Force:  false,
							  })
							  if err != nil {
							    return fmt.Errorf("failed to checkout branch %s: %w", branchName, err)
							  }

							  return nil
							}
