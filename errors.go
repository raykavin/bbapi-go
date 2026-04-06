package bbapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ErrorDetail represents a single error entry from the Banco do Brasil API.
type ErrorDetail struct {
	Codigo     string `json:"codigo"`
	Versao     string `json:"versao"`
	Mensagem   string `json:"mensagem"`
	Ocorrencia string `json:"ocorrencia"`
}

// ErrorResponse is the top-level BB API error body.
type ErrorResponse struct {
	Erros []ErrorDetail `json:"erros"`
}

// OAuthError represents an error returned by the OAuth2 server.
type OAuthError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
	Attributes struct {
		Error string `json:"error"`
	} `json:"attributes"`
}

// APIError represents an error returned by the BB API or OAuth server.
type APIError struct {
	StatusCode  int
	Message     string
	RawResponse string
	Details     []ErrorDetail
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if len(e.Details) == 0 {
		msg := e.RawResponse
		if e.Message != "" {
			msg = e.Message
		}
		return fmt.Sprintf("bbapi: HTTP %d: %s", e.StatusCode, msg)
	}

	parts := make([]string, 0, len(e.Details))
	for _, d := range e.Details {
		s := d.Mensagem
		if d.Codigo != "" {
			s = fmt.Sprintf("[%s] %s", d.Codigo, s)
		}
		if d.Ocorrencia != "" {
			s += " (" + d.Ocorrencia + ")"
		}
		parts = append(parts, s)
	}
	return fmt.Sprintf("bbapi: HTTP %d: %s", e.StatusCode, strings.Join(parts, "; "))
}

// asAPIError unwraps err into an *APIError, returning it and true if successful.
func asAPIError(err error) (*APIError, bool) {
	var e *APIError
	return e, errors.As(err, &e)
}

// IsNotFound returns true if err is an *APIError with HTTP 404.
func IsNotFound(err error) bool {
	e, ok := asAPIError(err)
	return ok && e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if err is an *APIError with HTTP 401.
func IsUnauthorized(err error) bool {
	e, ok := asAPIError(err)
	return ok && e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if err is an *APIError with HTTP 403.
func IsForbidden(err error) bool {
	e, ok := asAPIError(err)
	return ok && e.StatusCode == http.StatusForbidden
}

// IsRateLimited returns true if err is an *APIError with HTTP 429.
func IsRateLimited(err error) bool {
	e, ok := asAPIError(err)
	return ok && e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true if err is an *APIError with a 5xx status code.
func IsServerError(err error) bool {
	e, ok := asAPIError(err)
	return ok && e.StatusCode >= 500
}

// isRetryableStatus reports whether the status code warrants a retry.
func isRetryableStatus(code int) bool {
	switch code {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// parseAPIError builds an *APIError from a status code and raw body bytes.
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode:  statusCode,
		RawResponse: string(body),
	}

	var errResp ErrorResponse
	if json.Unmarshal(body, &errResp) == nil && len(errResp.Erros) > 0 {
		apiErr.Details = errResp.Erros
		apiErr.Message = errResp.Erros[0].Mensagem
		return apiErr
	}

	var oauthErr OAuthError
	if json.Unmarshal(body, &oauthErr) == nil && oauthErr.Error != "" {
		apiErr.Message = oauthErr.Message
		if apiErr.Message == "" {
			apiErr.Message = oauthErr.Error
		}
		return apiErr
	}

	return apiErr
}
