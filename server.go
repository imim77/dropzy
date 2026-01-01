package main

import (
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Port string
}

type Logger struct {
	out io.Writer
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{out: out}
}

func (l *Logger) InfoMess(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func (l *Logger) Error(msg string, err error) {
	fmt.Fprintf(l.out, "[ERROR] %s: %v\n", msg, err)
}

func NewServer(cfg *Config, logger *Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, logger, *cfg)
	var handler http.Handler = mux
	return handler
}
