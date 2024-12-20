	package checker

	import (
	  "github.com/fatih/color"
	  "github.com/hassanzouhar/gitstatus/internal/config"
	  "github.com/hassanzouhar/gitstatus/pkg/git"
	)

	type Status struct {
	    Branch          string
	    HasUncommitted bool
	    HasUnpushed    bool
	    IsProtected    bool
	    NeedsPull      bool
	    Config         *config.Config  // Made public for handler access
					    repo           git.Repository // Store repository instance
	}

	func Check(cfg *config.Config) (*Status, error) {
	  repo, err := git.OpenRepository(".", cfg.ProtectedBranches)
	  if err != nil {
	    return nil, err
	  }
	  
	  status := &Status{Config: cfg, repo: repo}

	  branch, err := repo.CurrentBranch()
	  if err != nil {
	    return nil, err
	  }
	  status.Branch = branch
	  status.IsProtected = repo.IsProtectedBranch(branch)

	  if hasUncommitted, err := repo.HasUncommittedChanges(); err != nil {
	    return nil, err
	  } else {
	    status.HasUncommitted = hasUncommitted
	  }

	  if hasUnpushed, err := repo.HasUnpushedCommits(); err != nil {
	    return nil, err
	  } else {
	    status.HasUnpushed = hasUnpushed
	  }

	  return status, nil
	}

	func (s *Status) Print() {
	  color.Blue("Current branch: %s", s.Branch)
	  
	  if s.IsProtected {
	    color.Yellow("⚠ Warning: You are on protected branch '%s'", s.Branch)
	  }
	  
	  if s.HasUncommitted {
	    color.Yellow("⚠ You have uncommitted changes")
	  }
	  
	  if s.HasUnpushed {
	    color.Yellow("⚠ You have unpushed commits")
	  }
	  
	  if s.NeedsPull {
	    color.Yellow("⚠ Branch needs to be pulled")
	  }
	}

		func (s *Status) HasIssues() bool {
		    return s.HasUncommitted || s.HasUnpushed || s.NeedsPull || s.IsProtected
		}

		func (s *Status) Error() error {
	  if !s.HasUncommitted && !s.HasUnpushed && !s.NeedsPull {
	    return nil
	  }
	  return &StatusError{status: s}
	}

	type StatusError struct {
	  status *Status
	}

	func (e *StatusError) Error() string {
	  return "repository has pending changes"
	}

	// Repository returns the underlying git.Repository instance
	func (s *Status) Repository() git.Repository {
	    return s.repo
	}
