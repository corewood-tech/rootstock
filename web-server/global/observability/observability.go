package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"rootstock/web-server/config"
)

// Logger provides structured logging methods.
type Logger interface {
	Info(ctx context.Context, msg string, attrs map[string]interface{})
	Error(ctx context.Context, msg string, attrs map[string]interface{})
	Warn(ctx context.Context, msg string, attrs map[string]interface{})
	Debug(ctx context.Context, msg string, attrs map[string]interface{})
}

var (
	mu             sync.RWMutex
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	logProvider    *sdklog.LoggerProvider
	initialized    bool
)

// Initialize sets up the global OTel providers based on the given config.
func Initialize(ctx context.Context, cfg config.ObservabilityConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}

	if cfg.EnableTraces {
		if err := initTracer(ctx, cfg, res); err != nil {
			return fmt.Errorf("init tracer: %w", err)
		}
	}

	if cfg.EnableMetrics {
		if err := initMeter(ctx, cfg, res); err != nil {
			return fmt.Errorf("init meter: %w", err)
		}
	}

	if cfg.EnableLogs {
		if err := initLogger(ctx, cfg, res); err != nil {
			return fmt.Errorf("init logger: %w", err)
		}
	}

	initialized = true
	return nil
}

// Shutdown gracefully shuts down all providers.
func Shutdown(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	var errs []error

	if tracerProvider != nil {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown tracer: %w", err))
		}
		tracerProvider = nil
	}

	if meterProvider != nil {
		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown meter: %w", err))
		}
		meterProvider = nil
	}

	if logProvider != nil {
		if err := logProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown logger: %w", err))
		}
		logProvider = nil
	}

	initialized = false

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}

// GetTracer returns a named tracer from the global provider.
func GetTracer(name string) trace.Tracer {
	mu.RLock()
	defer mu.RUnlock()

	if tracerProvider != nil {
		return tracerProvider.Tracer(name)
	}
	return otel.Tracer(name)
}

// GetMeter returns a named meter from the global provider.
func GetMeter(name string) metric.Meter {
	mu.RLock()
	defer mu.RUnlock()

	if meterProvider != nil {
		return meterProvider.Meter(name)
	}
	return otel.Meter(name)
}

// GetLogger returns a Logger backed by OTel if available, otherwise slog.
func GetLogger(name string) Logger {
	mu.RLock()
	defer mu.RUnlock()

	var handler slog.Handler
	if logProvider != nil {
		handler = otelslog.NewHandler(name, otelslog.WithLoggerProvider(logProvider))
	} else {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}

	return &slogLogger{
		logger: slog.New(handler).With("component", name),
	}
}

// slogLogger implements Logger using the standard library slog package.
type slogLogger struct {
	logger *slog.Logger
}

func (l *slogLogger) Info(ctx context.Context, msg string, attrs map[string]interface{}) {
	l.logger.InfoContext(ctx, msg, attrsToSlogArgs(attrs)...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, attrs map[string]interface{}) {
	l.logger.ErrorContext(ctx, msg, attrsToSlogArgs(attrs)...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, attrs map[string]interface{}) {
	l.logger.WarnContext(ctx, msg, attrsToSlogArgs(attrs)...)
}

func (l *slogLogger) Debug(ctx context.Context, msg string, attrs map[string]interface{}) {
	l.logger.DebugContext(ctx, msg, attrsToSlogArgs(attrs)...)
}

func attrsToSlogArgs(attrs map[string]interface{}) []any {
	if len(attrs) == 0 {
		return nil
	}
	args := make([]any, 0, len(attrs)*2)
	for k, v := range attrs {
		args = append(args, k, v)
	}
	return args
}

func initTracer(ctx context.Context, cfg config.ObservabilityConfig, res *resource.Resource) error {
	var exporter sdktrace.SpanExporter
	var err error

	switch cfg.TraceExporter {
	case "stdout":
		exporter, err = stdouttrace.New()
	case "otlp":
		exporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(cfg.Endpoint), otlptracehttp.WithInsecure())
	case "none":
		return nil
	default:
		return fmt.Errorf("unknown trace exporter: %s", cfg.TraceExporter)
	}
	if err != nil {
		return fmt.Errorf("create trace exporter: %w", err)
	}

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	return nil
}

func initMeter(ctx context.Context, cfg config.ObservabilityConfig, res *resource.Resource) error {
	var exporter sdkmetric.Exporter
	var err error

	switch cfg.TraceExporter {
	case "stdout":
		exporter, err = stdoutmetric.New()
	case "otlp":
		exporter, err = otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(cfg.Endpoint), otlpmetrichttp.WithInsecure())
	case "none":
		return nil
	default:
		return fmt.Errorf("unknown metric exporter: %s", cfg.TraceExporter)
	}
	if err != nil {
		return fmt.Errorf("create metric exporter: %w", err)
	}

	meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)
	return nil
}

func initLogger(ctx context.Context, cfg config.ObservabilityConfig, res *resource.Resource) error {
	var exporter sdklog.Exporter
	var err error

	switch cfg.TraceExporter {
	case "stdout":
		exporter, err = stdoutlog.New()
	case "otlp":
		exporter, err = otlploghttp.New(ctx, otlploghttp.WithEndpoint(cfg.Endpoint), otlploghttp.WithInsecure())
	case "none":
		return nil
	default:
		return fmt.Errorf("unknown log exporter: %s", cfg.TraceExporter)
	}
	if err != nil {
		return fmt.Errorf("create log exporter: %w", err)
	}

	logProvider = sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	return nil
}
