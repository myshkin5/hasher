package handlers

import (
	"net/http"

	"github.com/myshkin5/hasher/logs"
)

const ShutdownPattern = "/shutdown"

func NewShutdownFunc(shuttingDown chan struct{}) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		logs.Logger.Info("Shutdown started...")
		close(shuttingDown)
	}
}
