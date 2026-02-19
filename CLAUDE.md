# Rootstock — Engineering Rules

## Metalogic and behavior
- If you cannot provide an example or a reference, you are guessing. Stop and look for something verifiable to build on.

## Patterns and Consistency
- Consistency is one of the most important design principles of this project. Verify all changes against design principles.
- USE THE MAKEFILE for everything we need to do, do not do one-off podman commands. Most of the time you should not need to touch runtime, everything hot reloads.

## Go Concurrency
- Never use `sync.Mutex` or `sync.RWMutex`. Use goroutines and channels for all concurrency.

## Verification
- The web server hot reloads via `air`. Do not run `go build` or `go test` locally — just save files and verify against the running stack. STILL CREATE TESTS, but use the container to run time.
- Health check: `curl -s -X POST http://localhost:8080/rootstock.v1.HealthService/Check -H "Content-Type: application/proto" --data-binary ''`
- Check server logs: `podman logs --tail 30 compose_web-server_1`
- Query span metrics via Grafana→Prometheus proxy: `curl -s 'http://localhost:9999/grafana/api/datasources/proxy/uid/prometheus/api/v1/query?query=span_metrics_calls_total'`
- Grafana dashboards: `http://localhost:9999/grafana/`
- The server uses Connect RPC with a `BinaryOnlyInterceptor` — JSON encoding is rejected, use `application/proto`.
