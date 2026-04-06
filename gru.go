package bbapi

import (
	"context"
	"fmt"
	"net/url"
)

const (
	endpointGRUBatch        = "/pagamentos-gru"
	endpointGRUBatchRequest = "/lotes-gru/%s/solicitacao" // NOTE: prefix differs from endpointGRUBatch — intentional per BB API docs.
	endpointGRUPayment      = "/gru/%s"
)

// GRUEntry represents a GRU payment entry.
type GRUEntry struct {
	Barcode             string   `json:"codigoBarras"`
	DueDate             *int64   `json:"dataVencimento,omitempty"`
	PaymentDate         int64    `json:"dataPagamento"`
	PaymentValue        float64  `json:"valorPagamento"`
	DebitDocumentNumber *int64   `json:"numeroDocumentoDebito,omitempty"`
	PaymentText         *string  `json:"textoPagamento,omitempty"`
	ReferenceNumber     *string  `json:"numeroReferencia,omitempty"`
	CompetenceMonthYear *int64   `json:"mesAnoCompetencia,omitempty"`
	TaxpayerID          *int64   `json:"idContribuinte,omitempty"`
	PrincipalValue      *float64 `json:"valorPrincipal,omitempty"`
	DiscountValue       *float64 `json:"valorDesconto,omitempty"`
	OtherDeductionValue *float64 `json:"valorOutraDeducao,omitempty"`
	FineValue           *float64 `json:"valorMulta,omitempty"`
	InterestValue       *float64 `json:"valorJuroEncargo,omitempty"`
	OtherIncreaseValue  *float64 `json:"valorOutroAcrescimo,omitempty"`
}

// CreateGRUBatchRequest is the request body for POST /pagamentos-gru.
type CreateGRUBatchRequest struct {
	RequestNumber     int64      `json:"numeroRequisicao"`
	ContractCode      *int64     `json:"codigoContrato,omitempty"`
	Agency            *int64     `json:"agencia,omitempty"`
	Account           *int64     `json:"conta,omitempty"`
	AccountCheckDigit *string    `json:"digitoConta,omitempty"`
	Entries           []GRUEntry `json:"listaRequisicao"`
}

// GRUBatchPaymentResult represents a payment returned by GRU batch creation.
type GRUBatchPaymentResult struct {
	PaymentID           int64   `json:"idPagamento,omitempty"`
	ReceiverName        string  `json:"nomeRecebedor,omitempty"`
	Barcode             string  `json:"codigoBarras,omitempty"`
	DueDate             int64   `json:"dataVencimento,omitempty"`
	PaymentDate         int64   `json:"dataPagamento,omitempty"`
	PaymentValue        float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64   `json:"numeroDocumentoDebito,omitempty"`
	PaymentDescription  string  `json:"descricaoPagamento,omitempty"`
	ReferenceNumber     string  `json:"numeroReferencia,omitempty"`
	CompetenceMonthYear int64   `json:"mesAnoCompetencia,omitempty"`
	TaxpayerID          int64   `json:"idContribuinte,omitempty"`
	PrincipalValue      float64 `json:"valorPrincipal,omitempty"`
	DiscountValue       float64 `json:"valorDesconto,omitempty"`
	OtherDeductionValue float64 `json:"valorOutraDeducao,omitempty"`
	FineValue           float64 `json:"valorMulta,omitempty"`
	InterestValue       float64 `json:"valorJuroEncargo,omitempty"`
	OtherIncreaseValue  float64 `json:"valorOutrosAcrescimos,omitempty"`
	AcceptanceIndicator string  `json:"indicadorMovimentoAceito,omitempty"`
	Errors              []int64 `json:"erros,omitempty"`
}

// GRUBatchRequestPayment represents a nested payment returned by GRU batch request lookups.
type GRUBatchRequestPayment struct {
	ID                  int64   `json:"id,omitempty"`
	Date                int64   `json:"data,omitempty"`
	Value               float64 `json:"valor,omitempty"`
	PrincipalValue      float64 `json:"valorPrincipal,omitempty"`
	DiscountValue       float64 `json:"valorDesconto,omitempty"`
	OtherDeductionValue float64 `json:"valorOutroDeducao,omitempty"`
	FineValue           float64 `json:"valorMulta,omitempty"`
	InterestValue       float64 `json:"valorJuroEncargo,omitempty"`
	OtherValue          float64 `json:"valorOutro,omitempty"`
	TaxpayerDocument    int64   `json:"cpfCnpjContribuinte,omitempty"`
	DebitDocumentNumber int64   `json:"numeroDocumentoDebito,omitempty"`
	DescriptionText     string  `json:"textoDescricao,omitempty"`
}

// GRUBatchRequestEntry represents a GRU batch entry returned by request lookups.
type GRUBatchRequestEntry struct {
	ConventionName      string                   `json:"nomeConvenente,omitempty"`
	BarcodeText         string                   `json:"textoCodigoBarras,omitempty"`
	ReferenceNumber     string                   `json:"numeroReferencia,omitempty"`
	DueDate             int64                    `json:"dataVencimento,omitempty"`
	CompetenceMonthYear int64                    `json:"mesAnoCompetencia,omitempty"`
	Payments            []GRUBatchRequestPayment `json:"pagamento,omitempty"`
	AcceptanceIndicator string                   `json:"indicadorMovimentoAceito,omitempty"`
	Errors              []int64                  `json:"erros,omitempty"`
}

// GRUPaymentItem represents a GRU payment detail item.
type GRUPaymentItem struct {
	Code                   string  `json:"codigo,omitempty"`
	ReceiverName           string  `json:"nomeRecebedor,omitempty"`
	ReferenceNumber        string  `json:"numeroReferencia,omitempty"`
	CompetenceMonthYear    int64   `json:"mesAnoCompetencia,omitempty"`
	DueDate                int64   `json:"dataVencimento,omitempty"`
	TaxpayerIdentification int64   `json:"identificacaoContribuinte,omitempty"`
	PrincipalValue         float64 `json:"valorPrincipal,omitempty"`
	DiscountValue          float64 `json:"valorDesconto,omitempty"`
	OtherDeductionValue    float64 `json:"valorOutroDeducao,omitempty"`
	FineValue              float64 `json:"valorMulta,omitempty"`
	InterestValue          float64 `json:"valorJuroEncargo,omitempty"`
	OtherValue             float64 `json:"valorOutro,omitempty"`
	Text                   string  `json:"texto,omitempty"`
}

// GRUOccurrenceItem represents a GRU occurrence entry.
type GRUOccurrenceItem struct {
	Code int64 `json:"codigo,omitempty"`
}

// CreateGRUBatchResponse is the response body for POST /pagamentos-gru.
type CreateGRUBatchResponse struct {
	RequestNumber   int64                   `json:"numeroRequisicao,omitempty"`
	RequestState    int64                   `json:"estadoRequisicao,omitempty"`
	TotalCount      int64                   `json:"quantidadeTotal,omitempty"`
	TotalValue      float64                 `json:"valorTotal,omitempty"`
	ValidTotalCount int64                   `json:"quantidadeTotalValido,omitempty"`
	ValidTotalValue float64                 `json:"valorTotalValido,omitempty"`
	Payments        []GRUBatchPaymentResult `json:"pagamentos,omitempty"`
}

// GetGRUBatchRequestResponse is the response body for GET /lotes-gru/{id}/solicitacao.
type GetGRUBatchRequestResponse struct {
	ID              int64                  `json:"id,omitempty"`
	StateCode       int64                  `json:"codigoEstado,omitempty"`
	EntryCount      int64                  `json:"quantidadeLancamentos,omitempty"`
	EntryValue      float64                `json:"valorLancamentos,omitempty"`
	ValidEntryCount int64                  `json:"quantidadeLancamentosValidos,omitempty"`
	ValidEntryValue float64                `json:"valorLancamentosValidos,omitempty"`
	Entries         []GRUBatchRequestEntry `json:"lancamentos,omitempty"`
}

// GetGRUPaymentResponse is the response body for GET /gru/{id}.
type GetGRUPaymentResponse struct {
	ID                  int64               `json:"id"`
	PaymentState        string              `json:"estadoPagamento,omitempty"`
	DebitAgency         int64               `json:"agenciaDebito,omitempty"`
	DebitAccount        int64               `json:"contaCorrenteDebito,omitempty"`
	DebitAccountDigit   string              `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CreditCardStart     int64               `json:"inicioCartaoCredito,omitempty"`
	CreditCardEnd       int64               `json:"fimCartaoCredito,omitempty"`
	PaymentDate         int64               `json:"dataPagamento,omitempty"`
	PaymentValue        float64             `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64               `json:"documentoDebito,omitempty"`
	AuthenticationCode  string              `json:"codigoAutenticacaoPagamento,omitempty"`
	PaymentItems        []GRUPaymentItem    `json:"listaPagamentos,omitempty"`
	OccurrenceItems     []GRUOccurrenceItem `json:"listaOcorrencias,omitempty"`
}

// CreateGRUBatch creates a GRU batch.
func (c *Client) CreateGRUBatch(
	ctx context.Context,
	req *CreateGRUBatchRequest,
) (*CreateGRUBatchResponse, error) {
	return post[*CreateGRUBatchResponse](c, ctx, endpointGRUBatch, req)
}

// GetGRUBatchRequest returns the request-stage representation of a GRU batch.
func (c *Client) GetGRUBatchRequest(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetGRUBatchRequestResponse, error) {
	if err := c.requireMTLS(); err != nil {
		return nil, err
	}

	query := url.Values{}
	setAccountLookupQuery(query, params, "numeroAgenciaDebito",
		"numeroContaCorrenteDebito", "digitoVerificadorContaCorrenteDebito")
	return get[*GetGRUBatchRequestResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointGRUBatchRequest, id), query),
	)
}

// GetGRUPayment returns a single GRU payment.
func (c *Client) GetGRUPayment(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetGRUPaymentResponse, error) {
	if err := c.requireMTLS(); err != nil {
		return nil, err
	}

	query := url.Values{}
	setAccountLookupQuery(query, params,
		"agencia", "contaCorrente", "digitoVerificador")
	return get[*GetGRUPaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointGRUPayment, id), query),
	)
}
