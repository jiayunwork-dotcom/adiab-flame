package cli

// exit codes returned by the CLI. User input mistakes use 2, runtime
// failures (bad config, non-convergence) use 1.
const (
	// ExitOK is a successful run.
	ExitOK = 0
	// ExitRuntime is a calculation or file failure.
	ExitRuntime = 1
	// ExitUsage is a command line mistake.
	ExitUsage = 2
)
