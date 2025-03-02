package functions

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

const (
	ServiceName     = "cbm-echo-test3"
	ServiceSubject1 = "cbm.echo.test3"
	VersionNumber   = "0.0.1"
)

// StartEchoService initializes and starts the echo service
func StartService(nc *nats.Conn) error {
	// Request handler
	handler := func(req micro.Request) {
		// Check if the request has a reply subject
		if req.Reply() == "" {
			fmt.Fprintf(os.Stdout, "Received message without reply subject: %s\n", string(req.Data()))
			return // No reply-to subject, so no need to respond
		}

		// Handle request and send response
		err := req.Respond(req.Data())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
		}
	}

	// Add the service to NATS
	_, err := micro.AddService(nc, micro.Config{
		Name:    ServiceName,
		Version: VersionNumber,
		// Base handler
		Endpoint: &micro.EndpointConfig{
			Subject: ServiceSubject1,
			Handler: micro.HandlerFunc(handler),
		},
	})

	if err != nil {
		return fmt.Errorf("error adding service: %w", err)
	}

	return nil
}
