package handler

import (
	"context"
	"net/http"
	"restgrpcproxy/internal/lib/logger/sl"

	gw "restgrpcproxy/gen/go/sso"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// collects incoming request and sends it to the output
func HandlerReturn(w http.ResponseWriter, r *http.Request, gRPCServer string, log *slog.Logger) {
	op := "restgrpcproxy.handler.HandlerReturn"

	ctx := context.Background()
	mux := runtime.NewServeMux()

	/*
		creds, err := credentials.NewClientTLSFromFile("path/to/server-cert.pem", "")
		if err != nil {
			log.Error("failed to load TLS credentials", sl.Err(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
			// without TLS:
			// grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	*/
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := gw.RegisterAuthHandlerFromEndpoint(ctx, mux, gRPCServer, opts)
	if err != nil {
		log.Error(op, "failed to register gateway", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Обслуживаем запрос через gateway
	mux.ServeHTTP(w, r)

}
