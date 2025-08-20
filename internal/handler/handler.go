package handler

import (
	"context"
	"crypto/tls"
	"net/http"
	"restgrpcproxy/internal/lib/logger/sl"

	gw "restgrpcproxy/gen/go/sso"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// start grpc connection with grpc-gateway
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
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	*/

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
	}

	err := gw.RegisterAuthHandlerFromEndpoint(ctx, mux, gRPCServer, opts)
	if err != nil {
		log.Error(op, "failed to register gateway", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	mux.ServeHTTP(w, r)

}
