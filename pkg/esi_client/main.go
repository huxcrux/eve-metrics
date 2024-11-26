package esiClient

import (
	"context"
	"net/http"

	"github.com/antihax/goesi"
)

type ESIClient struct {
	Ctx    context.Context
	Client *goesi.APIClient
}

func NewESIClient(token string) ESIClient {

	// Create a custom transport with the token
	customTransport := &customTransport{
		Transport: http.DefaultTransport,
		Token:     token,
	}

	// Create an HTTP client using the custom transport
	httpClient := &http.Client{
		Transport: customTransport,
	}

	// Use the custom client with goesi
	esiClient := goesi.NewAPIClient(httpClient, "eve-notifier")
	esiClient.ChangeBasePath("http://127.0.0.1:8182")

	ctx := context.TODO()

	return ESIClient{
		Client: esiClient,
		Ctx:    ctx,
	}
}

// CustomTransport adds a "token" parameter to every request.
type customTransport struct {
	Transport http.RoundTripper
	Token     string
}

func (ct *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to modify it safely
	newReq := req.Clone(context.Background())

	// Disable gzip compression
	newReq.Header.Set("Accept-Encoding", "identity")

	// Parse the URL and add the token parameter
	q := newReq.URL.Query()
	q.Add("token", ct.Token)
	newReq.URL.RawQuery = q.Encode()

	// Proceed with the original transport
	return ct.Transport.RoundTrip(newReq)
}
