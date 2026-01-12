package main

import (
	"context"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/smart-safety-hub/backend/internal/app"
)

func decodeKey(envVar string) string {
	encoded := os.Getenv(envVar)
	if encoded == "" {
		log.Fatalf("Environment variable %s is empty", envVar)
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Failed to decode %s: %v", envVar, err)
	}

	return string(decoded)
}

func main() {
	godotenv.Load()

	privateKey := decodeKey("PRIVATE_KEY_BASE64")
	publicKey := decodeKey("PUBLIC_KEY_BASE64")

	cfg := app.Config{
		GrpcAddr:   ":50051",
		HTTPAddr:   ":8080",
		DBURL:      "postgres://postgres:root@localhost:5432/smart_safety_hub?sslmode=disable",
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	container, close := app.Bootstrap(cfg)
	defer close()

	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen gRPC server")
	}

	go func() {
		container.Logger.Sugar().Infof("GRPC listening %s", cfg.GrpcAddr)
		if err := container.GRPCServer.Serve(lis); err != nil {
			container.Logger.Sugar().Fatalf("Failed to start gRPC server %v", err)
		}
	}()

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: container.HTTPRouter}

	go func() {
		container.Logger.Sugar().Infof("Rest listening %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Sugar().Fatalf("Failed to start http server %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	container.Logger.Sugar().Info("Shutdown signal receive")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	container.GRPCServer.GracefulStop()
	if err := srv.Shutdown(ctx); err != nil {
		container.Logger.Sugar().Error("http shutdown err %v", err)
	}
	container.Logger.Sugar().Info("Stopped")
}
