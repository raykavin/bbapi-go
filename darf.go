package bbapi

import (
	"context"
	"fmt"
	"net/url"
)

const (
	endpointDARFBatch        = "/lotes-darf-normal-preto"
	endpointDARFBatchRequest = "/lotes-darf-preto-normal/%s/solicitacao"
	endpointDARFPayment      = "/darf-preto/%s"
)

// DARFEntry represents a DARF payment entry.
type DARFEntry struct {
	PaymentDate            int64    `json:"dataPagamento"`
	PaymentValue           float64  `json:"valorPagamento"`
	DebitDocumentNumber    *int64   `json:"numeroDocumentoDebito,omitempty"`
	YourDocumentCode       *string  `json:"codigoSeuDocumento,omitempty"`
	PaymentDescription     *string  `json:"textoDescricaoPagamento,omitempty"`
	TaxRevenueCode         *int64   `json:"codigoReceitaTributo,omitempty"`
	TaxpayerTypeCode       *int64   `json:"codigoTipoContribuinte,omitempty"`
	TaxpayerIdentification *int64   `json:"numeroIdentificacaoContribuinte,omitempty"`
	TaxIdentifierCode      *string  `json:"codigoIdentificadorTributo,omitempty"`
	AssessmentDate         *int64   `json:"dataApuracao,omitempty"`
	ReferenceNumber        *int64   `json:"numeroReferencia,omitempty"`
	PrincipalValue         *float64 `json:"valorPrincipal,omitempty"`
	FineValue              *float64 `json:"valorMulta,omitempty"`
	InterestValue          *float64 `json:"valorJuroEncargo,omitempty"`
	DueDate                *int64   `json:"dataVencimento,omitempty"`
}

// CreateDARFBatchRequest is the request body for POST /lotes-darf-normal-preto.
type CreateDARFBatchRequest struct {
	RequestID              int64       `json:"id"`
	ContractCode           *int64      `json:"codigoContrato,omitempty"`
	DebitAgencyNumber      *int64      `json:"numeroAgenciaDebito,omitempty"`
	DebitAccountNumber     *int64      `json:"numeroContaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string     `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	Entries                []DARFEntry `json:"lancamentos"`
}

// DARFBatchEntryResult represents a DARF batch entry returned by batch creation.
type DARFBatchEntryResult struct {
	PaymentIdentifier      int64   `json:"codigoIdentificadorPagamento,omitempty"`
	ConventionName         string  `json:"nomeConvenente,omitempty"`
	PaymentDate            int64   `json:"dataPagamento,omitempty"`
	PaymentValue           float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber    int64   `json:"numeroDocumentoDebito,omitempty"`
	YourDocumentCode       string  `json:"codigoSeuDocumento,omitempty"`
	PaymentDescription     string  `json:"textoDescricaoPagamento,omitempty"`
	TaxRevenueCode         int64   `json:"codigoReceitaTributo,omitempty"`
	TaxpayerTypeCode       int64   `json:"codigoTipoContribuinte,omitempty"`
	TaxpayerIdentification int64   `json:"numeroIdentificacaoContribuinte,omitempty"`
	TaxIdentifierCode      string  `json:"codigoIdentificadorTributo,omitempty"`
	AssessmentDate         int64   `json:"dataApuracao,omitempty"`
	ReferenceNumber        int64   `json:"numeroReferencia,omitempty"`
	PrincipalValue         float64 `json:"valorPrincipal,omitempty"`
	FineValue              float64 `json:"valorMulta,omitempty"`
	InterestValue          float64 `json:"valorJuroEncargo,omitempty"`
	DueDate                int64   `json:"dataVencimento,omitempty"`
	AcceptanceIndicator    string  `json:"indicadorMovimentoAceito,omitempty"`
	Errors                 []int64 `json:"erros,omitempty"`
}

// DARFBatchRequestPayment represents a payment nested inside a DARF batch request lookup.
type DARFBatchRequestPayment struct {
	ID                  int64   `json:"id,omitempty"`
	Date                int64   `json:"data,omitempty"`
	Value               float64 `json:"valor,omitempty"`
	PrincipalValue      float64 `json:"valorPrincipal,omitempty"`
	FineValue           float64 `json:"valorMulta,omitempty"`
	InterestValue       float64 `json:"valorJuroEncargo,omitempty"`
	TaxpayerDocument    int64   `json:"cpfCnpjContribuinte,omitempty"`
	DebitDocumentNumber int64   `json:"numeroDocumentoDebito,omitempty"`
	DescriptionText     string  `json:"textoDescricao,omitempty"`
}

// DARFBatchRequestEntry represents a DARF batch entry returned by request lookups.
type DARFBatchRequestEntry struct {
	ConventionName       string                    `json:"nomeConvenente,omitempty"`
	ReferenceNumber      string                    `json:"numeroReferencia,omitempty"`
	DueDate              int64                     `json:"dataVencimento,omitempty"`
	Payments             []DARFBatchRequestPayment `json:"pagamento,omitempty"`
	AcceptanceIndicator  string                    `json:"indicadorMovimentoAceito,omitempty"`
	Errors               []int64                   `json:"erros,omitempty"`
	ContractCustomerCode int64                     `json:"codigoClienteContrato,omitempty"`
	DocumentCode         string                    `json:"codigoDocumento,omitempty"`
	TaxRevenueCode       int64                     `json:"codigoReceitaTributo,omitempty"`
	TaxIdentifierCode    string                    `json:"codigoIdentificadorTributo,omitempty"`
	TaxpayerTypeCode     int64                     `json:"codigoTipoContribuinte,omitempty"`
	AssessmentDate       int64                     `json:"dataApuracao,omitempty"`
}

// DARFPaymentItem represents a DARF payment detail item.
type DARFPaymentItem struct {
	Code                   int64   `json:"codigo,omitempty"`
	TaxpayerType           int64   `json:"tipoContribuinte,omitempty"`
	TaxpayerIdentification int64   `json:"identificacaoContribuinte,omitempty"`
	TaxIdentifier          string  `json:"identificacaoTributo,omitempty"`
	AssessmentDate         string  `json:"dataApuracao,omitempty"`
	ReferenceNumber        int64   `json:"numeroReferencia,omitempty"`
	PrincipalValue         float64 `json:"valorPrincipal,omitempty"`
	FineValue              float64 `json:"valorMulta,omitempty"`
	InterestValue          float64 `json:"valorJuroEncargo,omitempty"`
	DueDate                string  `json:"dataVencimento,omitempty"`
	FreeText               string  `json:"textoLivre,omitempty"`
}

// DARFReturnItem represents a DARF return entry.
type DARFReturnItem struct {
	ReasonCode int64 `json:"codigoMotivo,omitempty"`
}

// CreateDARFBatchResponse is the response body for POST /lotes-darf-normal-preto.
type CreateDARFBatchResponse struct {
	ID              int64                  `json:"id,omitempty"`
	StateCode       int64                  `json:"codigoEstado,omitempty"`
	EntryCount      int64                  `json:"quantidadeLancamentos,omitempty"`
	EntryValue      float64                `json:"valorLancamentos,omitempty"`
	ValidEntryCount int64                  `json:"quantidadeLancamentosValidos,omitempty"`
	ValidEntryValue float64                `json:"valorLancamentosValidos,omitempty"`
	Entries         []DARFBatchEntryResult `json:"lancamentos,omitempty"`
}

// GetDARFBatchRequestResponse is the response body for GET /lotes-darf-preto-normal/{id}/solicitacao.
type GetDARFBatchRequestResponse struct {
	ID              int64                   `json:"id,omitempty"`
	StateCode       int64                   `json:"codigoEstado,omitempty"`
	EntryCount      int64                   `json:"quantidadeLancamentos,omitempty"`
	EntryValue      float64                 `json:"valorLancamentos,omitempty"`
	ValidEntryCount int64                   `json:"quantidadeLancamentosValidos,omitempty"`
	ValidEntryValue float64                 `json:"valorLancamentosValidos,omitempty"`
	Entries         []DARFBatchRequestEntry `json:"lancamentos,omitempty"`
}

// GetDARFPaymentResponse is the response body for GET /darf-preto/{id}.
type GetDARFPaymentResponse struct {
	ID                  int64             `json:"id"`
	PaymentState        string            `json:"estadoPagamento,omitempty"`
	DebitAgency         int64             `json:"agenciaDebito,omitempty"`
	DebitAccount        int64             `json:"contaCorrenteDebito,omitempty"`
	DebitAccountDigit   string            `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CreditCardStart     int64             `json:"inicioCartaoCredito,omitempty"`
	CreditCardEnd       int64             `json:"fimCartaoCredito,omitempty"`
	PaymentDate         string            `json:"dataPagamento,omitempty"`
	PaymentValue        float64           `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64             `json:"documentoDebito,omitempty"`
	AuthenticationCode  string            `json:"codigoAutenticacaoPagamento,omitempty"`
	PaymentItems        []DARFPaymentItem `json:"listaPagamentos,omitempty"`
	ReturnItems         []DARFReturnItem  `json:"listaDevolucao,omitempty"`
}

// CreateDARFBatch creates a DARF batch.
// Requires mutual TLS — see Config.MTLSCertFile / Config.MTLSCertPEM.
func (c *Client) CreateDARFBatch(
	ctx context.Context,
	req *CreateDARFBatchRequest,
) (*CreateDARFBatchResponse, error) {
	if err := c.requireMTLS(); err != nil {
		return nil, err
	}
	return post[*CreateDARFBatchResponse](c, ctx, endpointDARFBatch, req)
}

// GetDARFBatchRequest returns the request-stage representation of a DARF batch.
// Requires mutual TLS see Config.MTLSCertFile / Config.MTLSCertPEM.
func (c *Client) GetDARFBatchRequest(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetDARFBatchRequestResponse, error) {
	if err := c.requireMTLS(); err != nil {
		return nil, err
	}

	query := url.Values{}
	setAccountLookupQuery(query, params, "numeroAgenciaDebito", "numeroContaCorrenteDebito", "digitoVerificadorContaCorrenteDebito")
	return get[*GetDARFBatchRequestResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointDARFBatchRequest, id), query),
	)
}

// GetDARFPayment returns a single DARF payment.
// Requires mutual TLS — see Config.MTLSCertFile / Config.MTLSCertPEM.
func (c *Client) GetDARFPayment(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetDARFPaymentResponse, error) {
	if err := c.requireMTLS(); err != nil {
		return nil, err
	}
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetDARFPaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointDARFPayment, id), query),
	)
}
