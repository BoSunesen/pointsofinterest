package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type PingHandler struct {
}

func (handler PingHandler) ServeHttp(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "I am alive, server time: %v", time.Now())
	return nil
}
