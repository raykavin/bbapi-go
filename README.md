# Unofficial Banco do Brasil API

[![Go Reference](https://pkg.go.dev/badge/github.com/raykavin/bbapi-go.svg)](https://pkg.go.dev/github.com/raykavin/bbapi-go)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/raykavin/bbapi-go)](https://goreportcard.com/report/github.com/raykavin/bbapi-go)

A Go library for the **Banco do Brasil Batch Payments API** (`Pagamentos em Lote`).

Covers all documented resources from the OpenAPI file in this repository: Payment Management, Transfers, Pix Transfers, Bank Slips, Barcode Guides, DARF, GPS, and GRU.

---

## Installation

```bash
go get github.com/raykavin/bbapi-go
```

Requires **Go 1.25+**.

The SDK uses `github.com/raykavin/gokit/http` for HTTP requests and includes retry support for transient failures.

---

## Quick Start

```go
package main

import (
	"context"
	"log"

	bbapi "github.com/raykavin/bbapi-go"
)

func main() {
	client, err := bbapi.NewClient(bbapi.Config{
		ClientID:     "your-client-id",
		ClientSecret: "your-client-secret",
		AppKey:       "your-app-key",
		Sandbox:      true,
		Scopes: []bbapi.Scope{
			bbapi.ScopeTransfersRequest,
			bbapi.ScopeBatchesRequest,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	tokenResp, err := client.Authenticate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	client.SetTokenResponse(tokenResp)

	resp, err := client.CreateTransferBatch(ctx, &bbapi.CreateTransferBatchRequest{
		RequestNumber: 123,
		PaymentType:   bbapi.PaymentTypeMiscellaneous,
		Transfers: []bbapi.Transfer{
			{
				TransferDate:  10042026,
				TransferValue: 150.75,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("request_state=%d transfer_count=%d", resp.RequestState, resp.TransferCount)
}
```

---

## Configuration

```go
client, err := bbapi.NewClient(bbapi.Config{
	// Required
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	AppKey:       "your-app-key",

	// Optional defaults shown
	Sandbox:      true,
	APIURL:       "",
	AuthURL:      "",
	AccessToken:  "",
	Scopes:       []bbapi.Scope{},
	HTTPClient:   &http.Client{},
	Timeout:      30 * time.Second,
	MaxRetries:   3,
	RetryWaitMin: 1 * time.Second,
	RetryWaitMax: 10 * time.Second,
})
```

### Default URLs

- Sandbox API: `https://homologa-api-ip.bb.com.br:7144/pagamentos-lote/v1`
- Production API: `https://api-ip.bb.com.br/pagamentos-lote/v1`
- Sandbox OAuth: `https://oauth.sandbox.bb.com.br/oauth/token`
- Production OAuth: `https://oauth.bb.com.br/oauth/token`

### Config fields

| Field | Description |
|---|---|
| `ClientID` | Banco do Brasil application client ID |
| `ClientSecret` | Banco do Brasil application client secret |
| `AppKey` | Application key required by the API |
| `Sandbox` | Switches the client to BB sandbox endpoints |
| `APIURL` | Optional API base URL override |
| `AuthURL` | Optional OAuth token URL override |
| `AccessToken` | Optional preloaded access token |
| `Scopes` | OAuth scopes requested during authentication |
| `HTTPClient` | Optional custom `*http.Client` |
| `Timeout` | Request timeout |
| `MaxRetries` | Retry attempts for transient failures |
| `RetryWaitMin` | Minimum retry backoff |
| `RetryWaitMax` | Maximum retry backoff |

---

## Authentication

The SDK supports OAuth2 using the `client_credentials` flow.

### Authenticate with configured credentials

```go
tokenResp, err := client.Authenticate(ctx)
if err != nil {
	log.Fatal(err)
}
client.SetTokenResponse(tokenResp)
```

### Authenticate with explicit credentials

```go
tokenResp, err := client.AuthenticateClientCredentials(ctx, bbapi.ClientCredentialsRequest{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	Scope:        "pagamentos-lote.transferencias-requisicao pagamentos-lote.lotes-requisicao",
})
if err != nil {
	log.Fatal(err)
}
client.SetTokenResponse(tokenResp)
```

### Set a token manually

```go
client.SetAccessToken("your-access-token")
```

### Token lifecycle

```go
token := client.GetAccessToken()
expiresAt := client.TokenExpiresAt()

_ = token
_ = expiresAt
```

If no valid token is cached, the client authenticates automatically before calling the API.

---

## OAuth Scopes

The package exposes typed scope constants, including:

- `bbapi.ScopeTransfersInfo`
- `bbapi.ScopeTransfersRequest`
- `bbapi.ScopePixInfo`
- `bbapi.ScopePixTransfersInfo`
- `bbapi.ScopePixTransfersRequest`
- `bbapi.ScopePaymentsInfo`
- `bbapi.ScopeBatchesInfo`
- `bbapi.ScopeBatchesRequest`
- `bbapi.ScopeBankSlipsInfo`
- `bbapi.ScopeBankSlipsRequest`
- `bbapi.ScopeBarcodeGuidesInfo`
- `bbapi.ScopeBarcodeGuidesRequest`
- `bbapi.ScopeReturnedPaymentsInfo`
- `bbapi.ScopeCancelRequest`

See [doc.go](/workspaces/app/doc.go) for the full list.

---

## Implemented API Coverage

The OpenAPI file in this repository defines **30 unique HTTP operations**, and this SDK implements all **30**.

### Payment Management

- `ReleasePayments`
- `CancelPayments`
- `UpdatePaymentDates`
- `ListReturnedPayments`
- `ListPaymentEntries`
- `GetBarcodePayments`

### Transfers

- `ListTransferBatches`
- `CreateTransferBatch`
- `GetBatch`
- `GetBatchRequest`
- `GetTransferPayment`
- `ListBeneficiaryTransfers`

### Pix Transfers

- `CreatePixTransferBatch`
- `GetPixTransferBatchRequest`
- `GetPixPayment`

### Bank Slips

- `CreateBankSlipBatch`
- `GetBankSlipBatchRequest`
- `GetBankSlipPayment`

### Barcode Guides

- `CreateBarcodeGuideBatch`
- `GetBarcodeGuideBatchRequest`
- `GetBarcodeGuidePayment`

### DARF

- `CreateDARFBatch`
- `GetDARFBatchRequest`
- `GetDARFPayment`

### GPS

- `CreateGPSBatch`
- `GetGPSBatchRequest`
- `GetGPSPayment`

### GRU

- `CreateGRUBatch`
- `GetGRUBatchRequest`
- `GetGRUPayment`

### Implemented resource prefixes

- `/pagamentos`
- `/liberar-pagamentos`
- `/cancelar-pagamentos`
- `/lancamentos-periodo`
- `/lotes-transferencias`
- `/transferencias`
- `/beneficiarios`
- `/lotes-transferencias-pix`
- `/pix`
- `/lotes-boletos`
- `/boletos`
- `/lotes-guias-codigo-barras`
- `/guias-codigo-barras`
- `/pagamentos-codigo-barras`
- `/lotes-darf-normal-preto`
- `/lotes-darf-preto-normal`
- `/darf-preto`
- `/lotes-gps`
- `/gps`
- `/pagamentos-gru`
- `/lotes-gru`
- `/gru`
- `/{id}`

---

## Useful Constants

### Payment types

- `bbapi.PaymentTypeSuppliers`
- `bbapi.PaymentTypeSalary`
- `bbapi.PaymentTypeMiscellaneous`

### Request states

- `bbapi.PaymentRequestStateConsistent`
- `bbapi.PaymentRequestStateInconsistent`
- `bbapi.PaymentRequestStateAllInconsistent`
- `bbapi.PaymentRequestStatePending`
- `bbapi.PaymentRequestStateProcessing`
- `bbapi.PaymentRequestStateProcessed`
- `bbapi.PaymentRequestStateRejected`
- `bbapi.PaymentRequestStatePreparingUnreleased`
- `bbapi.PaymentRequestStateReleasedByAPI`
- `bbapi.PaymentRequestStatePreparingReleased`

---

## Error Handling

Errors returned by Banco do Brasil's API and OAuth server are normalized into `*bbapi.APIError`.

```go
resp, err := client.GetGRUPayment(ctx, "123", nil)
if err != nil {
	var apiErr *bbapi.APIError
	if errors.As(err, &apiErr) {
		log.Printf("status=%d message=%s", apiErr.StatusCode, apiErr.Message)
	}
	log.Fatal(err)
}

_ = resp
```

Helper functions are available for common checks:

- `bbapi.IsNotFound(err)`
- `bbapi.IsUnauthorized(err)`
- `bbapi.IsForbidden(err)`
- `bbapi.IsRateLimited(err)`
- `bbapi.IsServerError(err)`

---

## Retry Behavior

The client retries transient failures, including:

- `429 Too Many Requests`
- `500 Internal Server Error`
- `502 Bad Gateway`
- `503 Service Unavailable`
- `504 Gateway Timeout`

The client also clears the cached token and re-authenticates after a `401 Unauthorized` response when appropriate.

---

## Testing

Run the test suite with:

```bash
go test ./...
```

Current tests cover:

- request construction
- model serialization and deserialization
- response parsing
- error handling

---

## Project Layout

| File | Responsibility |
|---|---|
| `client.go` | Core client, token reuse, request execution, retry integration |
| `config.go` | SDK configuration and defaults |
| `auth.go` | OAuth authentication |
| `errors.go` | API error parsing and helpers |
| `helpers.go` | Shared request and decoding helpers |
| `payments.go` | General payment management endpoints |
| `transfers.go` | Transfers and Pix transfers |
| `bank_slips.go` | Bank slip endpoints |
| `barcode_guides.go` | Barcode guide endpoints |
| `darf.go` | DARF endpoints |
| `gps.go` | GPS endpoints |
| `gru.go` | GRU endpoints |

---

## Notes

- All Go identifiers are intentionally written in English.
- Banco do Brasil JSON field names remain mapped through `json` tags.
- The codebase is organized so new Banco do Brasil APIs can be added later without changing the current client structure.
