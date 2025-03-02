package main

import (
	"context"
	"fmt"
	"os"

	"cryptobotmanager.com/cbm-backend/microservices/utils/functions"
	"cryptobotmanager.com/cbm-backend/shared/messaging"
	"github.com/nats-io/nats.go"
)

func main() {
	ctx := context.Background()

	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = "nats://nats:4222"
	}

	nc, err := nats.Connect(natsUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to NATS server with: %v Err %v\n", natsUrl, err)
		return
	}
	messaging.SetupSignalHandlers(nc)

	// Loop through all registered services
	for name, startFunc := range functions.GetServices() {
		fmt.Printf("Starting service: %s\n", name)
		if err := startFunc(nc); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting service %s: %v\n", name, err)

			return
		}
	}

	<-ctx.Done()

}
