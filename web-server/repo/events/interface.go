package events

// Repository defines the interface for workflow event operations.
type Repository interface {
	GetContext() WorkflowContext
	Shutdown()
}
