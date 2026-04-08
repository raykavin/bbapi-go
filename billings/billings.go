package billings

import (
	"fmt"

	"github.com/raykavin/bbapi-go"
)

// Client wraps bbapi.Client with billing operations.
type Client struct {
	*bbapi.Client
}

// NewClient returns a billing client backed by the provided bbapi.Client.
func NewClient(bbClient *bbapi.Client) (*Client, error) {
	if bbClient == nil {
		return nil, fmt.Errorf("bb client is required for billings initialization")
	}

	return &Client{Client: bbClient}, nil
}
