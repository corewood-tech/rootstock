package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"connectrpc.com/connect"

	"rootstock/web-server/config"
	connecthandlers "rootstock/web-server/handlers/connect"
	"rootstock/web-server/proto/rootstock/v1/rootstockv1connect"
	"rootstock/web-server/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load("config.yaml", nil)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	interceptors := connect.WithInterceptors(server.BinaryOnlyInterceptor())

	healthHandler := connecthandlers.NewHealthServiceHandler()
	path, handler := rootstockv1connect.NewHealthServiceHandler(healthHandler, interceptors)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	httpServer := &http.Server{Handler: mux}

	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("server listening on %s\n", addr)
		if err := httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("shutting down server...")
		return httpServer.Close()
	case err := <-errChan:
		return err
	}
}
