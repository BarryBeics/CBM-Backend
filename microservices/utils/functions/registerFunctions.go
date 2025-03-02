package functions

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

// ServiceHandler defines the interface for all service handlers.
type ServiceHandler interface {
	Handler(req micro.Request)                              // Handles incoming requests
	BuildResponse(req micro.Request) map[string]interface{} // Builds a response
}

type Utils struct{}

// Define a type for service registration functions
type ServiceFunc func(nc *nats.Conn) error

// Central registry for all service functions
var serviceRegistry = make(map[string]ServiceFunc)

// RegisterService adds a service to the registry
func RegisterService(name string, fn ServiceFunc) {
	serviceRegistry[name] = fn
}

// GetServices returns all registered services
func GetServices() map[string]ServiceFunc {
	return serviceRegistry
}
