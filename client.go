package bbapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	internalretry "github.com/raykavin/bbapi-go/internal/retry"
	gkhttp "github.com/raykavin/gokit/http"
)

const (
	// tokenExpiryBuffer is subtracted from the token TTL to avoid using a
	// token that expires mid-request.
	tokenExpiryBuffer = 30 * time.Second

	// gwAppKeyParam is the query parameter name required by the BB gateway.
	gwAppKeyParam = "gw-dev-app-key"
)

type tokenState struct {
	accessToken string
	expiresAt   time.Time
}

func (t tokenState) isValid(now time.Time) bool {
	if t.accessToken == "" {
		return false
	}
	return t.expiresAt.IsZero() || now.Before(t.expiresAt)
}

// Client is the Banco do Brasil SDK client. It is safe for concurrent use.
type Client struct {
	config      Config
	httpClient  *http.Client
	mtlsEnabled bool
	sandboxMode bool
	tokenMu     sync.RWMutex
	token       tokenState
}

// NewClient creates and validates a new Banco do Brasil API client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.ClientID == "" {
		return nil, errors.New("bbapi: ClientID is required")
	}
	if cfg.ClientSecret == "" {
		return nil, errors.New("bbapi: ClientSecret is required")
	}
	if cfg.AppKey == "" {
		return nil, errors.New("bbapi: AppKey is required")
	}

	if err := cfg.setDefaults(); err != nil {
		return nil, err
	}

	c := &Client{
		config:      cfg,
		httpClient:  cfg.HTTPClient,
		mtlsEnabled: cfg.MTLSEnabled,
		sandboxMode: cfg.Sandbox,
	}
	if cfg.AccessToken != "" {
		c.token.accessToken = cfg.AccessToken
	}

	return c, nil
}

// SetAccessToken stores a raw access token with no expiry.
func (c *Client) SetAccessToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token = tokenState{accessToken: token}
}

// SetTokenResponse caches the token from an OAuth2 response.
func (c *Client) SetTokenResponse(resp *TokenResponse) {
	if resp == nil {
		return
	}

	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	var expiry time.Time
	if resp.ExpiresIn > 0 {
		ttl := time.Duration(resp.ExpiresIn)*time.Second - tokenExpiryBuffer
		if ttl < 0 {
			ttl = 0
		}
		expiry = time.Now().Add(ttl)
	}

	c.token = tokenState{
		accessToken: resp.AccessToken,
		expiresAt:   expiry,
	}
}

// GetAccessToken returns the currently cached access token string.
func (c *Client) GetAccessToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token.accessToken
}

// TokenExpiresAt returns the current access-token expiration time.
func (c *Client) TokenExpiresAt() time.Time {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token.expiresAt
}

// MTLSErr returns ErrMTLSRequired if the client was not configured with
// mutual TLS and not in sandbox mode. Call this at the top of any method
// that the BB API mandates a client certificate for.
func (c *Client) MTLSErr(apiName string) error {
	if !c.mtlsEnabled && !c.sandboxMode {
		return fmt.Errorf("%s api require a mTLS authentication certificates", apiName)
	}
	return nil
}

func (c *Client) apiURL(path string) string {
	base := strings.TrimRight(c.config.APIURL, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

func (c *Client) authURL() string {
	return c.config.AuthURL
}

func (c *Client) ensureToken(ctx context.Context) (string, error) {
	c.tokenMu.RLock()
	state := c.token
	c.tokenMu.RUnlock()

	if state.isValid(time.Now()) {
		return state.accessToken, nil
	}

	tr, err := c.Authenticate(ctx)
	if err != nil {
		return "", err
	}
	return tr.AccessToken, nil
}

func (c *Client) do(
	ctx context.Context,
	method, rawURL string,
	payload []byte,
	contentType string,
) ([]byte, int, error) {
	var (
		responseBody []byte
		statusCode   int
	)

	reauthDone := false

	callErr := internalretry.Do(
		ctx,
		c.config.MaxRetries+1,
		c.config.RetryWaitMin,
		c.config.RetryWaitMax,
		func(_ int, err error) bool {
			if err == nil || ctx.Err() != nil {
				return false
			}
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				return true
			}
			if apiErr.StatusCode == http.StatusUnauthorized && !reauthDone {
				return true
			}
			return isRetryableStatus(apiErr.StatusCode)
		},
		func() error {
			token, err := c.ensureToken(ctx)
			if err != nil {
				return err
			}

			headers := gkhttp.DefaultJSONHeaders()
			headers.Set(gkhttp.HeaderAuthorization, "Bearer "+token)
			headers.Set(gkhttp.HeaderUserAgent, userAgent)
			if contentType != "" {
				headers.Set(gkhttp.HeaderContentType, contentType)
			}

			q := gkhttp.MapParams{}
			q.Set(gwAppKeyParam, c.config.AppKey)

			rb, sc, err := gkhttp.NewRequestWithContext(
				ctx,
				method,
				rawURL,
				q,
				headers,
				payload,
				c.httpClient,
			)
			if err != nil {
				return fmt.Errorf("bbapi: http: %w", err)
			}

			statusCode = sc

			if sc == http.StatusUnauthorized && !reauthDone {
				reauthDone = true
				c.tokenMu.Lock()
				c.token = tokenState{}
				c.tokenMu.Unlock()
				return parseAPIError(sc, rb)
			}

			if sc >= 400 {
				return parseAPIError(sc, rb)
			}

			responseBody = rb
			return nil
		},
	)

	return responseBody, statusCode, callErr
}

func (c *Client) doJSON(
	ctx context.Context,
	method, rawURL string,
	body any,
) ([]byte, error) {
	var (
		payload     []byte
		contentType string
	)

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("bbapi: marshal body: %w", err)
		}
		payload = b
		contentType = gkhttp.MIMEApplicationJSON
	}

	result, _, err := c.do(ctx, method, rawURL, payload, contentType)
	return result, err
}

func (c *Client) doFormDirect(
	ctx context.Context,
	targetURL string,
	payload []byte,
	extraHeaders map[string]string,
) ([]byte, int, error) {
	headers := gkhttp.DefaultFormHeaders()
	headers.Set(gkhttp.HeaderUserAgent, userAgent)
	for k, v := range extraHeaders {
		headers.Set(k, v)
	}

	rb, sc, err := gkhttp.NewRequestWithContext(
		ctx,
		http.MethodPost,
		targetURL,
		nil,
		headers,
		payload,
		c.httpClient,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("bbapi: auth http: %w", err)
	}
	return rb, sc, nil
}
