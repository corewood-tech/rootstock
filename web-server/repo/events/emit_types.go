package events

// WorkflowContext provides access to the workflow runtime.
// It wraps vendor-specific context types so callers never import the vendor.
type WorkflowContext interface{}
