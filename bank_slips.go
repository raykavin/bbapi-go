package bbapi

import (
	"context"
	"fmt"
	"net/url"
)

// BankSlipEntry represents a bank-slip payment entry.
type BankSlipEntry struct {
	DebitDocumentNumber *int64   `json:"numeroDocumentoDebito,omitempty"`
	BarcodeNumber       string   `json:"numeroCodigoBarras"`
	PaymentDate         int64    `json:"dataPagamento"`
	PaymentValue        float64  `json:"valorPagamento"`
	PaymentDescription  *string  `json:"descricaoPagamento,omitempty"`
	YourDocumentCode    *string  `json:"codigoSeuDocumento,omitempty"`
	OurDocumentCode     *string  `json:"codigoNossoDocumento,omitempty"`
	NominalValue        *float64 `json:"valorNominal,omitempty"`
	DiscountValue       *float64 `json:"valorDesconto,omitempty"`
	LateFeeValue        *float64 `json:"valorMoraMulta,omitempty"`
	PayerTypeCode       *int64   `json:"codigoTipoPagador,omitempty"`
	PayerDocument       *int64   `json:"documentoPagador,omitempty"`
	BeneficiaryTypeCode *int64   `json:"codigoTipoBeneficiario,omitempty"`
	BeneficiaryDocument *int64   `json:"documentoBeneficiario,omitempty"`
	EndorserTypeCode    *int64   `json:"codigoTipoAvalista,omitempty"`
	EndorserDocument    *int64   `json:"documentoAvalista,omitempty"`
}

// CreateBankSlipBatchRequest is the request body for POST /lotes-boletos.
type CreateBankSlipBatchRequest struct {
	RequestNumber          int64           `json:"numeroRequisicao"`
	ContractCode           *int64          `json:"codigoContrato,omitempty"`
	DebitAgencyNumber      *int64          `json:"numeroAgenciaDebito,omitempty"`
	DebitAccountNumber     *int64          `json:"numeroContaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string         `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	Entries                []BankSlipEntry `json:"lancamentos"`
}

// bankSlipBatchEntryBase holds the fields shared between BankSlipBatchEntryResult
// and BankSlipBatchLookupEntry.
type bankSlipBatchEntryBase struct {
	PaymentIdentifier   int64   `json:"codigoIdentificadorPagamento,omitempty"`
	DebitDocumentNumber int64   `json:"numeroDocumentoDebito,omitempty"`
	BarcodeNumber       string  `json:"numeroCodigoBarras,omitempty"`
	PaymentDate         int64   `json:"dataPagamento,omitempty"`
	PaymentValue        float64 `json:"valorPagamento,omitempty"`
	PaymentDescription  string  `json:"descricaoPagamento,omitempty"`
	YourDocumentCode    string  `json:"codigoSeuDocumento,omitempty"`
	OurDocumentCode     string  `json:"codigoNossoDocumento,omitempty"`
	NominalValue        float64 `json:"valorNominal,omitempty"`
	DiscountValue       float64 `json:"valorDesconto,omitempty"`
	LateFeeValue        float64 `json:"valorMoraMulta,omitempty"`
	PayerTypeCode       int64   `json:"codigoTipoPagador,omitempty"`
	PayerDocument       int64   `json:"documentoPagador,omitempty"`
	PayerName           string  `json:"nomePagador,omitempty"`
	BeneficiaryTypeCode int64   `json:"codigoTipoBeneficiario,omitempty"`
	BeneficiaryDocument int64   `json:"documentoBeneficiario,omitempty"`
	BeneficiaryName     string  `json:"nomeBeneficiario,omitempty"`
	EndorserTypeCode    int64   `json:"codigoTipoAvalista,omitempty"`
	EndorserDocument    int64   `json:"documentoAvalista,omitempty"`
	EndorserName        string  `json:"nomeAvalista,omitempty"`
	AcceptanceIndicator string  `json:"indicadorAceite,omitempty"`
}

// BankSlipBatchEntryResult represents a batch entry returned by bank-slip batch creation.
type BankSlipBatchEntryResult struct {
	bankSlipBatchEntryBase
	ErrorCodes []int64 `json:"errorCodes,omitempty"`
}

// BankSlipBatchLookupEntry represents a batch entry returned by bank-slip batch lookups.
type BankSlipBatchLookupEntry struct {
	bankSlipBatchEntryBase
	Errors []int64 `json:"erros,omitempty"`
}

// BankSlipPaymentItem represents a bank-slip payment detail item.
type BankSlipPaymentItem struct {
	Code                string  `json:"codigo,omitempty"`
	OurDocument         string  `json:"nossoDocumento,omitempty"`
	YourDocument        string  `json:"seuDocumento,omitempty"`
	BeneficiaryType     int64   `json:"tipoPessoaBeneficiario,omitempty"`
	BeneficiaryDocument int64   `json:"documentoBeneficiario,omitempty"`
	BeneficiaryName     string  `json:"nomeBeneficiario,omitempty"`
	PayerType           int64   `json:"tipoPessoaPagador,omitempty"`
	PayerDocument       int64   `json:"documentoPagador,omitempty"`
	PayerName           string  `json:"nomePagador,omitempty"`
	EndorserType        int64   `json:"tipoPessoaAvalista,omitempty"`
	EndorserDocument    int64   `json:"documentoAvalista,omitempty"`
	EndorserName        string  `json:"nomeAvalista,omitempty"`
	DueDate             string  `json:"dataVencimento,omitempty"`
	ScheduledDate       string  `json:"dataAgendamento,omitempty"`
	NominalValue        float64 `json:"valorNominal,omitempty"`
	LateFeeValue        float64 `json:"valorMoraMulta,omitempty"`
	DiscountValue       float64 `json:"valorDesconto,omitempty"`
	Text                string  `json:"texto,omitempty"`
}

// BankSlipReturnItem represents a bank-slip return entry.
type BankSlipReturnItem struct {
	ReasonCode  int64   `json:"codigoMotivo,omitempty"`
	ReturnDate  string  `json:"dataDevolucao,omitempty"`
	ReturnValue float64 `json:"valorDevolucao,omitempty"`
}

// CreateBankSlipBatchResponse is the response body for POST /lotes-boletos.
type CreateBankSlipBatchResponse struct {
	RequestNumber   int64                      `json:"numeroRequisicao,omitempty"`
	RequestState    int64                      `json:"estadoRequisicao,omitempty"`
	EntryCount      int64                      `json:"quantidadeLancamentos,omitempty"`
	EntryValue      float64                    `json:"valorLancamentos,omitempty"`
	ValidEntryCount int64                      `json:"quantidadeLancamentosValidos,omitempty"`
	ValidEntryValue float64                    `json:"valorLancamentosValidos,omitempty"`
	Entries         []BankSlipBatchEntryResult `json:"lancamentos,omitempty"`
}

// GetBankSlipBatchRequestResponse is the response body for GET /lotes-boletos/{id}/solicitacao.
type GetBankSlipBatchRequestResponse struct {
	RequestState    int64                      `json:"estadoRequisicao,omitempty"`
	EntryCount      int64                      `json:"quantidadeLancamentos,omitempty"`
	EntryValue      float64                    `json:"valorLancamentos,omitempty"`
	ValidEntryCount int64                      `json:"quantidadeLancamentosValidos,omitempty"`
	ValidEntryValue float64                    `json:"valorLancamentosValidos,omitempty"`
	Entries         []BankSlipBatchLookupEntry `json:"lancamentos,omitempty"`
}

// GetBankSlipPaymentResponse is the response body for GET /boletos/{id}.
type GetBankSlipPaymentResponse struct {
	ID                     int64                 `json:"id"`
	PaymentState           string                `json:"estadoPagamento,omitempty"`
	CreditType             int64                 `json:"tipoCredito,omitempty"`
	DebitAgency            int64                 `json:"agenciaDebito,omitempty"`
	DebitAccount           int64                 `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit string                `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CreditCardStart        int64                 `json:"inicioCartaoCredito,omitempty"`
	CreditCardEnd          int64                 `json:"fimCartaoCredito,omitempty"`
	PaymentDate            string                `json:"dataPagamento,omitempty"`
	PaymentValue           float64               `json:"valorPagamento,omitempty"`
	DebitDocumentNumber    int64                 `json:"documentoDebito,omitempty"`
	PaymentAuthentication  string                `json:"codigoAutenticacaoPagamento,omitempty"`
	PaymentItems           []BankSlipPaymentItem `json:"listaPagamentos,omitempty"`
	ReturnItems            []BankSlipReturnItem  `json:"listaDevolucao,omitempty"`
}

const (
	endpointBankSlipBatch        = "/lotes-boletos"
	endpointBankSlipBatchRequest = "/lotes-boletos/%s/solicitacao"
	endpointBankSlipPayment      = "/boletos/%s"
)

// CreateBankSlipBatch creates a bank-slip batch.
func (c *Client) CreateBankSlipBatch(
	ctx context.Context,
	req *CreateBankSlipBatchRequest,
) (*CreateBankSlipBatchResponse, error) {
	return post[*CreateBankSlipBatchResponse](c, ctx, endpointBankSlipBatch, req)
}

// GetBankSlipBatchRequest returns the request-stage representation of a bank-slip batch.
func (c *Client) GetBankSlipBatchRequest(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetBankSlipBatchRequestResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetBankSlipBatchRequestResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointBankSlipBatchRequest, id), query),
	)
}

// GetBankSlipPayment returns a single bank-slip payment.
func (c *Client) GetBankSlipPayment(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetBankSlipPaymentResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetBankSlipPaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointBankSlipPayment, id), query),
	)
}
