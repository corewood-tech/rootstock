module rootstock/web-server

go 1.25rc2

require (
	connectrpc.com/connect v1.19.1
	connectrpc.com/otelconnect v0.9.0
	github.com/dbos-inc/dbos-transact-golang v0.11.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/env v1.1.0
	github.com/knadh/koanf/providers/file v1.2.1
	github.com/knadh/koanf/providers/posflag v1.0.1
	github.com/knadh/koanf/providers/structs v1.0.0
	github.com/knadh/koanf/v2 v2.3.2
	github.com/spf13/pflag v1.0.10
	go.opentelemetry.io/contrib/bridges/otelslog v0.15.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.65.0
	go.opentelemetry.io/otel v1.40.0
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.16.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.40.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.40.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.16.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.40.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.40.0
	go.opentelemetry.io/otel/metric v1.40.0
	go.opentelemetry.io/otel/sdk v1.40.0
	go.opentelemetry.io/otel/sdk/log v0.16.0
	go.opentelemetry.io/otel/sdk/metric v1.40.0
	go.opentelemetry.io/otel/trace v1.40.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7 // indirect
	github.com/jackc/pgerrcode v0.0.0-20250907135507-afb5586c32a6 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.40.0 // indirect
	go.opentelemetry.io/otel/log v0.16.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260128011058-8636f8732409 // indirect
	google.golang.org/grpc v1.78.0 // indirect
)
