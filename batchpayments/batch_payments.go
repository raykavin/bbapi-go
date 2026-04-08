package batchpayments

import (
	"fmt"

	"github.com/raykavin/bbapi-go"
)

// Client wraps bbapi.Client with batch payment operations.
type Client struct {
	*bbapi.Client
}

// NewClient returns a Client backed by the provided bbapi.Client.
// The underlying client must be configured with mTLS.
func NewClient(bbClient *bbapi.Client) (*Client, error) {
	if bbClient == nil {
		return nil, fmt.Errorf("bb client is required for batch payments initialization")
	}

	if err := bbClient.MTLSErr("batch payments"); err != nil {
		return nil, err
	}

	return &Client{Client: bbClient}, nil
}
