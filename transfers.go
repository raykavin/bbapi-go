package bbapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	endpointTransferBatches         = "/lotes-transferencias"
	endpointTransferPayment         = "/transferencias/%s"
	endpointBatchRequest            = "/%s/solicitacao"
	endpointBatch                   = "/%s"
	endpointBeneficiaryTransfers    = "/beneficiarios/%s/transferencias"
	endpointPixTransferBatches      = "/lotes-transferencias-pix"
	endpointPixTransferBatchRequest = "/lotes-transferencias-pix/%s/solicitacao"
	endpointPixPayment              = "/pix/%s"
)

// ListTransferBatchesParams holds query parameters for GET /lotes-transferencias.
type ListTransferBatchesParams struct {
	PaymentContractNumber  *int64
	DebitAgency            *int64
	DebitAccount           *int64
	DebitAccountCheckDigit *string
	StartDate              *int64
	EndDate                *int64
	PaymentType            *int64
	RequestState           *int64
	Index                  *int64
}

// TransferBatchSummary represents a transfer batch returned by GET /lotes-transferencias.
type TransferBatchSummary struct {
	RequestNumber      int64   `json:"numeroRequisicao,omitempty"`
	RequestState       int64   `json:"estadoRequisicao,omitempty"`
	DebitAgency        int64   `json:"agenciaDebito,omitempty"`
	DebitAccount       int64   `json:"contaCorrenteDebito,omitempty"`
	DebitAccountDigit  string  `json:"digitoVerificadorContaCorrente,omitempty"`
	RequestDate        int64   `json:"dataRequisicao,omitempty"`
	PaymentType        int64   `json:"tipoPagamento,omitempty"`
	RequesterID        string  `json:"identificacaoRequisitante,omitempty"`
	TransferCount      int64   `json:"quantidadeTransferencias,omitempty"`
	TransferValue      float64 `json:"totalTransferencias,omitempty"`
	ValidTransferCount int64   `json:"quantidadeTransferenciasValidas,omitempty"`
	ValidTransferValue float64 `json:"totalTransferenciasValidas,omitempty"`
}

// ListTransferBatchesResponse is the response body for GET /lotes-transferencias.
type ListTransferBatchesResponse struct {
	Index     int64                  `json:"indice"`
	Transfers []TransferBatchSummary `json:"transferencias,omitempty"`
}

// CreateTransferBatchRequest is the request body for POST /lotes-transferencias.
type CreateTransferBatchRequest struct {
	RequestNumber          int64      `json:"numeroRequisicao"`
	PaymentContractNumber  *int64     `json:"numeroContratoPagamento,omitempty"`
	DebitAgency            *int64     `json:"agenciaDebito,omitempty"`
	DebitAccount           *int64     `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string    `json:"digitoVerificadorContaCorrente,omitempty"`
	PaymentType            int64      `json:"tipoPagamento"`
	Transfers              []Transfer `json:"listaTransferencias"`
}

// Transfer represents a single transfer entry.
type Transfer struct {
	COMPENumber             *int64  `json:"numeroCOMPE,omitempty"`
	ISPBNumber              *int64  `json:"numeroISPB,omitempty"`
	CreditAgency            *int64  `json:"agenciaCredito,omitempty"`
	CreditAccount           *int64  `json:"contaCorrenteCredito,omitempty"`
	CreditAccountCheckDigit *string `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditPaymentAccount    *string `json:"contaPagamentoCredito,omitempty"`
	BeneficiaryCPF          *int64  `json:"cpfBeneficiario,omitempty"`
	BeneficiaryCNPJ         *int64  `json:"cnpjBeneficiario,omitempty"`
	TransferDate            int64   `json:"dataTransferencia"`
	TransferValue           float64 `json:"valorTransferencia"`
	DebitDocumentNumber     *int64  `json:"documentoDebito,omitempty"`
	CreditDocumentNumber    *int64  `json:"documentoCredito,omitempty"`
	DOCPurposeCode          *int64  `json:"codigoFinalidadeDOC,omitempty"`
	TEDPurposeCode          *int64  `json:"codigoFinalidadeTED,omitempty"`
	JudicialDepositNumber   *string `json:"numeroDepositoJudicial,omitempty"`
	TransferDescription     *string `json:"descricaoTransferencia,omitempty"`
}

// TransferResult represents a transfer returned in request responses.
type TransferResult struct {
	TransferIdentifier      int64   `json:"identificadorTransferencia,omitempty"`
	CreditType              int64   `json:"tipoCredito,omitempty"`
	COMPENumber             int64   `json:"numeroCOMPE,omitempty"`
	ISPBNumber              int64   `json:"numeroISPB,omitempty"`
	CreditAgency            int64   `json:"agenciaCredito,omitempty"`
	CreditAccount           int64   `json:"contaCorrenteCredito,omitempty"`
	CreditAccountCheckDigit string  `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditPaymentAccount    string  `json:"contaPagamentoCredito,omitempty"`
	BeneficiaryCPF          int64   `json:"cpfBeneficiario,omitempty"`
	BeneficiaryCNPJ         int64   `json:"cnpjBeneficiario,omitempty"`
	TransferDate            int64   `json:"dataTransferencia,omitempty"`
	TransferValue           float64 `json:"valorTransferencia,omitempty"`
	DebitDocumentNumber     int64   `json:"documentoDebito,omitempty"`
	CreditDocumentNumber    int64   `json:"documentoCredito,omitempty"`
	JudicialDepositNumber   string  `json:"numeroDepositoJudicial,omitempty"`
	TransferDescription     string  `json:"descricaoTransferencia,omitempty"`
	AcceptanceIndicator     string  `json:"indicadorAceite,omitempty"`
	DOCPurposeCode          string  `json:"codigoFinalidadeDOC,omitempty"`
	TEDPurposeCode          string  `json:"codigoFinalidadeTED,omitempty"`
	Errors                  []int64 `json:"erros,omitempty"`
}

// TransferPaymentItem represents a transfer item returned by a payment lookup.
type TransferPaymentItem struct {
	COMPENumber             int64  `json:"numeroCOMPE,omitempty"`
	ISPBNumber              int64  `json:"numeroISPB,omitempty"`
	CreditAgency            int64  `json:"agenciaCredito,omitempty"`
	CreditAccount           int64  `json:"contaCorrenteCredito,omitempty"`
	CreditAccountCheckDigit string `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditPaymentAccount    string `json:"numeroContaCredito,omitempty"`
	BeneficiaryType         int64  `json:"tipoBeneficiario,omitempty"`
	BeneficiaryDocument     int64  `json:"cpfCnpjBeneficiario,omitempty"`
	BeneficiaryName         string `json:"nomeBeneficiario,omitempty"`
	CreditDocumentNumber    int64  `json:"documentoCredito,omitempty"`
	Text                    string `json:"texto,omitempty"`
}

// TransferReturnItem represents a transfer return entry.
type TransferReturnItem struct {
	ReasonCode  int64   `json:"codigoMotivo,omitempty"`
	ReturnDate  string  `json:"dataDevolucao,omitempty"`
	ReturnValue float64 `json:"valorDevolucao,omitempty"`
}

// TransferBatchRequestPayment represents a payment returned by a batch request lookup.
type TransferBatchRequestPayment struct {
	PaymentIdentifier       int64   `json:"identificadorPagamento,omitempty"`
	COMPENumber             int64   `json:"numeroCOMPE,omitempty"`
	ISPBNumber              int64   `json:"numeroISPB,omitempty"`
	CreditAgency            int64   `json:"agenciaCredito,omitempty"`
	CreditAccount           int64   `json:"contaCorrenteCredito,omitempty"`
	CreditAccountCheckDigit string  `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditPaymentAccount    string  `json:"contaPagamentoCredito,omitempty"`
	BeneficiaryCPF          int64   `json:"cpfBeneficiario,omitempty"`
	BeneficiaryCNPJ         int64   `json:"cnpjBeneficiario,omitempty"`
	PaymentDate             int64   `json:"dataPagamento,omitempty"`
	PaymentValue            float64 `json:"valorPagamento,omitempty"`
	DebitDocumentNumber     int64   `json:"documentoDebito,omitempty"`
	CreditDocumentNumber    int64   `json:"documentoCredito,omitempty"`
	CreditType              int64   `json:"tipoCredito,omitempty"`
	DOCPurposeCode          int64   `json:"codigoFinalidadeDOC,omitempty"`
	TEDPurposeCode          int64   `json:"codigoFinalidadeTED,omitempty"`
	JudicialDepositNumber   string  `json:"numeroDepositoJudicial,omitempty"`
	PaymentDescription      string  `json:"descricaoPagamento,omitempty"`
	AcceptanceIndicator     string  `json:"indicadorAceite,omitempty"`
	Errors                  []int64 `json:"erros,omitempty"`
}

// TransferBatchPayment represents a payment returned by a batch lookup.
type TransferBatchPayment struct {
	PaymentIdentifier   int64   `json:"identificadorPagamento,omitempty"`
	PaymentDate         int64   `json:"dataPagamento,omitempty"`
	PaymentValue        float64 `json:"valorPagamento,omitempty"`
	CreditType          int64   `json:"tipoCredito,omitempty"`
	BeneficiaryType     int64   `json:"tipoBeneficiario,omitempty"`
	BeneficiaryDocument int64   `json:"cpfCnpjBeneficiario,omitempty"`
	BeneficiaryName     string  `json:"nomeBeneficiario,omitempty"`
	PaymentState        string  `json:"estadoPagamento,omitempty"`
	PaymentDescription  string  `json:"descricaoPagamento,omitempty"`
	Errors              []int64 `json:"erros,omitempty"`
}

// BeneficiaryTransfer represents a beneficiary transfer entry.
type BeneficiaryTransfer struct {
	Identifier              int64   `json:"identificador,omitempty"`
	PaymentState            string  `json:"estadoPagamento,omitempty"`
	PaymentType             int64   `json:"tipoPagamento,omitempty"`
	CreditType              int64   `json:"tipoCredito,omitempty"`
	TransferDate            int64   `json:"dataTransferencia,omitempty"`
	TransferValue           float64 `json:"valorTransferencia,omitempty"`
	DebitDocumentNumber     int64   `json:"documentoDebito,omitempty"`
	COMPENumber             int64   `json:"numeroCOMPE,omitempty"`
	ISPBNumber              int64   `json:"numeroISPB,omitempty"`
	CreditAgency            int64   `json:"agenciaCredito,omitempty"`
	CreditAccount           int64   `json:"contaCorrenteCredito,omitempty"`
	CreditAccountCheckDigit string  `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditPaymentAccount    string  `json:"contaPagamentoCredito,omitempty"`
	BeneficiaryType         int64   `json:"tipoBeneficiario,omitempty"`
	BeneficiaryDocument     int64   `json:"cpfCnpjBeneficiario,omitempty"`
	BeneficiaryName         string  `json:"nomeBeneficiario,omitempty"`
	AuthenticationCode      string  `json:"codigoAutenticacaoPagamento,omitempty"`
	DOCPurposeCode          string  `json:"codigoFinalidadeDOC,omitempty"`
	TEDPurposeCode          string  `json:"codigoFinalidadeTED,omitempty"`
	JudicialDepositNumber   string  `json:"numeroDepositoJudicial,omitempty"`
	RequestNumber           int64   `json:"numeroRequisicao,omitempty"`
	PaymentFileNumber       int64   `json:"numeroArquivoPagamento,omitempty"`
	DebitAgency             int64   `json:"agenciaDebito,omitempty"`
	DebitAccount            int64   `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit  string  `json:"digitoVerificadorContaCorrenteDebito,omitempty"`
	CardStart               int64   `json:"inicioCartao,omitempty"`
	CardEnd                 int64   `json:"fimCartao,omitempty"`
	TransferDescription     string  `json:"descricaoTransferencia,omitempty"`
	TransmissionType        int64   `json:"formaTransmissao,omitempty"`
}

// PixTransferBatchItem represents a Pix transfer returned by batch responses.
type PixTransferBatchItem struct {
	PaymentIdentifier         int64   `json:"identificadorPagamento,omitempty"`
	Date                      int64   `json:"data,omitempty"`
	Value                     float64 `json:"valor,omitempty"`
	DebitDocumentNumber       int64   `json:"documentoDebito,omitempty"`
	CreditDocumentNumber      int64   `json:"documentoCredito,omitempty"`
	PaymentDescription        string  `json:"descricaoPagamento,omitempty"`
	InstantPaymentDescription string  `json:"descricaoPagamentoInstantaneo,omitempty"`
	IdentificationMode        int64   `json:"formaIdentificacao,omitempty"`
	PhoneAreaCode             int64   `json:"dddTelefone,omitempty"`
	PhoneNumber               int64   `json:"telefone,omitempty"`
	Email                     string  `json:"email,omitempty"`
	BeneficiaryCPF            int64   `json:"cpf,omitempty"`
	BeneficiaryCNPJ           int64   `json:"cnpj,omitempty"`
	RandomIdentifier          string  `json:"identificacaoAleatoria,omitempty"`
	COMPENumber               int64   `json:"numeroCOMPE,omitempty"`
	ISPBNumber                int64   `json:"numeroISPB,omitempty"`
	AccountType               int64   `json:"tipoConta,omitempty"`
	Agency                    int64   `json:"agencia,omitempty"`
	Account                   int64   `json:"conta,omitempty"`
	AccountCheckDigit         string  `json:"digitoVerificadorConta,omitempty"`
	PaymentAccount            string  `json:"contaPagamento,omitempty"`
	AcceptanceIndicator       string  `json:"indicadorMovimentoAceito,omitempty"`
	Errors                    []int64 `json:"erros,omitempty"`
}

// PixPaymentItem represents a Pix payment item returned by a payment lookup.
type PixPaymentItem struct {
	COMPENumber                int64  `json:"numeroCOMPE,omitempty"`
	ISPBNumber                 int64  `json:"numeroISPB,omitempty"`
	CreditAgency               int64  `json:"agenciaCredito,omitempty"`
	CreditAccount              int64  `json:"contaCorrenteCredito,omitempty"`
	CreditAccountCheckDigit    string `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditPaymentAccountNumber string `json:"numeroContaPagamentoCredito,omitempty"`
	BeneficiaryType            int64  `json:"tipoBeneficiario,omitempty"`
	BeneficiaryDocument        int64  `json:"cpfCnpjBeneficiario,omitempty"`
	BeneficiaryName            string `json:"nomeBeneficiario,omitempty"`
	CreditDocumentNumber       int64  `json:"documentoCredito,omitempty"`
	InstantPaymentDescription  string `json:"descricaoPagamentoInstantaneo,omitempty"`
	AccountType                int64  `json:"tipoConta,omitempty"`
	IdentificationMode         string `json:"formaIdentificacao,omitempty"`
	PhoneAreaCode              int64  `json:"dddTelefone,omitempty"`
	PhoneNumber                int64  `json:"telefone,omitempty"`
	Email                      string `json:"email,omitempty"`
	RandomIdentifier           string `json:"identificacaoAleatoria,omitempty"`
	PixText                    string `json:"textoPix,omitempty"`
}

// PixReturnItem represents a Pix return entry.
type PixReturnItem struct {
	ReasonCode  int64   `json:"codigoMotivo,omitempty"`
	ReturnDate  int64   `json:"dataDevolucao,omitempty"`
	ReturnValue float64 `json:"valorDevolucao,omitempty"`
}

// CreateTransferBatchResponse is the response body for POST /lotes-transferencias.
type CreateTransferBatchResponse struct {
	RequestState       int64            `json:"estadoRequisicao,omitempty"`
	TransferCount      int64            `json:"quantidadeTransferencias,omitempty"`
	TransferValue      float64          `json:"valorTransferencias,omitempty"`
	ValidTransferCount int64            `json:"quantidadeTransferenciasValidas,omitempty"`
	ValidTransferValue float64          `json:"valorTransferenciasValidas,omitempty"`
	Transfers          []TransferResult `json:"transferencias,omitempty"`
}

// GetTransferPaymentResponse is the response body for GET /transferencias/{id}.
type GetTransferPaymentResponse struct {
	ID                    int64                 `json:"id"`
	PaymentState          string                `json:"estadoPagamento,omitempty"`
	CreditType            int64                 `json:"tipoCredito,omitempty"`
	DebitAgency           int64                 `json:"agenciaDebito,omitempty"`
	DebitAccount          int64                 `json:"contaCorrenteDebito,omitempty"`
	DebitAccountDigit     string                `json:"digitoVerificadorContaCorrente,omitempty"`
	CreditCardStart       int64                 `json:"inicioCartaoCredito,omitempty"`
	CreditCardEnd         int64                 `json:"fimCartaoCredito,omitempty"`
	PaymentDate           string                `json:"dataPagamento,omitempty"`
	PaymentValue          float64               `json:"valorPagamento,omitempty"`
	DebitDocumentNumber   int64                 `json:"documentoDebito,omitempty"`
	PaymentAuthentication string                `json:"codigoAutenticacaoPagamento,omitempty"`
	JudicialDepositNumber string                `json:"numeroDepositoJudicial,omitempty"`
	DOCPurposeCode        string                `json:"codigoFinalidadeDOC,omitempty"`
	TEDPurposeCode        string                `json:"codigoFinalidadeTED,omitempty"`
	PaymentItems          []TransferPaymentItem `json:"listaPagamentos,omitempty"`
	ReturnItems           []TransferReturnItem  `json:"listaDevolucao,omitempty"`
}

// GetBatchRequestResponse is the response body for GET /{id}/solicitacao.
type GetBatchRequestResponse struct {
	RequestState      int64                         `json:"estadoRequisicao,omitempty"`
	PaymentCount      int64                         `json:"quantidadePagamentos,omitempty"`
	PaymentValue      float64                       `json:"valorPagamentos,omitempty"`
	ValidPaymentCount int64                         `json:"quantidadePagamentosValidos,omitempty"`
	ValidPaymentValue float64                       `json:"valorPagamentosValidos,omitempty"` // fixed: was int64
	Payments          []TransferBatchRequestPayment `json:"pagamentos,omitempty"`
}

// GetBatchResponse is the response body for GET /{id}.
type GetBatchResponse struct {
	Index        int64                  `json:"indice,omitempty"`
	RequestState int64                  `json:"estadoRequisicao,omitempty"`
	PaymentType  int64                  `json:"tipoPagamento,omitempty"`
	RequestDate  int64                  `json:"dataRequisicao,omitempty"`
	PaymentCount int64                  `json:"quantidadePagamentos,omitempty"`
	PaymentValue float64                `json:"valorPagamentos,omitempty"`
	Payments     []TransferBatchPayment `json:"pagamentos,omitempty"`
}

// ListBeneficiaryTransfersParams holds query parameters for GET /beneficiarios/{id}/transferencias.
type ListBeneficiaryTransfersParams struct {
	DebitAgency             *int64
	DebitAccount            *int64
	DebitAccountCheckDigit  *string
	PaymentType             *int64
	COMPENumber             *int64
	ISPBNumber              *int64
	CreditAgency            *int64
	CreditAccount           *int64
	CreditAccountCheckDigit *string
	CreditPaymentAccount    *string
	StartDate               int64
	EndDate                 int64
	Index                   int64
	BeneficiaryType         int64
}

// ListBeneficiaryTransfersResponse is the response body for GET /beneficiarios/{id}/transferencias.
type ListBeneficiaryTransfersResponse struct {
	Index              int64                 `json:"indice"`
	TotalTransferCount int64                 `json:"quantidadeTotalTransferencias,omitempty"`
	Transfers          []BeneficiaryTransfer `json:"transferencias,omitempty"`
}

// CreatePixTransferBatchRequest is the request body for POST /lotes-transferencias-pix.
type CreatePixTransferBatchRequest struct {
	RequestNumber          int64         `json:"numeroRequisicao"`
	ContractNumber         *int64        `json:"numeroContrato,omitempty"`
	DebitAgency            *int64        `json:"agenciaDebito,omitempty"`
	DebitAccount           *int64        `json:"contaCorrenteDebito,omitempty"`
	DebitAccountCheckDigit *string       `json:"digitoVerificadorContaCorrente,omitempty"`
	PaymentType            int64         `json:"tipoPagamento"`
	Transfers              []PixTransfer `json:"listaTransferencias"`
}

// PixTransfer represents a Pix transfer entry.
type PixTransfer struct {
	Date                      int64   `json:"data"`
	Value                     float64 `json:"valor"`
	DebitDocumentNumber       *int64  `json:"documentoDebito,omitempty"`
	CreditDocumentNumber      *int64  `json:"documentoCredito,omitempty"`
	PaymentDescription        *string `json:"descricaoPagamento,omitempty"`
	InstantPaymentDescription *string `json:"descricaoPagamentoInstantaneo,omitempty"`
	IdentificationMode        int64   `json:"formaIdentificacao"`
	PhoneAreaCode             *int64  `json:"dddTelefone,omitempty"`
	PhoneNumber               *int64  `json:"telefone,omitempty"`
	Email                     *string `json:"email,omitempty"`
	BeneficiaryCPF            *int64  `json:"cpf,omitempty"`
	BeneficiaryCNPJ           *int64  `json:"cnpj,omitempty"`
	RandomIdentifier          *string `json:"identificacaoAleatoria,omitempty"`
	COMPENumber               *int64  `json:"numeroCOMPE,omitempty"`
	ISPBNumber                *int64  `json:"numeroISPB,omitempty"`
	AccountType               *int64  `json:"tipoConta,omitempty"`
	Agency                    *int64  `json:"agencia,omitempty"`
	Account                   *int64  `json:"conta,omitempty"`
	AccountCheckDigit         *string `json:"digitoVerificadorConta,omitempty"`
	PaymentAccount            *string `json:"contaPagamento,omitempty"`
}

// CreatePixTransferBatchResponse is the response body for POST /lotes-transferencias-pix.
type CreatePixTransferBatchResponse struct {
	RequestNumber      int64                  `json:"numeroRequisicao,omitempty"`
	RequestState       int64                  `json:"estadoRequisicao,omitempty"`
	TransferCount      int64                  `json:"quantidadeTransferencias,omitempty"`
	TransferValue      float64                `json:"valorTransferencias,omitempty"`
	ValidTransferCount int64                  `json:"quantidadeTransferenciasValidas,omitempty"`
	ValidTransferValue float64                `json:"valorTransferenciasValidas,omitempty"`
	Transfers          []PixTransferBatchItem `json:"listaTransferencias,omitempty"`
}

// GetPixTransferBatchRequestResponse is the response body for GET /lotes-transferencias-pix/{id}/solicitacao.
type GetPixTransferBatchRequestResponse struct {
	RequestNumber      int64                  `json:"numeroRequisicao,omitempty"`
	RequestState       int64                  `json:"estadoRequisicao,omitempty"`
	TransferCount      int64                  `json:"quantidadeTransferencias,omitempty"`
	TransferValue      float64                `json:"valorTransferencias,omitempty"`
	ValidTransferCount int64                  `json:"quantidadeTransferenciasValidas,omitempty"`
	ValidTransferValue float64                `json:"valorTransferenciasValidas,omitempty"`
	Transfers          []PixTransferBatchItem `json:"listaTransferencias,omitempty"`
}

// GetPixPaymentParams holds query parameters for GET /pix/{id}.
type GetPixPaymentParams struct {
	Agency      *int64
	Account     *int64
	CheckDigit  *string
	ExtraFields *string
}

// GetPixPaymentResponse is the response body for GET /pix/{id}.
type GetPixPaymentResponse struct {
	ID                   int64            `json:"id"`
	PaymentState         string           `json:"estadoPagamento,omitempty"`
	DebitAgency          int64            `json:"agencia,omitempty"`
	DebitAccount         int64            `json:"contaDebito,omitempty"`
	DebitAccountDigit    string           `json:"digitoContaDebito,omitempty"`
	CardStart            int64            `json:"numeroCartaoInicio,omitempty"`
	CardEnd              int64            `json:"numeroCartaoFim,omitempty"`
	PaymentRequestNumber int64            `json:"requisicaoPagamento,omitempty"`
	PaymentFile          string           `json:"arquivoPagamento,omitempty"`
	PaymentDate          int64            `json:"dataPagamento,omitempty"`
	PaymentValue         float64          `json:"valorPagamento,omitempty"`
	DebitDocumentNumber  int64            `json:"numeroDocumentoDebito,omitempty"`
	AuthenticationCode   string           `json:"autenticacaoPagamento,omitempty"`
	PaymentDescription   string           `json:"descricaoPagamento,omitempty"`
	PixOccurrenceCount   int64            `json:"quantidadeOcorrenciaPix,omitempty"`
	PixItems             []PixPaymentItem `json:"listaPix,omitempty"`
	ReturnItems          []PixReturnItem  `json:"listaDevolucao,omitempty"`
}

// ListTransferBatches lists transfer batches.
func (c *Client) ListTransferBatches(
	ctx context.Context,
	params *ListTransferBatchesParams,
) (*ListTransferBatchesResponse, error) {
	query := url.Values{}
	if params != nil {
		setInt64(query, "numeroContratoPagamento", params.PaymentContractNumber)
		setInt64(query, "agenciaDebito", params.DebitAgency)
		setInt64(query, "contaCorrenteDebito", params.DebitAccount)
		setString(query, "digitoVerificadorContaCorrente", params.DebitAccountCheckDigit)
		setInt64(query, "dataInicio", params.StartDate)
		setInt64(query, "dataFim", params.EndDate)
		setInt64(query, "tipoPagamento", params.PaymentType)
		setInt64(query, "estadoRequisicao", params.RequestState)
		setInt64(query, "indice", params.Index)
	}
	return get[*ListTransferBatchesResponse](c, ctx, buildPath(endpointTransferBatches, query))
}

// CreateTransferBatch creates a transfer batch.
func (c *Client) CreateTransferBatch(
	ctx context.Context,
	req *CreateTransferBatchRequest,
) (*CreateTransferBatchResponse, error) {
	return post[*CreateTransferBatchResponse](c, ctx, endpointTransferBatches, req)
}

// GetTransferPayment returns a single transfer payment.
func (c *Client) GetTransferPayment(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetTransferPaymentResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetTransferPaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointTransferPayment, id), query),
	)
}

// GetBatchRequest returns the request-stage representation of a batch.
func (c *Client) GetBatchRequest(ctx context.Context, id string) (*GetBatchRequestResponse, error) {
	return get[*GetBatchRequestResponse](c, ctx, fmt.Sprintf(endpointBatchRequest, id))
}

// GetBatch returns a payment batch by identifier.
func (c *Client) GetBatch(ctx context.Context, id string) (*GetBatchResponse, error) {
	return get[*GetBatchResponse](c, ctx, fmt.Sprintf(endpointBatch, id))
}

// ListBeneficiaryTransfers lists transfers by beneficiary.
func (c *Client) ListBeneficiaryTransfers(
	ctx context.Context,
	id string,
	params *ListBeneficiaryTransfersParams,
) (*ListBeneficiaryTransfersResponse, error) {
	query := url.Values{}
	if params != nil {
		setInt64(query, "agenciaDebito", params.DebitAgency)
		setInt64(query, "contaCorrenteDebito", params.DebitAccount)
		setString(query, "digitoVerificadorContaCorrente", params.DebitAccountCheckDigit)
		setInt64(query, "tipoPagamento", params.PaymentType)
		setInt64(query, "numeroCOMPE", params.COMPENumber)
		setInt64(query, "numeroISPB", params.ISPBNumber)
		setInt64(query, "agenciaCredito", params.CreditAgency)
		setInt64(query, "contaCorrenteCredito", params.CreditAccount)
		setString(query, "digitoVerificadorContaCorrenteCredito", params.CreditAccountCheckDigit)
		setString(query, "contaPagamentoCredito", params.CreditPaymentAccount)
		query.Set("dataInicio", strconv.FormatInt(params.StartDate, 10))
		query.Set("dataFim", strconv.FormatInt(params.EndDate, 10))
		query.Set("indice", strconv.FormatInt(params.Index, 10))
		query.Set("tipoBeneficiario", strconv.FormatInt(params.BeneficiaryType, 10))
	}
	return get[*ListBeneficiaryTransfersResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointBeneficiaryTransfers, id), query),
	)
}

// CreatePixTransferBatch creates a Pix transfer batch.
func (c *Client) CreatePixTransferBatch(
	ctx context.Context,
	req *CreatePixTransferBatchRequest,
) (*CreatePixTransferBatchResponse, error) {
	return post[*CreatePixTransferBatchResponse](c, ctx, endpointPixTransferBatches, req)
}

// GetPixTransferBatchRequest returns the request-stage representation of a Pix batch.
func (c *Client) GetPixTransferBatchRequest(
	ctx context.Context,
	id string,
	params *AccountLookupParams,
) (*GetPixTransferBatchRequestResponse, error) {
	query := url.Values{}
	setAccountLookupQuery(query, params, "agencia", "contaCorrente", "digitoVerificador")
	return get[*GetPixTransferBatchRequestResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointPixTransferBatchRequest, id), query),
	)
}

// GetPixPayment returns a single Pix payment.
func (c *Client) GetPixPayment(
	ctx context.Context,
	id string,
	params *GetPixPaymentParams,
) (*GetPixPaymentResponse, error) {
	query := url.Values{}
	if params != nil {
		setInt64(query, "agencia", params.Agency)
		setInt64(query, "contaCorrente", params.Account)
		setString(query, "digitoVerificador", params.CheckDigit)
		setString(query, "camposExtras", params.ExtraFields)
	}
	return get[*GetPixPaymentResponse](
		c, ctx,
		buildPath(fmt.Sprintf(endpointPixPayment, id), query),
	)
}
