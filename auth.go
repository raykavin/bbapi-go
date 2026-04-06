package bbapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	gkhttp "github.com/raykavin/gokit/http"
)

const defaultClientCredentialsGrantType = "client_credentials"

// TokenResponse is the OAuth2 access token response from Banco do Brasil.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// ClientCredentialsRequest holds OAuth client credentials for Banco do Brasil.
type ClientCredentialsRequest struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	Scope        string
}

// Authenticate uses the credentials configured on the client and caches the
// access token for future API calls.
func (c *Client) Authenticate(ctx context.Context) (*TokenResponse, error) {
	return c.AuthenticateClientCredentials(ctx, ClientCredentialsRequest{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scope:        c.config.scopeString(),
	})
}

// AuthenticateClientCredentials performs the OAuth2 client_credentials flow.
func (c *Client) AuthenticateClientCredentials(
	ctx context.Context,
	req ClientCredentialsRequest,
) (*TokenResponse, error) {
	clientID := strings.TrimSpace(req.ClientID)
	if clientID == "" {
		clientID = c.config.ClientID
	}

	clientSecret := strings.TrimSpace(req.ClientSecret)
	if clientSecret == "" {
		clientSecret = c.config.ClientSecret
	}

	form := url.Values{}
	form.Set("grant_type", defaultString(req.GrantType, defaultClientCredentialsGrantType))
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	if scope := strings.TrimSpace(req.Scope); scope != "" {
		form.Set("scope", scope)
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	headers := map[string]string{gkhttp.HeaderAuthorization: "Basic " + credentials}

	raw, statusCode, err := c.doFormDirect(ctx, c.authURL(), []byte(form.Encode()), headers)
	if err != nil {
		return nil, fmt.Errorf("bbapi: authenticate: %w", err)
	}
	if statusCode >= 400 {
		return nil, parseAPIError(statusCode, raw)
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(raw, &tokenResponse); err != nil {
		return nil, fmt.Errorf("bbapi: authenticate: decode response: %w", err)
	}

	c.SetTokenResponse(&tokenResponse)
	return &tokenResponse, nil
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
