package observability

import "context"

// Repository defines the interface for observability operations.
type Repository interface {
	GetTracer(name string) Tracer
	GetMeter(name string) Meter
	GetLogger(name string) Logger
	Shutdown(ctx context.Context) error
}
