package bbapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	endpointReleasePayments  = "/liberar-pagamentos"
	endpointCancelPayments   = "/cancelar-pagamentos"
	endpointUpdateDates      = "/%s/data-pagamentos"
	endpointReturnedPayments = "/pagamentos"
	endpointPaymentEntries   = "/lancamentos-periodo"
	endpointBarcodePayments  = "/pagamentos-codigo-barras/%s"
)

// ReleasePaymentsRequest is the request body for POST /liberar-pagamentos.
type ReleasePaymentsRequest struct {
	RequestNumber  int64  `json:"numeroRequisicao"`
	FloatIndicator string `json:"indicadorFloat"`
}

// ReleasePaymentsResponse is the response body for POST /liberar-pagamentos.
type ReleasePaymentsResponse struct {
	ReturnMessage string `json:"mensagemRetorno"`
}

// CancelPayment identifies a payment to be cancelled.
type CancelPayment struct {
	PaymentCode string `json:"codigoPagamento"`
}

// CancelPaymentsRequest is the request body for POST /cancelar-pagamentos.
type CancelPaymentsRequest struct {
	DebitAgency            *int64          `json:"agenciaDebito,omitempty"`
	DebitAccount           *int64          `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string         `json:"digitoVerificadorContaCorrente,omitempty"`
	Payments               []CancelPayment `json:"listaPagamentos,omitempty"`
	PaymentContractNumber  *int64          `json:"numeroContratoPagamento,omitempty"`
}

// CancelPaymentResult describes the cancellation result of a payment.
type CancelPaymentResult struct {
	PaymentCode           int64  `json:"codigoPagamento"`
	CancellationIndicator string `json:"indicadorCancelamento,omitempty"`
	PaymentState          string `json:"estadoPagamento,omitempty"`
	CancellationState     string `json:"estadoCancelamento,omitempty"`
}

// CancelPaymentsResponse is the response body for POST /cancelar-pagamentos.
type CancelPaymentsResponse struct {
	Payments []CancelPaymentResult `json:"pagamentos,omitempty"`
}

// UpdatePaymentDatesRequest is the request body for PUT /{id}/data-pagamentos.
type UpdatePaymentDatesRequest struct {
	DebitAgencyNumber      int64  `json:"numeroAgenciaDebito"`
	DebitAccountNumber     int64  `json:"numeroContaCorrenteDebito"`
	DebitAccountCheckDigit string `json:"digitoVerificadorContaCorrenteDebito"`
	ProductCode            *int64 `json:"codigoProduto,omitempty"`
	OriginalPaymentDate    int64  `json:"dataOriginalPagamento"`
	NewPaymentDate         int64  `json:"dataNovoPagamento"`
}

// UpdatePaymentDatesResponse is the response body for PUT /{id}/data-pagamentos.
type UpdatePaymentDatesResponse struct {
	OriginalEntryCount int64 `json:"quantidadeLancamentoOriginal"`
	UpdatedEntryCount  int64 `json:"quantidadeLancamentoAlterado"`
}

// ListReturnedPaymentsParams holds query parameters for GET /pagamentos.
type ListReturnedPaymentsParams struct {
	DebitAgency            *int64
	DebitAccount           *int64
	DebitAccountCheckDigit *string
	PaymentContractNumber  *int64
	StartDate              int64
	EndDate                int64
	PaymentState           *string
	Index                  int64
}

// ReturnedPayment represents a payment returned by GET /pagamentos.
type ReturnedPayment struct {
	PaymentIdentifier         int64   `json:"identificadorPagamento"`
	PaymentType               int64   `json:"tipoPagamento,omitempty"`
	CreditType                int64   `json:"tipoCredito,omitempty"`
	PaymentDate               int64   `json:"dataPagamento,omitempty"`
	CreditCOMPE               int64   `json:"compeCredito,omitempty"`
	CreditISPB                int64   `json:"ispbCredito,omitempty"`
	CreditAgency              int64   `json:"agenciaCredito,omitempty"`
	CreditAccount             int64   `json:"contaCredito,omitempty"`
	CreditAccountCheckDigit   string  `json:"digitoVerificadorContaCredito,omitempty"`
	CreditPaymentAccount      string  `json:"contaPagamentoCredito,omitempty"`
	BeneficiaryType           int64   `json:"tipoBeneficiario,omitempty"`
	TaxID                     int64   `json:"cpfCnpj,omitempty"`
	Name                      string  `json:"nome,omitempty"`
	PaymentValue              float64 `json:"valorPagamento,omitempty"`
	BankSlipBarcode           string  `json:"codigoBarrasBoleto,omitempty"`
	BankSlipValue             float64 `json:"valorBoleto,omitempty"`
	RequestNumber             int64   `json:"numeroRequisicao,omitempty"`
	DebitAgency               int64   `json:"agenciaDebito,omitempty"`
	DebitAccountNumber        int64   `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit    string  `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CardStart                 int64   `json:"inicioCartao,omitempty"`
	CardEnd                   int64   `json:"fimCartao,omitempty"`
	DebitDocumentNumber       int64   `json:"documentoDebito,omitempty"`
	ReturnDate                int64   `json:"dataDevolucao,omitempty"`
	ReturnValue               float64 `json:"valorDevolucao,omitempty"`
	ReturnCode                int64   `json:"codigoDevolucao,omitempty"`
	ReturnSequence            int64   `json:"sequenciaDevolucao,omitempty"`
	AccountType               int64   `json:"tipoConta,omitempty"`
	InstantPaymentDescription string  `json:"descricaoPagamentoInstantaneo,omitempty"`
	IdentificationMode        int64   `json:"formaIdentificacao,omitempty"`
}

// ListReturnedPaymentsResponse is the response body for GET /pagamentos.
type ListReturnedPaymentsResponse struct {
	Index            int64             `json:"indice"`
	TotalRecordCount int64             `json:"quantidadeTotalRegistros,omitempty"`
	RecordCount      int64             `json:"quantidadeRegistros,omitempty"`
	Payments         []ReturnedPayment `json:"pagamentos,omitempty"`
}

// ListPaymentEntriesParams holds query parameters for GET /lancamentos-periodo.
type ListPaymentEntriesParams struct {
	ClientAgreementCode    *int64
	DebitAgencyNumber      int64
	DebitAccountNumber     int64
	DebitAccountCheckDigit string
	RequestSentStartDate   int64
	RequestSentEndDate     *int64
	PaymentStateCode       *int64
	ProductCode            *int64
	SearchPosition         *int64
}

// PaymentEntry represents a payment entry returned by GET /lancamentos-periodo.
type PaymentEntry struct {
	RequestNumber             int64   `json:"numeroRequisicaoPagamento"`
	PaymentStateText          string  `json:"textoEstadoPagamento,omitempty"`
	PaymentIdentifier         int64   `json:"codigoIdentificadorDoPagamento,omitempty"`
	BeneficiaryName           string  `json:"nomeDoFavorecido,omitempty"` // fixed: had a spurious "." before the comma
	PersonTypeCode            int64   `json:"codigoDoTipoDePessoa,omitempty"`
	TaxID                     int64   `json:"numeroCPFouCNPJ,omitempty"`
	PaymentDate               int64   `json:"dataPagamento,omitempty"`
	PaymentValue              float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber       int64   `json:"numeroDocumentoDebito,omitempty"`
	CreditDocumentNumber      int64   `json:"numeroDocumentoCredito,omitempty"`
	CreditTypeCode            int64   `json:"codigoFormaCredito,omitempty"`
	PaymentAuthenticationCode string  `json:"codigoAutenticacaoPagamento,omitempty"`
	DebitDate                 int64   `json:"dataDebito,omitempty"`
	PaymentTypeCode           int64   `json:"codigoTipoPagamento,omitempty"`
}

// ListPaymentEntriesResponse is the response body for GET /lancamentos-periodo.
type ListPaymentEntriesResponse struct {
	SearchPosition     int64          `json:"numeroDaposicaoDePesquisa"`
	TotalOccurrences   int64          `json:"quantidadeOcorrenciasTotal,omitempty"`
	IndexedOccurrences int64          `json:"quantidadeOcorrenciasTabeladas,omitempty"`
	Entries            []PaymentEntry `json:"listaMovimento,omitempty"`
}

// BarcodeLinkedPayment represents a barcode-linked payment.
type BarcodeLinkedPayment struct {
	PaymentIdentifier      int64   `json:"identificadorPagamento"`
	PaymentAuthentication  string  `json:"autenticacaoPagamento,omitempty"`
	DebitAgency            int64   `json:"agenciaDebito,omitempty"`
	DebitAccountNumber     int64   `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit string  `json:"digitoVerificadorContaCorrente,omitempty"`
	RequestNumber          int64   `json:"numeroRequisicao,omitempty"`
	DebitDocumentNumber    int64   `json:"documentoDebito,omitempty"`
	ScheduledDate          int64   `json:"dataAgendamento,omitempty"`
	PaymentDate            int64   `json:"dataPagamento,omitempty"`
	PaymentValue           float64 `json:"valorPagamento,omitempty"`
	BeneficiaryName        string  `json:"nomeBeneficiario,omitempty"`
	PaymentState           string  `json:"estadoPagamento,omitempty"`
	PaymentDescription     string  `json:"descricaoPagamento,omitempty"`
	Errors                 []int64 `json:"erros,omitempty"`
}

// BarcodePaymentsResponse is the response body for GET /pagamentos-codigo-barras/{id}.
type BarcodePaymentsResponse struct {
	Payments []BarcodeLinkedPayment `json:"pagamentos,omitempty"`
}

// ReleasePayments releases a payment batch.
func (c *Client) ReleasePayments(ctx context.Context, req *ReleasePaymentsRequest) (*ReleasePaymentsResponse, error) {
	return post[*ReleasePaymentsResponse](c, ctx, endpointReleasePayments, req)
}

// CancelPayments requests payment cancellation.
func (c *Client) CancelPayments(ctx context.Context, req *CancelPaymentsRequest) (*CancelPaymentsResponse, error) {
	return post[*CancelPaymentsResponse](c, ctx, endpointCancelPayments, req)
}

// UpdatePaymentDates updates scheduled payment dates inside a batch.
func (c *Client) UpdatePaymentDates(
	ctx context.Context,
	id string,
	req *UpdatePaymentDatesRequest,
) (*UpdatePaymentDatesResponse, error) {
	return put[*UpdatePaymentDatesResponse](c, ctx, fmt.Sprintf(endpointUpdateDates, id), req)
}

// ListReturnedPayments lists returned or reversed payments.
func (c *Client) ListReturnedPayments(
	ctx context.Context,
	params *ListReturnedPaymentsParams,
) (*ListReturnedPaymentsResponse, error) {
	query := url.Values{}
	if params != nil {
		setInt64(query, "agenciaDebito", params.DebitAgency)
		setInt64(query, "contaCorrenteDebito", params.DebitAccount)
		setString(query, "digitoVerificadorContaCorrente", params.DebitAccountCheckDigit)
		setInt64(query, "numeroContratoPagamento", params.PaymentContractNumber)
		query.Set("dataInicio", strconv.FormatInt(params.StartDate, 10))
		query.Set("dataFim", strconv.FormatInt(params.EndDate, 10))
		query.Set("indice", strconv.FormatInt(params.Index, 10))
		setString(query, "estadoPagamento", params.PaymentState)
	}
	return get[*ListReturnedPaymentsResponse](c, ctx, buildPath(endpointReturnedPayments, query))
}

// ListPaymentEntries lists payment entries by period.
func (c *Client) ListPaymentEntries(
	ctx context.Context,
	params *ListPaymentEntriesParams,
) (*ListPaymentEntriesResponse, error) {
	query := url.Values{}
	if params != nil {
		setInt64(query, "codigoClienteConveniado", params.ClientAgreementCode)
		query.Set("numeroAgenciaDebito", strconv.FormatInt(params.DebitAgencyNumber, 10))
		query.Set("numeroContaCorrenteDebito", strconv.FormatInt(params.DebitAccountNumber, 10))
		query.Set("digitoVerificadorContaCorrenteDebito", params.DebitAccountCheckDigit)
		query.Set("dataInicialdeEnviodaRequisicao", strconv.FormatInt(params.RequestSentStartDate, 10))
		setInt64(query, "dataFinaldeEnviodaRequisicao", params.RequestSentEndDate)
		setInt64(query, "codigodoEstadodoPagamento", params.PaymentStateCode)
		setInt64(query, "codigoProduto", params.ProductCode)
		setInt64(query, "numeroDaPosicaoDePesquisa", params.SearchPosition)
	}
	return get[*ListPaymentEntriesResponse](c, ctx, buildPath(endpointPaymentEntries, query))
}

// GetBarcodePayments returns payments linked to a barcode.
func (c *Client) GetBarcodePayments(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*BarcodePaymentsResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agenciaDebito", "contaCorrenteDebito", "digitoVerificadorContaCorrente")
	return get[*BarcodePaymentsResponse](c, ctx, buildPath(fmt.Sprintf(endpointBarcodePayments, id), query))
}
