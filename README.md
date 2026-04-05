# bbapi-go

`bbapi-go` is a Go SDK for Banco do Brasil's Batch Payments API (`Pagamentos em Lote`).

This repository currently focuses on the API described in `OpenAPI_BB_Pagamentos em Lote_v1.json`, while keeping the package layout ready for future Banco do Brasil APIs.

## Features

- Idiomatic Go client API
- English code identifiers with JSON tags mapped to Banco do Brasil field names
- OAuth2 `client_credentials` authentication
- Automatic token reuse and refresh
- Retry support for transient HTTP failures
- Typed request and response models
- Unit tests for request building, model serialization, response parsing, and error handling

## Installation

```bash
go get github.com/raykavin/bbapi-go
```

## Requirements

You will need Banco do Brasil API credentials and configuration:

- `client_id`
- `client_secret`
- `app_key`
- OAuth scopes for the endpoints you want to call

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

	log.Printf("state=%d transfers=%d", resp.RequestState, resp.TransferCount)
}
```

## Configuration

The client is configured through `bbapi.Config`.

| Field | Description |
|---|---|
| `ClientID` | Banco do Brasil application client ID |
| `ClientSecret` | Banco do Brasil application client secret |
| `AppKey` | Application key sent to the API |
| `Sandbox` | Uses BB sandbox URLs when `true` |
| `APIURL` | Optional API base URL override |
| `AuthURL` | Optional OAuth token URL override |
| `AccessToken` | Optional initial token |
| `Scopes` | OAuth scopes requested during authentication |
| `HTTPClient` | Optional custom `*http.Client` |
| `Timeout` | HTTP timeout |
| `MaxRetries` | Retry count for transient failures |
| `RetryWaitMin` | Minimum retry backoff |
| `RetryWaitMax` | Maximum retry backoff |

### Default URLs

- Sandbox API: `https://homologa-api-ip.bb.com.br:7144/pagamentos-lote/v1`
- Production API: `https://api-ip.bb.com.br/pagamentos-lote/v1`
- Sandbox OAuth: `https://oauth.sandbox.bb.com.br/oauth/token`
- Production OAuth: `https://oauth.bb.com.br/oauth/token`

## Authentication

The SDK supports OAuth2 using the `client_credentials` flow.

You can authenticate explicitly:

```go
token, err := client.Authenticate(ctx)
if err != nil {
	log.Fatal(err)
}

log.Println(token.AccessToken)
```

Or set a token manually:

```go
client.SetAccessToken("your-access-token")
```

If no valid token is cached, the client authenticates automatically before calling the API.

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

See `doc.go` for the full scope list.

## Implemented API Coverage

The OpenAPI file in this repository contains 30 unique HTTP operations, and this SDK implements all 30 of them.

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

## Implemented Resource Prefixes

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

## Useful Constants

### Payment Types

- `bbapi.PaymentTypeSuppliers`
- `bbapi.PaymentTypeSalary`
- `bbapi.PaymentTypeMiscellaneous`

### Request States

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

## Error Handling

Errors returned by the API and OAuth server are normalized into `*bbapi.APIError`.

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

## Retry Behavior

The client retries transient failures, including:

- `429 Too Many Requests`
- `500 Internal Server Error`
- `502 Bad Gateway`
- `503 Service Unavailable`
- `504 Gateway Timeout`

The client also clears the cached token and re-authenticates after a `401 Unauthorized` response when appropriate.

## Testing

Run the test suite with:

```bash
go test ./...
```

The current tests cover:

- request construction
- model serialization and deserialization
- response parsing
- error handling

## Project Layout

| File | Responsibility |
|---|---|
| `client.go` | Core client, request flow, token reuse, retry integration |
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

## Notes

- All Go identifiers are intentionally written in English.
- Banco do Brasil JSON field names remain mapped through `json` tags.
- The codebase is organized so new Banco do Brasil APIs can be added later without reshaping the current client design.
