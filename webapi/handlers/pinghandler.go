package handlers

import (
	"fmt"
	"html"
	"net/http"
)

type PingHandler struct {
}

func (handler PingHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	return nil
}
