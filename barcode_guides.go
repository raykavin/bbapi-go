package bbapi

import (
	"context"
	"fmt"
	"net/url"
)

const (
	endpointBarcodeGuideBatch        = "/lotes-guias-codigo-barras"
	endpointBarcodeGuideBatchRequest = "/lotes-guias-codigo-barras/%s/solicitacao"
	endpointBarcodeGuidePayment      = "/guias-codigo-barras/%s"
)

// BarcodeGuideEntry represents a barcode-guide payment entry.
type BarcodeGuideEntry struct {
	Barcode             string  `json:"codigoBarras"`
	PaymentDate         int64   `json:"dataPagamento"`
	PaymentValue        float64 `json:"valorPagamento"`
	DebitDocumentNumber *int64  `json:"numeroDocumentoDebito,omitempty"`
	YourDocumentCode    *string `json:"codigoSeuDocumento,omitempty"`
	PaymentDescription  *string `json:"descricaoPagamento,omitempty"`
}

// CreateBarcodeGuideBatchRequest is the request body for POST /lotes-guias-codigo-barras.
type CreateBarcodeGuideBatchRequest struct {
	RequestNumber          int64               `json:"numeroRequisicao"`
	ContractCode           *int64              `json:"codigoContrato,omitempty"`
	DebitAgencyNumber      *int64              `json:"numeroAgenciaDebito,omitempty"`
	DebitAccountNumber     *int64              `json:"numeroContaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string             `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	Entries                []BarcodeGuideEntry `json:"lancamentos"`
}

// BarcodeGuideBatchEntryResult represents a batch entry returned by barcode-guide batch creation.
type BarcodeGuideBatchEntryResult struct {
	PaymentIdentifier   string  `json:"codigoIdentificadorPagamento,omitempty"`
	BeneficiaryName     string  `json:"nomeBeneficiario,omitempty"`
	Barcode             string  `json:"codigoBarras,omitempty"`
	PaymentDate         string  `json:"dataPagamento,omitempty"`
	PaymentValue        float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64   `json:"numeroDocumentoDebito,omitempty"`
	YourDocumentCode    string  `json:"codigoSeuDocumento,omitempty"`
	PaymentDescription  string  `json:"descricaoPagamento,omitempty"`
	AcceptanceIndicator string  `json:"indicadorAceite,omitempty"`
	Errors              []int64 `json:"erros,omitempty"`
}

// BarcodeGuideBatchPayment represents a payment returned by barcode-guide batch lookups.
type BarcodeGuideBatchPayment struct {
	PaymentCode         string  `json:"codigoPagamento,omitempty"`
	BeneficiaryName     string  `json:"nomeBeneficiario,omitempty"`
	Barcode             string  `json:"codigoBarras,omitempty"`
	PaymentDate         string  `json:"dataPagamento,omitempty"`
	PaymentValue        float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64   `json:"documentoDebito,omitempty"`
	YourDocumentCode    string  `json:"codigoSeuDocumento,omitempty"`
	PaymentDescription  string  `json:"descricaoPagamento,omitempty"`
	AcceptanceIndicator string  `json:"indicadorAceite,omitempty"`
	Errors              []int64 `json:"erros,omitempty"`
}

// BarcodeGuidePaymentItem represents a barcode-guide payment detail item.
type BarcodeGuidePaymentItem struct {
	Code         string `json:"codigo,omitempty"`
	ReceiverName string `json:"nomeRecebedor,omitempty"`
	YourNumber   string `json:"seuNumero,omitempty"`
	Text         string `json:"texto,omitempty"`
}

// BarcodeGuideReturnItem represents a barcode-guide return entry.
type BarcodeGuideReturnItem struct {
	ReasonCode int64 `json:"codigoMotivo,omitempty"`
}

// CreateBarcodeGuideBatchResponse is the response body for POST /lotes-guias-codigo-barras.
type CreateBarcodeGuideBatchResponse struct {
	RequestNumber   int64                          `json:"numeroRequisicao,omitempty"`
	StateCode       int64                          `json:"codigoEstado,omitempty"`
	EntryCount      int64                          `json:"quantidadeLancamentos,omitempty"`
	EntryValue      float64                        `json:"valorLancamentos,omitempty"`
	ValidEntryCount int64                          `json:"quantidadeLancamentosValidos,omitempty"`
	ValidEntryValue float64                        `json:"valorLancamentosValidos,omitempty"`
	Entries         []BarcodeGuideBatchEntryResult `json:"lancamentos,omitempty"`
}

// GetBarcodeGuideBatchRequestResponse is the response body for GET /lotes-guias-codigo-barras/{id}/solicitacao.
type GetBarcodeGuideBatchRequestResponse struct {
	RequestNumber     int64                      `json:"numeroRequisicao,omitempty"`
	RequestState      int64                      `json:"estadoRequisicao,omitempty"`
	PaymentCount      int64                      `json:"quantidadePagamentos,omitempty"`
	PaymentValue      float64                    `json:"valorPagamentos,omitempty"`
	ValidPaymentCount int64                      `json:"quantidadePagamentosValidos,omitempty"`
	ValidPaymentValue float64                    `json:"valorPagamentosValidos,omitempty"`
	Payments          []BarcodeGuideBatchPayment `json:"pagamentos,omitempty"`
}

// GetBarcodeGuidePaymentResponse is the response body for GET /guias-codigo-barras/{id}.
type GetBarcodeGuidePaymentResponse struct {
	ID                  int64                     `json:"id"`
	PaymentState        string                    `json:"estadoPagamento,omitempty"`
	DebitAgency         int64                     `json:"agenciaDebito,omitempty"`
	DebitAccount        int64                     `json:"contaCorrenteDebito,omitempty"`
	DebitAccountDigit   string                    `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CreditCardStart     int64                     `json:"inicioCartaoCredito,omitempty"`
	CreditCardEnd       int64                     `json:"fimCartaoCredito,omitempty"`
	PaymentDate         string                    `json:"dataPagamento,omitempty"`
	PaymentValue        float64                   `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64                     `json:"documentoDebito,omitempty"`
	AuthenticationCode  string                    `json:"codigoAutenticacaoPagamento,omitempty"`
	PaymentItems        []BarcodeGuidePaymentItem `json:"listaPagamentos,omitempty"`
	ReturnItems         []BarcodeGuideReturnItem  `json:"listaDevolucao,omitempty"`
}

// CreateBarcodeGuideBatch creates a barcode-guide batch.
func (c *Client) CreateBarcodeGuideBatch(
	ctx context.Context,
	req *CreateBarcodeGuideBatchRequest,
) (*CreateBarcodeGuideBatchResponse, error) {
	return post[*CreateBarcodeGuideBatchResponse](c, ctx, endpointBarcodeGuideBatch, req)
}

// GetBarcodeGuideBatchRequest returns the request-stage representation of a barcode-guide batch.
func (c *Client) GetBarcodeGuideBatchRequest(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetBarcodeGuideBatchRequestResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetBarcodeGuideBatchRequestResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointBarcodeGuideBatchRequest, id), query),
	)
}

// GetBarcodeGuidePayment returns a single barcode-guide payment.
func (c *Client) GetBarcodeGuidePayment(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetBarcodeGuidePaymentResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetBarcodeGuidePaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointBarcodeGuidePayment, id), query),
	)
}
