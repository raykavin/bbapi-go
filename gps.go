package bbapi

import (
	"context"
	"fmt"
	"net/url"
)

const (
	endpointGPSBatch        = "/lotes-gps"
	endpointGPSBatchRequest = "/lotes-gps/%s/solicitacao"
	endpointGPSPayment      = "/gps/%s"
)

// GPSEntry represents a GPS payment entry.
type GPSEntry struct {
	PaymentDate             int64    `json:"dataPagamento"`
	PaymentValue            float64  `json:"valorPagamento"`
	DebitDocumentNumber     *int64   `json:"numeroDocumentoDebito,omitempty"`
	YourDocumentCode        *string  `json:"codigoSeuDocumento,omitempty"`
	PaymentDescription      *string  `json:"textoDescricaoPagamento,omitempty"`
	TaxRevenueCode          *int64   `json:"codigoReceitaTributoGuiaPrevidenciaSocial,omitempty"`
	TaxpayerTypeCode        *int64   `json:"codigoTipoContribuinteGuiaPrevidenciaSocial,omitempty"`
	TaxpayerIdentification  *int64   `json:"numeroIdentificacaoContribuinteGuiaPrevidenciaSocial,omitempty"`
	TaxIdentifierCode       *string  `json:"codigoIdentificadorTributoGuiaPrevidenciaSocial,omitempty"`
	CompetenceMonthYear     *int64   `json:"mesAnoCompetenciaGuiaPrevidenciaSocial,omitempty"`
	INSSValue               *float64 `json:"valorPrevistoInstNacSeguridadeSocialGuiaPrevidenciaSocial,omitempty"`
	OtherEntriesValue       *float64 `json:"valorOutroEntradaGuiaPrevidenciaSocial,omitempty"`
	MonetaryAdjustmentValue *float64 `json:"valorAtualizacaoMonetarioGuiaPrevidenciaSocial,omitempty"`
}

// CreateGPSBatchRequest is the request body for POST /lotes-gps.
type CreateGPSBatchRequest struct {
	RequestNumber          int64      `json:"numeroRequisicao"`
	ContractCode           *int64     `json:"codigoContrato,omitempty"`
	DebitAgencyNumber      *int64     `json:"numeroAgenciaDebito,omitempty"`
	DebitAccountNumber     *int64     `json:"numeroContaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string    `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	Entries                []GPSEntry `json:"lancamentos"`
}

// gpsBatchEntryBase holds the fields shared between GPSBatchEntryResult
// and GPSBatchRequestEntry.
type gpsBatchEntryBase struct {
	PaymentIdentifier       int64   `json:"codigoIdentificadorPagamento,omitempty"`
	ConventionName          string  `json:"nomeConvenente,omitempty"`
	PaymentDate             int64   `json:"dataPagamento,omitempty"`
	PaymentValue            float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber     int64   `json:"numeroDocumentoDebito,omitempty"`
	YourDocumentCode        string  `json:"codigoSeuDocumento,omitempty"`
	TaxRevenueCode          int64   `json:"codigoReceitaTributoGuiaPrevidenciaSocial,omitempty"`
	TaxpayerTypeCode        int64   `json:"codigoTipoContribuinteGuiaPrevidenciaSocial,omitempty"`
	TaxpayerIdentification  int64   `json:"numeroIdentificacaoContribuinteGuiaPrevidenciaSocial,omitempty"`
	TaxIdentifierCode       string  `json:"codigoIdentificadorTributoGuiaPrevidenciaSocial,omitempty"`
	CompetenceMonthYear     int64   `json:"mesAnoCompetenciaGuiaPrevidenciaSocial,omitempty"`
	INSSValue               float64 `json:"valorPrevistoInstNacSeguridadeSocialGuiaPrevidenciaSocial,omitempty"`
	OtherEntriesValue       float64 `json:"valorOutroEntradaGuiaPrevidenciaSocial,omitempty"`
	MonetaryAdjustmentValue float64 `json:"valorAtualizacaoMonetarioGuiaPrevidenciaSocial,omitempty"`
	AcceptanceIndicator     string  `json:"indicadorMovimentoAceito,omitempty"`
	Errors                  []int64 `json:"erros,omitempty"`
}

// GPSBatchEntryResult represents a GPS batch entry returned by batch creation.
type GPSBatchEntryResult struct {
	gpsBatchEntryBase
}

// GPSBatchRequestEntry represents a GPS batch entry returned by request lookups.
type GPSBatchRequestEntry struct {
	gpsBatchEntryBase
	PaymentDescription string `json:"textoDescricaoPagamento,omitempty"`
	GPSText            string `json:"textoGuiaPrevidenciaSocial,omitempty"`
}

// GPSPaymentItem represents a GPS payment detail item.
type GPSPaymentItem struct {
	Code                    int64   `json:"codigo,omitempty"`
	TaxpayerType            int64   `json:"tipoContribuinte,omitempty"`
	TaxpayerIdentification  int64   `json:"identificacaoContribuinte,omitempty"`
	GPSIdentification       string  `json:"identificacaoGPS,omitempty"`
	CompetenceMonthYear     int64   `json:"mesAnoCompetencia,omitempty"`
	INSSValue               float64 `json:"valorINSS,omitempty"`
	MonetaryAdjustmentValue float64 `json:"valorAtualizacaoMonetaria,omitempty"`
	Text                    string  `json:"texto,omitempty"`
}

// GPSReturnItem represents a GPS return entry.
type GPSReturnItem struct {
	ReasonCode int64 `json:"codigoMotivo,omitempty"`
}

// CreateGPSBatchResponse is the response body for POST /lotes-gps.
type CreateGPSBatchResponse struct {
	RequestNumber    int64                 `json:"numeroRequisicao,omitempty"`
	RequestStateCode int64                 `json:"codigoEstadoRequisicao,omitempty"`
	TotalEntryCount  int64                 `json:"quantidadeTotalLancamento,omitempty"`
	TotalEntryValue  float64               `json:"valorTotalLancamento,omitempty"`
	ValidTotalCount  int64                 `json:"quantidadeTotalValido,omitempty"`
	ValidEntryValue  float64               `json:"valorLancamentosValidos,omitempty"`
	Entries          []GPSBatchEntryResult `json:"lancamentos,omitempty"`
}

// GetGPSBatchRequestResponse is the response body for GET /lotes-gps/{id}/solicitacao.
type GetGPSBatchRequestResponse struct {
	RequestNumber    int64                  `json:"numeroRequisicao,omitempty"`
	RequestStateCode int64                  `json:"codigoEstadoRequisicao,omitempty"`
	TotalEntryCount  int64                  `json:"quantidadeTotalLancamento,omitempty"`
	TotalEntryValue  float64                `json:"valorTotalLancamento,omitempty"`
	ValidTotalCount  int64                  `json:"quantidadeTotalValido,omitempty"`
	ValidEntryValue  float64                `json:"valorLancamentosValidos,omitempty"`
	Entries          []GPSBatchRequestEntry `json:"lancamentos,omitempty"`
}

// GetGPSPaymentResponse is the response body for GET /gps/{id}.
type GetGPSPaymentResponse struct {
	ID                  int64            `json:"id"`
	PaymentState        string           `json:"estadoPagamento,omitempty"`
	DebitAgency         int64            `json:"agenciaDebito,omitempty"`
	DebitAccount        int64            `json:"contaCorrenteDebito,omitempty"`
	DebitAccountDigit   string           `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CreditCardStart     int64            `json:"inicioCartaoCredito,omitempty"`
	CreditCardEnd       int64            `json:"fimCartaoCredito,omitempty"`
	PaymentDate         int64            `json:"dataPagamento,omitempty"`
	PaymentValue        float64          `json:"valorPagamento,omitempty"`
	DebitDocumentNumber int64            `json:"documentoDebito,omitempty"`
	AuthenticationCode  string           `json:"codigoAutenticacaoPagamento,omitempty"`
	PaymentItems        []GPSPaymentItem `json:"listaPagamentos,omitempty"`
	ReturnItems         []GPSReturnItem  `json:"listaDevolucao,omitempty"`
}

// CreateGPSBatch creates a GPS batch.
func (c *Client) CreateGPSBatch(
	ctx context.Context,
	req *CreateGPSBatchRequest,
) (*CreateGPSBatchResponse, error) {
	return post[*CreateGPSBatchResponse](c, ctx, endpointGPSBatch, req)
}

// GetGPSBatchRequest returns the request-stage representation of a GPS batch.
func (c *Client) GetGPSBatchRequest(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetGPSBatchRequestResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "numeroAgenciaDebito",
		"numeroContaCorrenteDebito", "digitoVerificadorContaCorrenteDebito")
	return get[*GetGPSBatchRequestResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointGPSBatchRequest, id), query),
	)
}

// GetGPSPayment returns a single GPS payment.
func (c *Client) GetGPSPayment(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetGPSPaymentResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetGPSPaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointGPSPayment, id), query),
	)
}
