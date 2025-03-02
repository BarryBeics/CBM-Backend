package functions

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

// StartService initializes and starts the service
func (u *Utils) ExampleService1(nc *nats.Conn) error {
	// Add the service to NATS
	_, err := micro.AddService(nc, micro.Config{
		Name:        "cbm-utils-example1",
		Version:     "0.0.1",
		Description: "Test 1",
		Endpoint: &micro.EndpointConfig{
			Subject: "cbm.utils.test1",
			Handler: micro.HandlerFunc(u.Handler1), // Use the separated handler function
		},
	})

	if err != nil {
		return fmt.Errorf("error adding service: %w", err)
	}

	return nil
}

// Handler processes incoming requests
func (u *Utils) Handler1(req micro.Request) {
	// Check if the request has a reply subject
	if req.Reply() == "" {
		fmt.Fprintf(os.Stdout, "Received message without reply subject: %s\n", string(req.Data()))
		return
	}

	// Generate the response
	response := u.BuildResponse1(req)

	// Convert response to JSON
	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling response: %v\n", err)
		return
	}

	// Send the response
	err = req.Respond(responseData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
	}
}

// BuildResponse constructs a response map
func (u *Utils) BuildResponse1(req micro.Request) map[string]interface{} {
	return map[string]interface{}{
		"originalMessage": string(req.Data()),
		"service":         "cbm-utils-test1",
		"additionalInfo":  "This is a test response from service test1",
	}
}
