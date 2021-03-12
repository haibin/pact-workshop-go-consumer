package main

import (
	"log"
	"net/url"

	"github.com/haibin/pact-workshop-go-consumer/consumer/client"
)

func main() {
	u, _ := url.Parse("http://localhost:8080")
	client := &client.Client{
		BaseURL: u,
	}

	users, err := client.GetUser(10)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(users)
}
