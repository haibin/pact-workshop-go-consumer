package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/haibin/pact-workshop-go-consumer/model"
)

// Client is our consumer interface to the Order API
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

// GetUser gets a single user from the API
func (c *Client) GetUser(id int) (*model.User, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/user/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var user model.User
	_, err = c.do(req, &user)

	if err != nil {
		return nil, ErrUnavailable
	}

	fmt.Printf("***** response user: %+v\n", user)

	return &user, err
}

// GetUsers gets all users from the API
func (c *Client) GetUsers() ([]model.User, error) {
	req, err := c.newRequest("GET", "/users", nil)
	if err != nil {
		return nil, err
	}
	var users []model.User
	_, err = c.do(req, &users)

	return users, err
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Admin Service")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("response: %s\n", body)
	fmt.Printf("response header Content-Type: %s\n", resp.Header["Content-Type"])
	fmt.Printf("response header X-Api-Correlation-Id: %s\n", resp.Header["X-Api-Correlation-Id"])

	err = json.Unmarshal(body, &v)

	return resp, err
}

var (
	ErrUnavailable = errors.New("api unavailable")
)
