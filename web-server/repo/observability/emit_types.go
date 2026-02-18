package observability

import "context"

// Logger provides structured logging methods.
type Logger interface {
	Info(ctx context.Context, msg string, attrs map[string]interface{})
	Error(ctx context.Context, msg string, attrs map[string]interface{})
	Warn(ctx context.Context, msg string, attrs map[string]interface{})
	Debug(ctx context.Context, msg string, attrs map[string]interface{})
}

// Tracer creates spans for distributed tracing.
type Tracer interface {
	StartSpan(ctx context.Context, name string) (context.Context, Span)
}

// Span represents a unit of work in a trace.
type Span interface {
	End()
	SetAttribute(key string, value interface{})
	RecordError(err error)
}

// Meter creates metric instruments.
type Meter interface {
	Counter(name string) Counter
	Histogram(name string) Histogram
}

// Counter is a monotonically increasing metric.
type Counter interface {
	Add(ctx context.Context, value float64)
}

// Histogram records a distribution of values.
type Histogram interface {
	Record(ctx context.Context, value float64)
}
