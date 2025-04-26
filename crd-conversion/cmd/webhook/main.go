package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/haoli/helm-crd-rollback/crd-conversion/pkg/server"
)

func main() {
	var tlsCertFile string
	var tlsKeyFile string
	var port int

	flag.StringVar(&tlsCertFile, "tls-cert-file", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS")
	flag.StringVar(&tlsKeyFile, "tls-key-file", "/etc/webhook/certs/tls.key", "File containing the x509 private key for HTTPS")
	flag.IntVar(&port, "port", 8443, "The port to listen on for HTTPS requests")
	flag.Parse()

	// Create a new webhook server
	webhookServer, err := server.NewWebhookServer(port, tlsCertFile, tlsKeyFile)
	if err != nil {
		log.Fatalf("Failed to create webhook server: %v", err)
	}

	// Start the server in a goroutine
	go func() {
		if err := webhookServer.Start(); err != nil {
			log.Fatalf("Failed to start webhook server: %v", err)
		}
	}()

	log.Printf("Conversion webhook server started on port %d with HTTPS", port)

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Received termination signal, shutting down gracefully...")
	webhookServer.Shutdown()
}
