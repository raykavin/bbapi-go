# Unofficial Banco do Brasil API

[![Go Reference](https://pkg.go.dev/badge/github.com/raykavin/bbapi-go.svg)](https://pkg.go.dev/github.com/raykavin/bbapi-go)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/raykavin/bbapi-go)](https://goreportcard.com/report/github.com/raykavin/bbapi-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

An unofficial Go SDK for [Banco do Brasil](https://developers.bb.com.br) APIs.

`bbapi-go` provides a clean, type-safe client for interacting with BB's platform. It handles OAuth2 authentication, automatic token renewal, transient-error retries, and response parsing, so you can focus on your application logic rather than HTTP plumbing.

> **Disclaimer:** This project is an unofficial client library for the
> [Banco do Brasil Open Finance APIs](https://developers.bb.com.br).
> It is not affiliated with, endorsed by, or maintained by Banco do Brasil S.A.
> "BB" and "Banco do Brasil" are registered trademarks of Banco do Brasil S.A.

---

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Authentication](#authentication)
- [Mutual TLS (mTLS)](#mutual-tls-mtls)
- [OAuth Scopes](#oauth-scopes)
- [API Coverage](#api-coverage)
  - [Batch Payments: Transfers (TED/DOC)](#batch-payments--transfers-teddoc)
  - [Batch Payments: Pix Transfers](#batch-payments--pix-transfers)
  - [Batch Payments: Payment Management](#batch-payments--payment-management)
  - [Batch Payments: Bank Slips](#batch-payments--bank-slips)
  - [Batch Payments: Barcode Guides](#batch-payments--barcode-guides)
  - [Batch Payments: DARF](#batch-payments--darf)
  - [Batch Payments: GPS](#batch-payments--gps)
  - [Batch Payments: GRU](#batch-payments--gru)
- [Error Handling](#error-handling)
- [Retry Behavior](#retry-behavior)
- [Constants Reference](#constants-reference)
- [Project Layout](#project-layout)
- [Testing](#testing)
- [License](#license)

---

## Features

- OAuth2 `client_credentials` flow with automatic token caching and renewal
- Mutual TLS (mTLS) support required for BB production endpoints that demand a client certificate
- Typed request and response structs for all supported endpoints
- Built-in exponential backoff retry for transient failures
- Sandbox and production environments switchable via a single config flag
- Thread-safe client, safe for concurrent use across goroutines
- Easily extendable, new API resources can be added without modifying the core client

---

## Requirements

- Go **1.25** or later

---

## Installation

```bash
go get github.com/raykavin/bbapi-go
```

---

## Quick Start

The example below authenticates using the OAuth2 `client_credentials` flow and submits a batch of TED/DOC transfers.

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

	// Authenticate and cache the token, subsequent calls reuse it automatically.
	if _, err = client.Authenticate(ctx); err != nil {
		log.Fatal(err)
	}

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
	Sandbox:      false,
	APIURL:       "",             // overrides the default base URL
	AuthURL:      "",             // overrides the default OAuth token URL
	AccessToken:  "",             // pre-load a token instead of authenticating
	Scopes:       nil,
	HTTPClient:   nil,            // uses a default *http.Client
	Timeout:      30 * time.Second,
	MaxRetries:   3,              // set to -1 to disable retries entirely
	RetryWaitMin: 1 * time.Second,
	RetryWaitMax: 10 * time.Second,

	// mTLS required for production endpoints that demand a client certificate
	MTLSCertFile: "",             // path to PEM-encoded client certificate file
	MTLSKeyFile:  "",             // path to PEM-encoded private key file
})
```

### Config Fields

| Field | Type | Description |
|---|---|---|
| `ClientID` | `string` | Banco do Brasil application client ID |
| `ClientSecret` | `string` | Banco do Brasil application client secret |
| `AppKey` | `string` | Application key required by the BB API gateway |
| `Sandbox` | `bool` | When `true`, routes all requests to BB sandbox endpoints |
| `APIURL` | `string` | Optional override for the API base URL |
| `AuthURL` | `string` | Optional override for the OAuth token URL |
| `AccessToken` | `string` | Pre-load an existing access token, skipping initial authentication |
| `Scopes` | `[]Scope` | OAuth scopes requested during authentication |
| `HTTPClient` | `*http.Client` | Optional custom HTTP client |
| `Timeout` | `time.Duration` | Per-request timeout (default: `30s`) |
| `MaxRetries` | `int` | Number of retry attempts for transient failures (default: `3`) |
| `RetryWaitMin` | `time.Duration` | Minimum backoff between retries (default: `1s`) |
| `RetryWaitMax` | `time.Duration` | Maximum backoff between retries (default: `10s`) |
| `MTLSCertFile` | `string` | Path to the PEM-encoded client certificate file |
| `MTLSKeyFile` | `string` | Path to the PEM-encoded private key file |
| `MTLSCertPEM` | `[]byte` | Client certificate as raw PEM bytes (alternative to file path) |
| `MTLSKeyPEM` | `[]byte` | Private key as raw PEM bytes (alternative to file path) |
| `MTLSCARootFile` | `string` | Path to a custom CA root certificate file (optional) |
| `MTLSCARootPEM` | `[]byte` | Custom CA root certificate as raw PEM bytes (optional) |
| `MTLSEnabled` | `bool` | Set to `true` when providing a pre-configured `HTTPClient` with mTLS |

### Default Endpoints

| Environment | API | OAuth |
|---|---|---|
| Sandbox | `https://homologa-api-ip.bb.com.br:7144/pagamentos-lote/v1` | `https://oauth.sandbox.bb.com.br/oauth/token` |
| Production | `https://api-ip.bb.com.br/pagamentos-lote/v1` | `https://oauth.bb.com.br/oauth/token` |

---

## Authentication

The SDK uses the OAuth2 `client_credentials` flow. A valid token is required before any API call. Once obtained, the client caches it and reuses it for subsequent requests. If the token expires or a `401` is returned, the client re-authenticates automatically.

### Authenticate with configured credentials

```go
// Uses ClientID, ClientSecret, and Scopes from Config.
// The token is cached internally; no need to call SetTokenResponse manually.
if _, err := client.Authenticate(ctx); err != nil {
	log.Fatal(err)
}
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
// Useful when your application manages token acquisition externally.
client.SetAccessToken("your-access-token")
```

### Inspect the current token

```go
token     := client.GetAccessToken()
expiresAt := client.TokenExpiresAt()
```

---

## Mutual TLS (mTLS)

Several Banco do Brasil production endpoints require the HTTP client to present a certificate during the TLS handshake ([see the official list](https://apoio.developers.bb.com.br/guias-e-tutoriais/seguranca/apis-que-exigem-certificado)). The SDK supports this natively when any mTLS field is set and `HTTPClient` is `nil`, a transport with the correct `tls.Config` is built automatically.

Calling a method that requires mTLS without configuring a certificate returns `bbapi.ErrMTLSRequired` immediately, before any network request is made.

### Using certificate files

```go
client, err := bbapi.NewClient(bbapi.Config{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	AppKey:       "your-app-key",
	MTLSCertFile: "/etc/ssl/bb/client.crt",
	MTLSKeyFile:  "/etc/ssl/bb/client.key",
})
```

### Using PEM bytes (e.g. from a secrets manager)

```go
certPEM, _ := secretsManager.GetSecret("bb-client-cert")
keyPEM,  _ := secretsManager.GetSecret("bb-client-key")

client, err := bbapi.NewClient(bbapi.Config{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	AppKey:       "your-app-key",
	MTLSCertPEM:  certPEM,
	MTLSKeyPEM:   keyPEM,
})
```

### Using a custom CA root

```go
client, err := bbapi.NewClient(bbapi.Config{
	ClientID:       "your-client-id",
	ClientSecret:   "your-client-secret",
	AppKey:         "your-app-key",
	MTLSCertFile:   "/etc/ssl/bb/client.crt",
	MTLSKeyFile:    "/etc/ssl/bb/client.key",
	MTLSCARootFile: "/etc/ssl/bb/ca-root.crt",
})
```

### Bringing your own http.Client

If you build and configure the `http.Client` yourself (e.g. with additional middleware), set `MTLSEnabled: true` so the SDK knows the client already presents a certificate:

```go
client, err := bbapi.NewClient(bbapi.Config{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	AppKey:       "your-app-key",
	HTTPClient:   myPreConfiguredMTLSClient,
	MTLSEnabled:  true,
})
```

### APIs that require mTLS

The following methods enforce the mTLS requirement and return `bbapi.ErrMTLSRequired` if no certificate is configured:

| Group | Methods |
|---|---|
| DARF | `CreateDARFBatch`, `GetDARFBatchRequest`, `GetDARFPayment` |
| GPS | `CreateGPSBatch`, `GetGPSBatchRequest`, `GetGPSPayment` |
| GRU | `CreateGRUBatch`, `GetGRUBatchRequest`, `GetGRUPayment` |
| Bank Slips | `CreateBankSlipBatch`, `GetBankSlipBatchRequest`, `GetBankSlipPayment` |
| Barcode Guides | `CreateBarcodeGuideBatch`, `GetBarcodeGuideBatchRequest`, `GetBarcodeGuidePayment` |
| Payments | `ReleasePayments`, `CancelPayments`, `UpdatePaymentDates`, `ListReturnedPayments`, `ListPaymentEntries`, `GetBarcodePayments` |
| Transfers | `ListTransferBatches`, `CreateTransferBatch`, `GetTransferPayment`, `GetBatchRequest`, `GetBatch`, `ListBeneficiaryTransfers`, `CreatePixTransferBatch`, `GetPixTransferBatchRequest`, `GetPixPayment` |

---

## OAuth Scopes

The package exposes typed `Scope` constants for every permission supported by the API. Pass the scopes your application needs during initialization or authentication.

| Constant | Value |
|---|---|
| `ScopeBatchesInfo` | `pagamentos-lote.lotes-info` |
| `ScopeBatchesRequest` | `pagamentos-lote.lotes-requisicao` |
| `ScopePaymentsInfo` | `pagamentos-lote.pagamentos-info` |
| `ScopeReturnedPaymentsInfo` | `pagamentos-lote.devolvidos-info` |
| `ScopeCancelRequest` | `pagamentos-lote.cancelar-requisicao` |
| `ScopeTransfersInfo` | `pagamentos-lote.transferencias-info` |
| `ScopeTransfersRequest` | `pagamentos-lote.transferencias-requisicao` |
| `ScopePixInfo` | `pagamentos-lote.pix-info` |
| `ScopePixTransfersInfo` | `pagamentos-lote.transferencias-pix-info` |
| `ScopePixTransfersRequest` | `pagamentos-lote.transferencias-pix-requisicao` |
| `ScopeBankSlipsInfo` | `pagamentos-lote.boletos-info` |
| `ScopeBankSlipsRequest` | `pagamentos-lote.boletos-requisicao` |
| `ScopeBarcodeGuidesInfo` | `pagamentos-lote.guias-codigo-barras-info` |
| `ScopeBarcodeGuidesRequest` | `pagamentos-lote.guias-codigo-barras-requisicao` |
| `ScopeBarcodePaymentsInfo` | `pagamentos-lote.pagamentos-codigo-barras-info` |
| `ScopeManualGuidePaymentsInfo` | `pagamentos-lote.pagamentos-guias-sem-codigo-barras-info` |
| `ScopeManualGuidePaymentsRequest` | `pagamentos-lote.pagamentos-guias-sem-codigo-barras-requisicao` |

---

## API Coverage

The SDK implements all **30** operations defined in the current Banco do Brasil Batch Payments API specification.

---

### Batch Payments Transfers (TED/DOC)

Submit and query batches of TED/DOC bank transfers.

| Method | Description |
|---|---|
| `CreateTransferBatch(ctx, *CreateTransferBatchRequest)` | Submit a new transfer batch |
| `ListTransferBatches(ctx, *ListTransferBatchesParams)` | List existing transfer batches |
| `GetBatch(ctx, id)` | Retrieve a batch by ID |
| `GetBatchRequest(ctx, id)` | Retrieve the request-stage details of a batch |
| `GetTransferPayment(ctx, id, *AccountLookupParams)` | Retrieve a single transfer payment |
| `ListBeneficiaryTransfers(ctx, id, *ListBeneficiaryTransfersParams)` | List transfers for a specific beneficiary |

**Example creating a transfer batch:**

```go
resp, err := client.CreateTransferBatch(ctx, &bbapi.CreateTransferBatchRequest{
	RequestNumber: 42,
	PaymentType:   bbapi.PaymentTypeSalary,
	Transfers: []bbapi.Transfer{
		{TransferDate: 15052026, TransferValue: 3500.00},
		{TransferDate: 15052026, TransferValue: 2800.50},
	},
})
if err != nil {
	log.Fatal(err)
}
log.Printf("batch_id=%d state=%d", resp.RequestIdentifier, resp.RequestState)
```

---

### Batch Payments Pix Transfers

Submit and query Pix transfer batches.

| Method | Description |
|---|---|
| `CreatePixTransferBatch(ctx, *CreatePixTransferBatchRequest)` | Submit a new Pix transfer batch |
| `GetPixTransferBatchRequest(ctx, id)` | Retrieve the request-stage details of a Pix batch |
| `GetPixPayment(ctx, id, *GetPixPaymentParams)` | Retrieve a single Pix payment |

**Example creating a Pix transfer batch:**

```go
resp, err := client.CreatePixTransferBatch(ctx, &bbapi.CreatePixTransferBatchRequest{
	RequestNumber: 99,
	PixTransfers: []bbapi.PixTransfer{
		{TransferDate: 15052026, TransferValue: 500.00},
	},
})
if err != nil {
	log.Fatal(err)
}
log.Printf("batch_id=%d state=%d", resp.RequestIdentifier, resp.RequestState)
```

---

### Batch Payments Payment Management

Cross-cutting operations that apply to payments regardless of type.

| Method | Description |
|---|---|
| `ReleasePayments(ctx, *ReleasePaymentsRequest)` | Release a batch for processing |
| `CancelPayments(ctx, *CancelPaymentsRequest)` | Cancel one or more payments |
| `UpdatePaymentDates(ctx, id, *UpdatePaymentDatesRequest)` | Reschedule payment dates |
| `ListReturnedPayments(ctx, *ListReturnedPaymentsParams)` | List returned/reversed payments |
| `ListPaymentEntries(ctx, *ListPaymentEntriesParams)` | List payment ledger entries for a period |
| `GetBarcodePayments(ctx, id, *AccountLookupParams)` | Retrieve barcode-linked payment details |

---

### Batch Payments Bank Slips

Submit and query bank slip (_boleto_) batches.

| Method | Description |
|---|---|
| `CreateBankSlipBatch(ctx, *CreateBankSlipBatchRequest)` | Submit a new bank slip batch |
| `GetBankSlipBatchRequest(ctx, id, *AccountLookupParams)` | Retrieve the request-stage details |
| `GetBankSlipPayment(ctx, id, *AccountLookupParams)` | Retrieve a single bank slip payment |

---

### Batch Payments Barcode Guides

Submit and query barcode guide (_guia de código de barras_) batches.

| Method | Description |
|---|---|
| `CreateBarcodeGuideBatch(ctx, *CreateBarcodeGuideBatchRequest)` | Submit a new barcode guide batch |
| `GetBarcodeGuideBatchRequest(ctx, id, *AccountLookupParams)` | Retrieve the request-stage details |
| `GetBarcodeGuidePayment(ctx, id, *AccountLookupParams)` | Retrieve a single barcode guide payment |

---

### Batch Payments DARF

Submit and query DARF (federal tax collection document) batches.

| Method | Description |
|---|---|
| `CreateDARFBatch(ctx, *CreateDARFBatchRequest)` | Submit a new DARF batch |
| `GetDARFBatchRequest(ctx, id, *AccountLookupParams)` | Retrieve the request-stage details |
| `GetDARFPayment(ctx, id, *AccountLookupParams)` | Retrieve a single DARF payment |

---

### Batch Payments GPS

Submit and query GPS (social security guide) batches.

| Method | Description |
|---|---|
| `CreateGPSBatch(ctx, *CreateGPSBatchRequest)` | Submit a new GPS batch |
| `GetGPSBatchRequest(ctx, id, *AccountLookupParams)` | Retrieve the request-stage details |
| `GetGPSPayment(ctx, id, *AccountLookupParams)` | Retrieve a single GPS payment |

---

### Batch Payments GRU

Submit and query GRU (federal government collection guide) batches.

| Method | Description |
|---|---|
| `CreateGRUBatch(ctx, *CreateGRUBatchRequest)` | Submit a new GRU batch |
| `GetGRUBatchRequest(ctx, id, *AccountLookupParams)` | Retrieve the request-stage details |
| `GetGRUPayment(ctx, id, *AccountLookupParams)` | Retrieve a single GRU payment |

---

## Error Handling

All API and OAuth errors are normalized into `*bbapi.APIError`, giving you a consistent interface regardless of which endpoint was called.

```go
resp, err := client.GetTransferPayment(ctx, "12345", nil)
if err != nil {
	var apiErr *bbapi.APIError
	if errors.As(err, &apiErr) {
		log.Printf("HTTP %d: %s", apiErr.StatusCode, apiErr.Message)
		for _, detail := range apiErr.Details {
			log.Printf("  [%s] %s", detail.Codigo, detail.Mensagem)
		}
	}
	log.Fatal(err)
}
```

The following helper functions are available for common status checks:

| Function | Condition |
|---|---|
| `bbapi.IsNotFound(err)` | HTTP 404 |
| `bbapi.IsUnauthorized(err)` | HTTP 401 |
| `bbapi.IsForbidden(err)` | HTTP 403 |
| `bbapi.IsRateLimited(err)` | HTTP 429 |
| `bbapi.IsServerError(err)` | HTTP 5xx |

```go
if bbapi.IsNotFound(err) {
	// handle missing resource
}
```

`bbapi.ErrMTLSRequired` is a sentinel error returned when a method that requires mutual TLS is called without a client certificate configured. Check for it with `errors.Is`:

```go
resp, err := client.CreateDARFBatch(ctx, req)
if errors.Is(err, bbapi.ErrMTLSRequired) {
	log.Fatal("configure MTLSCertFile/MTLSKeyFile to use this API in production")
}
```

---

## Retry Behavior

The client automatically retries requests that fail due to transient conditions:

| Status | Meaning |
|---|---|
| `429` | Too Many Requests |
| `500` | Internal Server Error |
| `502` | Bad Gateway |
| `503` | Service Unavailable |
| `504` | Gateway Timeout |

Retries use exponential backoff bounded by `RetryWaitMin` and `RetryWaitMax`. The number of attempts is controlled by `MaxRetries` (default: `3`). Set `MaxRetries: -1` to disable retries entirely.

A `401 Unauthorized` response causes the client to clear the cached token and re-authenticate before the next attempt.

---

## Constants Reference

### Payment Types

| Constant | Value | Description |
|---|---|---|
| `PaymentTypeSuppliers` | `126` | Supplier payments |
| `PaymentTypeSalary` | `127` | Payroll / salary payments |
| `PaymentTypeMiscellaneous` | `128` | General-purpose payments |

### Request States

| Constant | Value | Description |
|---|---|---|
| `PaymentRequestStateConsistent` | `1` | Request is consistent |
| `PaymentRequestStateInconsistent` | `2` | Request has inconsistencies |
| `PaymentRequestStateAllInconsistent` | `3` | All entries are inconsistent |
| `PaymentRequestStatePending` | `4` | Awaiting release |
| `PaymentRequestStateProcessing` | `5` | Being processed |
| `PaymentRequestStateProcessed` | `6` | Successfully processed |
| `PaymentRequestStateRejected` | `7` | Rejected |
| `PaymentRequestStatePreparingUnreleased` | `8` | Preparing not yet released |
| `PaymentRequestStateReleasedByAPI` | `9` | Released via API |
| `PaymentRequestStatePreparingReleased` | `10` | Preparing already released |

---

## Project Layout

| File | Responsibility |
|---|---|
| `client.go` | Core client, token lifecycle, request execution, retry integration |
| `config.go` | SDK configuration struct and default values |
| `auth.go` | OAuth2 authentication |
| `errors.go` | API error types and helper functions |
| `helpers.go` | Generic request/response helpers |
| `payments.go` | General payment management endpoints |
| `transfers.go` | TED/DOC and Pix transfer endpoints |
| `bank_slips.go` | Bank slip (_boleto_) endpoints |
| `barcode_guides.go` | Barcode guide endpoints |
| `darf.go` | DARF (federal tax) endpoints |
| `gps.go` | GPS (social security) endpoints |
| `gru.go` | GRU (public sector) endpoints |
| `doc.go` | Package-level documentation and constants |

---

## Testing

Run the full test suite with:

```bash
go test ./...
```

The tests cover client initialization, token management, OAuth flows, request construction, model serialization, response parsing, error handling, retry behavior, and mTLS enforcement.

---
## Contributing

Contributions to bbapi-go are welcome! Here are some ways you can help improve the project:

- **Report bugs and suggest features** by opening issues on GitHub
- **Submit pull requests** with bug fixes or new features
- **Improve documentation** to help other users and developers
- **Share your custom strategies** with the community

---

## License
bbapi-go is distributed under the **MIT License**.  
For complete license terms and conditions, see the [LICENSE](LICENSE.md) file in the repository.

---

## Contact

For support, collaboration, or questions about bbapi-go:

**Email**: [raykavin.meireles@gmail.com](mailto:raykavin.meireles@gmail.com)  
**LinkedIn**: [@raykavin.dev](https://www.linkedin.com/in/raykavin-dev)  
**GitHub**: [@raykavin](https://github.com/raykavin)  
