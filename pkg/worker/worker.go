package worker

import "context"

// ExitCode represents exit status.
type ExitCode int

const (
	NormalExit ExitCode = iota
	PrematureExit
	RuntimeErrorExit
	PanicExit
)

// Worker is the interface that concrete workers implement.
type Worker interface {
	Run(ctx context.Context, config map[string]interface{}) (ExitCode, error)
	GetWorkerName() string
}
