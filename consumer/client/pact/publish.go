package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

func main() {
	// Publish the Pacts...
	p := dsl.Publisher{}

	fmt.Println("Publishing Pact files to broker", os.Getenv("PACT_DIR"), os.Getenv("PACT_BROKER_URL"))
	err := p.Publish(types.PublishRequest{
		PactURLs:        []string{filepath.FromSlash(fmt.Sprintf("%s/dineinconsumer-dineinprovider.json", os.Getenv("PACT_DIR")))},
		PactBroker:      "https://test.pact.dius.com.au",
		ConsumerVersion: "1.1.1",
		Tags:            []string{"master", "prod"},

		BrokerUsername:  os.Getenv("PACT_BROKER_USERNAME"),
		BrokerPassword:  os.Getenv("PACT_BROKER_PASSWORD"),
	})

	if err != nil {
		log.Println("ERROR: ", err)
		os.Exit(1)
	}
}
