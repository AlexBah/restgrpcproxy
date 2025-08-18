package handler

import (
	"fmt"
	"net/http"
	"restgrpcproxy/internal/lib/logger/sl"

	"golang.org/x/exp/slog"
)

// collects incoming request and sends it to the output
func HandlerReturn(w http.ResponseWriter, r *http.Request, gRPCServer string, log *slog.Logger) {
	op := "restgrpcproxy.handler.HandlerReturn"

	requestString := fmt.Sprintf("%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		requestString += fmt.Sprintf("%q : %q\n", k, v)
	}
	requestString += fmt.Sprintf("\nHost = %q\nRemoteAddr = %q\n\n", r.Host, r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Error(op, sl.Err(err))
	}
	for k, v := range r.Form {
		requestString += fmt.Sprintf("Form[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "%s", requestString)
}
