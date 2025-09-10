package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"restgrpcproxy/internal/lib/logger/sl"
	"strconv"
	"strings"

	gw "github.com/AlexBah/Protos/gen/go/sso"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// extractToken get token from different sourse
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		if strings.HasPrefix(authHeader, "Token ") {
			return strings.TrimPrefix(authHeader, "Token ")
		}
		return authHeader
	}

	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}

	if cookie, err := r.Cookie("auth_token"); err == nil {
		return cookie.Value
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		r.ParseForm()
		if token := r.Form.Get("token"); token != "" {
			return token
		}
	}

	return ""
}

func customUpdateUserHandler(ctx context.Context, client gw.AuthClient, w http.ResponseWriter, r *http.Request, pathParams map[string]string, log *slog.Logger) {
	op := "restgrpcproxy.handler.customUpdateUserHandler"
	log = log.With(slog.String("op", op))

	token := extractToken(r)
	if token == "" {
		log.Error("missing authorization token")
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	userID := pathParams["user_id"]
	if userID == "" {
		log.Error("missing user_id in path")
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var updateReq gw.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	updateReq.Token = token
	temp, err := strconv.Atoi(userID)
	if err != nil {
		log.Error("userID is not integer", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	updateReq.UserId = int64(temp)

	resp, err := client.UpdateUser(ctx, &updateReq)
	if err != nil {
		log.Error("failed to update user via gRPC", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", sl.Err(err))
	}
}

func customDeleteUserHandler(ctx context.Context, client gw.AuthClient, w http.ResponseWriter, r *http.Request, pathParams map[string]string, log *slog.Logger) {
	op := "restgrpcproxy.handler.customDeleteUserHandler"
	log = log.With(slog.String("op", op))

	token := extractToken(r)
	if token == "" {
		log.Error("missing authorization token")
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	phone := pathParams["phone"]
	if phone == "" {
		log.Error("missing phone in path")
		http.Error(w, "Phone required", http.StatusBadRequest)
		return
	}

	deleteReq := &gw.DeleteUserRequest{
		Phone: phone,
		Token: token,
	}

	resp, err := client.DeleteUser(ctx, deleteReq)
	if err != nil {
		log.Error("failed to delete user via gRPC", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error("failed to encode response", sl.Err(err))
	}
}

// HandlerReturn start grpc connection with grpc-gateway
func HandlerReturn(w http.ResponseWriter, r *http.Request, gRPCServer string, log *slog.Logger) {
	op := "restgrpcproxy.handler.HandlerReturn"
	log = log.With(slog.String("op", op))

	log.Debug("request", r.Method, r.URL.Path)

	ctx := context.Background()
	mux := runtime.NewServeMux()

	var opts []grpc.DialOption
	if strings.Contains(gRPCServer, "127.0.0.1") {
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	} else {
		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})),
		}
	}

	clientConn, err := grpc.NewClient(gRPCServer, opts...)
	if err != nil {
		log.Error("failed to create gRPC client", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()
	client := gw.NewAuthClient(clientConn)

	err = gw.RegisterAuthHandlerFromEndpoint(ctx, mux, gRPCServer, opts)
	if err != nil {
		log.Error("failed to register gateway", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = mux.HandlePath("PUT", "/v1/users/{user_id}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		customUpdateUserHandler(ctx, client, w, r, pathParams, log)
	})
	if err != nil {
		log.Error("failed to register update handler", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = mux.HandlePath("DELETE", "/v1/users/{phone}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		customDeleteUserHandler(ctx, client, w, r, pathParams, log)
	})
	if err != nil {
		log.Error("failed to register delete handler", sl.Err(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Debug("serving request", r.URL.Path)
	mux.ServeHTTP(w, r)
}
