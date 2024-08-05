package vkpns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

// NewClient creates a new instance of the VKPNS client.
func NewClient(ctx context.Context, options ClientOptions) (*Client, error) {
	if options.ProjectID == "" {
		return nil, fmt.Errorf("%w: ProjectID", ErrNoData)
	}

	if options.ServiceToken == "" {
		return nil, fmt.Errorf("%w: ServiceToken", ErrNoData)
	}

	options.ServiceToken = "Bearer " + options.ServiceToken
	options.VKPNSEndpoint = fmt.Sprintf(defaultMessagingEndpoint, options.ProjectID)

	return &Client{
		options: options,
		client: &http.Client{
			Transport: &http2.Transport{
				DialTLS:         DialTLS,
				ReadIdleTimeout: ReadIdleTimeout,
			},
			Timeout: HTTPClientTimeout,
		},
	}, nil
}

// Send sends a Notification to the VKPNs gateway. Context carries a
// deadline and a cancellation signal and allows you to close long running
// requests when the context timeout is exceeded. Context can be nil, for
// backwards compatibility.
//
// It will return a Response indicating whether the notification was accepted or
// rejected by the VKPNs gateway, or an error if something goes wrong.
func (c *Client) Send(ctx context.Context, message *Push) (*Response, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.options.VKPNSEndpoint,
		io.NopCloser(bytes.NewBuffer(body)),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", c.options.ServiceToken)

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(all, &result); err != nil {
		return nil, err
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("result:", result)
	fmt.Println("all:", string(all))

	return &result, nil
}

func (c *Client) SendDryRun(ctx context.Context, message *Push) (*Response, error) {
	return nil, ErrNotImplemented
}
