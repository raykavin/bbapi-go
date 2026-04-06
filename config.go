package bbapi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
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
//
// mTLS (mutual TLS) required for Banco do Brasil production endpoints:
// Provide the client certificate and key either as file paths (MTLSCertFile /
// MTLSKeyFile) or as PEM-encoded bytes (MTLSCertPEM / MTLSKeyPEM). File paths
// take precedence when both are set. Optionally supply a custom CA root via
// MTLSCARootFile or MTLSCARootPEM. When any mTLS field is set and HTTPClient
// is nil, the SDK builds an http.Client with the appropriate tls.Config
// automatically.
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

	// mTLS client certificate file path (takes precedence over PEM bytes).
	MTLSCertFile string
	MTLSKeyFile  string

	// mTLS client certificate raw PEM-encoded bytes.
	MTLSCertPEM []byte
	MTLSKeyPEM  []byte

	// Optional custom CA root used to verify the server certificate.
	// File path takes precedence over PEM bytes.
	MTLSCARootFile string
	MTLSCARootPEM  []byte

	// MTLSEnabled marks the client as mTLS-capable. It is set automatically
	// when MTLSCertFile/MTLSKeyFile or MTLSCertPEM/MTLSKeyPEM are provided.
	// Set it to true manually when passing a pre-configured HTTPClient that
	// already presents a client certificate.
	MTLSEnabled bool
}

func (c *Config) setDefaults() error {
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
		if c.hasMTLS() {
			var err error
			c.HTTPClient, err = c.buildMTLSClient()
			if err != nil {
				return err
			}
			c.MTLSEnabled = true
		} else {
			c.HTTPClient = &http.Client{Timeout: c.Timeout}
		}
	}

	return nil
}

func (c *Config) hasMTLS() bool {
	return c.MTLSCertFile != "" || len(c.MTLSCertPEM) > 0
}

// buildMTLSClient constructs an *http.Client whose TLS transport presents the
// configured client certificate to the server (mutual TLS).
func (c *Config) buildMTLSClient() (*http.Client, error) {
	var (
		certPEM []byte
		keyPEM  []byte
		err     error
	)

	if c.MTLSCertFile != "" {
		certPEM, err = os.ReadFile(c.MTLSCertFile)
		if err != nil {
			return nil, fmt.Errorf("bbapi: mtls: read cert file: %w", err)
		}
		keyPEM, err = os.ReadFile(c.MTLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("bbapi: mtls: read key file: %w", err)
		}
	} else {
		certPEM = c.MTLSCertPEM
		keyPEM = c.MTLSKeyPEM
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("bbapi: mtls: parse key pair: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if c.MTLSCARootFile != "" || len(c.MTLSCARootPEM) > 0 {
		var caPEM []byte
		if c.MTLSCARootFile != "" {
			caPEM, err = os.ReadFile(c.MTLSCARootFile)
			if err != nil {
				return nil, fmt.Errorf("bbapi: mtls: read ca root file: %w", err)
			}
		} else {
			caPEM = c.MTLSCARootPEM
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caPEM) {
			return nil, fmt.Errorf("bbapi: mtls: failed to parse CA root certificate")
		}
		tlsCfg.RootCAs = pool
	}

	transport := &http.Transport{TLSClientConfig: tlsCfg}
	return &http.Client{
		Transport: transport,
		Timeout:   c.Timeout,
	}, nil
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
