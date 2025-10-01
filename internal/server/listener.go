package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"restgrpcproxy/internal/handler"
	"time"

	"golang.org/x/exp/slog"
)

// listens to a port, choosing between a secure or unsecured connection
func ListenPort(port, gRPCServer, tlsPath string, shutdownCh <-chan struct{}, log *slog.Logger, timeout time.Duration) {
	srv := &http.Server{Addr: port, Handler: http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) { handler.HandlerReturn(w, r, gRPCServer, log) },
	)}
	log.Info(fmt.Sprintf("Starting listen on port %s", srv.Addr))

	if tlsPath == "not exist" {
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Error("Port", srv.Addr, err)
			}
		}()
	} else {
		go func() {
			srv.TLSConfig = &tls.Config{NextProtos: []string{"h2", "http/1.1"}}
			certFile := tlsPath + "fullchain.pem"
			keyFile := tlsPath + "privkey.pem"
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
				log.Error("Port", srv.Addr, err)
			}
		}()
	}

	stopListen(srv, shutdownCh, log, timeout)
}

// stop listen port, then come signal close application
func stopListen(srv *http.Server, shutdownCh <-chan struct{}, log *slog.Logger, timeout time.Duration) {
	go func() {
		<-shutdownCh

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		log.Info(fmt.Sprintf("Shutting down server on port %s", srv.Addr))
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Server shutdown failed on port", srv.Addr, err)
		} else {
			log.Info(fmt.Sprintf("Server shutdown gracefully on port %s ", srv.Addr))
		}
	}()
}
