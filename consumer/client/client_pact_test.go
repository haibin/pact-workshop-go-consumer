// +build integration

package client

import (
	"fmt"
	"os"
	"testing"

	"net/url"

	"github.com/haibin/pact-workshop-go-consumer/model"
	"github.com/pact-foundation/pact-go/dsl"
)

var pact dsl.Pact

func TestMain(m *testing.M) {
	var exitCode int

	// Setup Pact and related test stuff
	pact = dsl.Pact{
		Consumer: os.Getenv("CONSUMER_NAME"),
		Provider: os.Getenv("PROVIDER_NAME"),
		LogDir:   os.Getenv("LOG_DIR"),
		PactDir:  os.Getenv("PACT_DIR"),
		LogLevel: "INFO",
	}
	defer pact.Teardown()

	// Proactively start service to get access to the port
	pact.Setup(true)

	// Run all the tests
	exitCode = m.Run()

	// Shutdown the Mock Service and Write pact files to disk
	if err := pact.WritePact(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func TestClientPact_GetUser(t *testing.T) {
	t.Run("the user exists", func(t *testing.T) {
		id := 10

		pact.
			AddInteraction().
			Given("User sally exists").
			UponReceiving("A request to login with user 'sally'").
			WithRequest(dsl.Request{
				Method: "GET",
				Path:   dsl.Term("/user/10", "/user/[0-9]+"),
			}).
			WillRespondWith(dsl.Response{
				Status: 200,
				Headers: dsl.MapMatcher{
					"Content-Type":         dsl.Term("application/json; charset=utf-8", `application\/json`),
					"X-Api-Correlation-Id": dsl.Like("100"),
				},
				Body: dsl.Match(model.User{}),
			})

		test := func() error {
			u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", pact.Server.Port))
			client := &Client{
				BaseURL: u,
			}

			user, err := client.GetUser(id)

			// Assert basic fact
			if user.ID != id {
				return fmt.Errorf("wanted user with ID %d but got %d", id, user.ID)
			}

			return err
		}

		if err := pact.Verify(test); err != nil {
			t.Fatalf("Error on Verify: %v", err)
		}
	})
}
