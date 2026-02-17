package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"

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

// providers holds the initialized OTel providers.
type providers struct {
	tracer *sdktrace.TracerProvider
	meter  *sdkmetric.MeterProvider
	logger *sdklog.LoggerProvider
}

type getTracerReq struct {
	name string
	resp chan trace.Tracer
}

type getMeterReq struct {
	name string
	resp chan metric.Meter
}

type getLoggerReq struct {
	name string
	resp chan Logger
}

type shutdownReq struct {
	ctx  context.Context
	resp chan error
}

var (
	tracerCh   = make(chan getTracerReq)
	meterCh    = make(chan getMeterReq)
	loggerCh   = make(chan getLoggerReq)
	shutdownCh = make(chan shutdownReq)
)

// Initialize sets up the global OTel providers and starts the manager goroutine.
func Initialize(ctx context.Context, cfg config.ObservabilityConfig) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}

	p := &providers{}

	if cfg.EnableTraces {
		if err := initTracer(ctx, cfg, res, p); err != nil {
			return fmt.Errorf("init tracer: %w", err)
		}
	}

	if cfg.EnableMetrics {
		if err := initMeter(ctx, cfg, res, p); err != nil {
			return fmt.Errorf("init meter: %w", err)
		}
	}

	if cfg.EnableLogs {
		if err := initLogger(ctx, cfg, res, p); err != nil {
			return fmt.Errorf("init logger: %w", err)
		}
	}

	go manage(p)
	return nil
}

// manage owns all provider state. All access goes through channels.
func manage(p *providers) {
	for {
		select {
		case req := <-tracerCh:
			if p.tracer != nil {
				req.resp <- p.tracer.Tracer(req.name)
			} else {
				req.resp <- otel.Tracer(req.name)
			}

		case req := <-meterCh:
			if p.meter != nil {
				req.resp <- p.meter.Meter(req.name)
			} else {
				req.resp <- otel.Meter(req.name)
			}

		case req := <-loggerCh:
			var handler slog.Handler
			if p.logger != nil {
				handler = otelslog.NewHandler(req.name, otelslog.WithLoggerProvider(p.logger))
			} else {
				handler = slog.NewJSONHandler(os.Stdout, nil)
			}
			req.resp <- &slogLogger{
				logger: slog.New(handler).With("component", req.name),
			}

		case req := <-shutdownCh:
			var errs []error
			if p.tracer != nil {
				if err := p.tracer.Shutdown(req.ctx); err != nil {
					errs = append(errs, fmt.Errorf("shutdown tracer: %w", err))
				}
				p.tracer = nil
			}
			if p.meter != nil {
				if err := p.meter.Shutdown(req.ctx); err != nil {
					errs = append(errs, fmt.Errorf("shutdown meter: %w", err))
				}
				p.meter = nil
			}
			if p.logger != nil {
				if err := p.logger.Shutdown(req.ctx); err != nil {
					errs = append(errs, fmt.Errorf("shutdown logger: %w", err))
				}
				p.logger = nil
			}
			if len(errs) > 0 {
				req.resp <- fmt.Errorf("shutdown errors: %v", errs)
			} else {
				req.resp <- nil
			}
			return
		}
	}
}

// Shutdown gracefully shuts down all providers.
func Shutdown(ctx context.Context) error {
	resp := make(chan error, 1)
	shutdownCh <- shutdownReq{ctx: ctx, resp: resp}
	return <-resp
}

// GetTracer returns a named tracer from the global provider.
func GetTracer(name string) trace.Tracer {
	resp := make(chan trace.Tracer, 1)
	tracerCh <- getTracerReq{name: name, resp: resp}
	return <-resp
}

// GetMeter returns a named meter from the global provider.
func GetMeter(name string) metric.Meter {
	resp := make(chan metric.Meter, 1)
	meterCh <- getMeterReq{name: name, resp: resp}
	return <-resp
}

// GetLogger returns a Logger backed by OTel if available, otherwise slog.
func GetLogger(name string) Logger {
	resp := make(chan Logger, 1)
	loggerCh <- getLoggerReq{name: name, resp: resp}
	return <-resp
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

func initTracer(ctx context.Context, cfg config.ObservabilityConfig, res *resource.Resource, p *providers) error {
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

	p.tracer = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(p.tracer)
	return nil
}

func initMeter(ctx context.Context, cfg config.ObservabilityConfig, res *resource.Resource, p *providers) error {
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

	p.meter = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(p.meter)
	return nil
}

func initLogger(ctx context.Context, cfg config.ObservabilityConfig, res *resource.Resource, p *providers) error {
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

	p.logger = sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
		sdklog.WithResource(res),
	)
	return nil
}
