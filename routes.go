package main

import (
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *Logger,
	config Config,
) {
	mux.Handle("/", http.NotFoundHandler())
}
