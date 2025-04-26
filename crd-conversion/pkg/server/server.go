package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/haoli/helm-crd-rollback/crd-conversion/pkg/conversion"
)

// WebhookServer handles the webhook requests
type WebhookServer struct {
	server *http.Server
}

// NewWebhookServer creates a new webhook server
func NewWebhookServer(port int, certFile, keyFile string) (*WebhookServer, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/convert", conversion.HandleConversion)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler:   mux,
	}

	return &WebhookServer{server: server}, nil
}

// Start starts the webhook server
func (ws *WebhookServer) Start() error {
	return ws.server.ListenAndServeTLS("", "")
}

// Shutdown gracefully shuts down the server
func (ws *WebhookServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ws.server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}
