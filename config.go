package bbapi

import (
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout      = 30 * time.Second
	defaultMaxRetries   = 3
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 10 * time.Second
	sdkVersion          = "0.1.0"
	userAgent           = "bbapi-go-sdk/" + sdkVersion
)

const (
	sandboxAPIURL     = "https://homologa-api-ip.bb.com.br:7144/pagamentos-lote/v1"
	productionAPIURL  = "https://api-ip.bb.com.br/pagamentos-lote/v1"
	sandboxAuthURL    = "https://oauth.sandbox.bb.com.br/oauth/token"
	productionAuthURL = "https://oauth.bb.com.br/oauth/token"
)

// Config holds all configuration for the BB API client.
//
// Zero values for duration and retry fields mean "use the default".
// To explicitly disable retries, set MaxRetries to -1.
type Config struct {
	ClientID     string
	ClientSecret string
	AppKey       string
	Sandbox      bool
	APIURL       string
	AuthURL      string
	AccessToken  string
	Scopes       []Scope
	HTTPClient   *http.Client
	Timeout      time.Duration
	// MaxRetries is the number of retry attempts after the first failure.
	// Set to -1 to disable retries entirely.
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

func (c *Config) setDefaults() {
	if c.Timeout <= 0 {
		c.Timeout = defaultTimeout
	}

	switch {
	case c.MaxRetries == 0:
		c.MaxRetries = defaultMaxRetries
	case c.MaxRetries < 0:
		c.MaxRetries = 0
	}

	if c.RetryWaitMin <= 0 {
		c.RetryWaitMin = defaultRetryWaitMin
	}

	if c.RetryWaitMax <= 0 {
		c.RetryWaitMax = defaultRetryWaitMax
	}

	if c.APIURL == "" {
		c.APIURL = pickURL(c.Sandbox, sandboxAPIURL, productionAPIURL)
	}

	if c.AuthURL == "" {
		c.AuthURL = pickURL(c.Sandbox, sandboxAuthURL, productionAuthURL)
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: c.Timeout}
	}
}

// pickURL returns sandboxVal when sandbox is true, otherwise productionVal.
func pickURL(sandbox bool, sandboxVal, productionVal string) string {
	if sandbox {
		return sandboxVal
	}
	return productionVal
}

func (c Config) scopeString() string {
	var sb strings.Builder
	for _, scope := range c.Scopes {
		if scope == "" {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(string(scope))
	}
	return sb.String()
}
