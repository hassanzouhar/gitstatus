	package config

	type Config struct {
	  QuietMode         bool
	  InteractiveMode   bool
	  ProtectedBranches []string
	}

	const (
	  ExitCodeOK                = 0
	  ExitCodeError            = 1
	  ExitCodeUncommitted      = 2
	  ExitCodeUnpushed         = 3
	  ExitCodeNeedsPull        = 4
	  ExitCodeProtectedBranch  = 5
	)

	func New(quiet, interactive bool) *Config {
	  return &Config{
	    QuietMode:       quiet,
	    InteractiveMode: interactive,
	    ProtectedBranches: []string{},
	  }
	}

	func (c *Config) SetProtectedBranches(branches []string) {
	  c.ProtectedBranches = branches
	}

	func (c *Config) UnprotectAllBranches() {
	  c.ProtectedBranches = []string{}
	}

