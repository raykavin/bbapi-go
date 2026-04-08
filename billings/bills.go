package billings

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/raykavin/bbapi-go"
)

const (
	endpointBillets              = "/boletos"
	endpointBillet               = "/boletos/%s"
	endpointBilletCancel         = "/boletos/%s/baixar"
	endpointBilletGeneratePix    = "/boletos/%s/gerar-pix"
	endpointBilletCancelPix      = "/boletos/%s/cancelar-pix"
	endpointBilletPix            = "/boletos/%s/pix"
	endpointOperationalWriteOffs = "/boletos-baixa-operacional"
)

// Discount defines a discount block for billing registration.
type Discount struct {
	Type       *int64   `json:"tipo,omitempty"`
	Expiration *string  `json:"dataExpiracao,omitempty"`
	Percentage *float64 `json:"porcentagem,omitempty"`
	Value      *float64 `json:"valor,omitempty"`
}

// Interest defines interest settings for a billing.
type Interest struct {
	Type       *int64   `json:"tipo,omitempty"`
	Percentage *float64 `json:"porcentagem,omitempty"`
	Value      *float64 `json:"valor,omitempty"`
}

// Fine defines fine settings for a billing.
type Fine struct {
	Type       *int64   `json:"tipo,omitempty"`
	Date       *string  `json:"data,omitempty"`
	Percentage *float64 `json:"porcentagem,omitempty"`
	Value      *float64 `json:"valor,omitempty"`
}

// Payer defines the payer information used when registering a billing.
type Payer struct {
	RegistrationType   int64   `json:"tipoInscricao"`
	RegistrationNumber int64   `json:"numeroInscricao"`
	Name               *string `json:"nome,omitempty"`
	Address            *string `json:"endereco,omitempty"`
	ZipCode            *int64  `json:"cep,omitempty"`
	City               *string `json:"cidade,omitempty"`
	Neighborhood       *string `json:"bairro,omitempty"`
	State              *string `json:"uf,omitempty"`
	Phone              *string `json:"telefone,omitempty"`
	Email              *string `json:"email,omitempty"`
}

// FinalBeneficiary defines the final beneficiary block for a billing.
type FinalBeneficiary struct {
	RegistrationType   *int64  `json:"tipoInscricao,omitempty"`
	RegistrationNumber *int64  `json:"numeroInscricao,omitempty"`
	Name               *string `json:"nome,omitempty"`
}

// CreateBillingRequest is the request body for POST /boletos.
type CreateBillingRequest struct {
	AgreementNumber                int64             `json:"numeroConvenio"`
	WalletNumber                   *int64            `json:"numeroCarteira,omitempty"`
	WalletVariationNumber          *int64            `json:"numeroVariacaoCarteira,omitempty"`
	ModalityCode                   *int64            `json:"codigoModalidade,omitempty"`
	IssueDate                      *string           `json:"dataEmissao,omitempty"`
	DueDate                        string            `json:"dataVencimento"`
	OriginalValue                  float64           `json:"valorOriginal"`
	RebateValue                    *float64          `json:"valorAbatimento,omitempty"`
	ProtestDays                    *float64          `json:"quantidadeDiasProtesto,omitempty"`
	NegativeReportingDays          *int64            `json:"quantidadeDiasNegativacao,omitempty"`
	NegativeReportingAgency        *int64            `json:"orgaoNegativador,omitempty"`
	AcceptExpiredTitleIndicator    *string           `json:"indicadorAceiteTituloVencido,omitempty"`
	ReceiptLimitDays               *int64            `json:"numeroDiasLimiteRecebimento,omitempty"`
	AcceptanceCode                 *string           `json:"codigoAceite,omitempty"`
	TitleTypeCode                  *int64            `json:"codigoTipoTitulo,omitempty"`
	TitleTypeDescription           *string           `json:"descricaoTipoTitulo,omitempty"`
	PartialPaymentAllowedIndicator *string           `json:"indicadorPermissaoRecebimentoParcial,omitempty"`
	BeneficiaryTitleNumber         *string           `json:"numeroTituloBeneficiario,omitempty"`
	BeneficiaryUsageField          *string           `json:"campoUtilizacaoBeneficiario,omitempty"`
	ClientTitleNumber              *string           `json:"numeroTituloCliente,omitempty"`
	BillMessage                    *string           `json:"mensagemBloquetoOcorrencia,omitempty"`
	Discount                       *Discount         `json:"desconto,omitempty"`
	SecondDiscount                 *Discount         `json:"segundoDesconto,omitempty"`
	ThirdDiscount                  *Discount         `json:"terceiroDesconto,omitempty"`
	Interest                       *Interest         `json:"jurosMora,omitempty"`
	Fine                           *Fine             `json:"multa,omitempty"`
	Payer                          *Payer            `json:"pagador,omitempty"`
	FinalBeneficiary               *FinalBeneficiary `json:"beneficiarioFinal,omitempty"`
	PixIndicator                   *string           `json:"indicadorPix,omitempty"`
	PixLocationID                  *int64            `json:"idLocationPix,omitempty"`
}

// BillingBeneficiaryResponse describes the beneficiary block returned after registration.
type BillingBeneficiaryResponse struct {
	Agency          *int64  `json:"agencia,omitempty"`
	CheckingAccount *int64  `json:"contaCorrente,omitempty"`
	AddressType     *int64  `json:"tipoEndereco,omitempty"`
	Street          *string `json:"logradouro,omitempty"`
	Neighborhood    *string `json:"bairro,omitempty"`
	City            *string `json:"cidade,omitempty"`
	CityCode        *int64  `json:"codigoCidade,omitempty"`
	State           *string `json:"uf,omitempty"`
	ZipCode         *int64  `json:"cep,omitempty"`
	ProofIndicator  *string `json:"indicadorComprovacao,omitempty"`
}

// PixQRCode contains the generated Pix QR Code data.
type PixQRCode struct {
	URL  *string `json:"url,omitempty"`
	TxID *string `json:"txId,omitempty"`
	EMV  *string `json:"emv,omitempty"`
	Type *int64  `json:"tipo,omitempty"`
}

// PixPayload represents Pix data associated with a billing.
type PixPayload struct {
	ReceivedValue      *float64 `json:"valorRecebido,omitempty"`
	Timestamp          *string  `json:"timestamp,omitempty"`
	Key                *string  `json:"chave,omitempty"`
	ReturnText         *string  `json:"textoRetorno,omitempty"`
	PayerInstitutionID *int64   `json:"idInstituicaoPagador,omitempty"`
	PayerAgency        *int64   `json:"agenciaPagador,omitempty"`
	PayerAccount       *int64   `json:"contaPagador,omitempty"`
	PayerPersonType    *int64   `json:"tipoPessoaPagador,omitempty"`
	PayerID            *int64   `json:"idPagador,omitempty"`
}

// CreateBillingResponse is the response body for POST /boletos.
type CreateBillingResponse struct {
	Number                string                      `json:"numero,omitempty"`
	WalletNumber          *int64                      `json:"numeroCarteira,omitempty"`
	WalletVariationNumber *int64                      `json:"numeroVariacaoCarteira,omitempty"`
	ClientCode            *int64                      `json:"codigoCliente,omitempty"`
	DigitableLine         *string                     `json:"linhaDigitavel,omitempty"`
	NumericBarcode        *string                     `json:"codigoBarraNumerico,omitempty"`
	BillingContractNumber *int64                      `json:"numeroContratoCobranca,omitempty"`
	Beneficiary           *BillingBeneficiaryResponse `json:"beneficiario,omitempty"`
	QRCode                *PixQRCode                  `json:"qrCode,omitempty"`
	BillImageURL          *string                     `json:"urlImagemBoleto,omitempty"`
	Observation           *string                     `json:"observacao,omitempty"`
}

// ListBillingsParams holds query parameters for GET /boletos.
type ListBillingsParams struct {
	StatusIndicator          string
	CollateralAccount        *int64
	BeneficiaryAgency        int64
	BeneficiaryAccount       int64
	AgreementWallet          *int64
	AgreementWalletVariation *int64
	ChargingMode             *int64
	PayerCNPJ                *int64
	PayerCNPJDigit           *int64
	PayerCPF                 *int64
	PayerCPFDigit            *int64
	StartDueDate             *string
	EndDueDate               *string
	StartRegisterDate        *string
	EndRegisterDate          *string
	StartMovementDate        *string
	EndMovementDate          *string
	BillingStateCode         *int64
	OverdueBilletIndicator   *string
	Index                    *int64
}

// BillingListItem describes an entry returned by GET /boletos.
type BillingListItem struct {
	BBNumber           string   `json:"numeroBoletoBB,omitempty"`
	BillingState       *string  `json:"estadoTituloCobranca,omitempty"`
	RegisterDate       *string  `json:"dataRegistro,omitempty"`
	DueDate            *string  `json:"dataVencimento,omitempty"`
	MovementDate       *string  `json:"dataMovimento,omitempty"`
	OriginalValue      *float64 `json:"valorOriginal,omitempty"`
	CurrentValue       *float64 `json:"valorAtual,omitempty"`
	PaidValue          *float64 `json:"valorPago,omitempty"`
	Contract           *int64   `json:"contrato,omitempty"`
	AgreementWallet    *int64   `json:"carteiraConvenio,omitempty"`
	AgreementVariation *int64   `json:"variacaoCarteiraConvenio,omitempty"`
	BillingStateCode   *int64   `json:"codigoEstadoTituloCobranca,omitempty"`
	CreditDate         *string  `json:"dataCredito,omitempty"`
}

// ListBillingsResponse is the response body for GET /boletos.
type ListBillingsResponse struct {
	ContinuationIndicator *string           `json:"indicadorContinuidade,omitempty"`
	RecordCount           *int64            `json:"quantidadeRegistros,omitempty"`
	NextIndex             *int64            `json:"proximoIndice,omitempty"`
	Billings              []BillingListItem `json:"boletos,omitempty"`
}

// GetBillingParams holds query parameters for GET /boletos/{id}.
type GetBillingParams struct {
	AgreementNumber int64
}

// BillingDetailResponse is the response body for GET /boletos/{id}.
// The API exposes a large contract; this struct captures the most relevant typed fields.
type BillingDetailResponse struct {
	DigitableLine           *string    `json:"codigoLinhaDigitavel,omitempty"`
	PayerEmail              *string    `json:"textoEmailPagador,omitempty"`
	BillMessage             *string    `json:"textoMensagemBloquetoTitulo,omitempty"`
	FineTypeCode            *int64     `json:"codigoTipoMulta,omitempty"`
	PaymentChannelCode      *int64     `json:"codigoCanalPagamento,omitempty"`
	BillingContractNumber   *int64     `json:"numeroContratoCobranca,omitempty"`
	PayerRegistrationType   *int64     `json:"codigoTipoInscricaoSacado,omitempty"`
	PayerRegistrationNumber *int64     `json:"numeroInscricaoSacadoCobranca,omitempty"`
	BillingStateCode        *int64     `json:"codigoEstadoTituloCobranca,omitempty"`
	BillingTypeCode         *int64     `json:"codigoTipoTituloCobranca,omitempty"`
	BillingModalityCode     *int64     `json:"codigoModalidadeTitulo,omitempty"`
	AcceptanceCode          *string    `json:"codigoAceiteTituloCobranca,omitempty"`
	ChargingAgencyPrefix    *int64     `json:"codigoPrefixoDependenciaCobrador,omitempty"`
	EconomicIndicatorCode   *int64     `json:"codigoIndicadorEconomico,omitempty"`
	BeneficiaryTitleNumber  *string    `json:"numeroTituloCedenteCobranca,omitempty"`
	InterestTypeCode        *int64     `json:"codigoTipoJuroMora,omitempty"`
	IssueDate               *string    `json:"dataEmissaoTituloCobranca,omitempty"`
	RegisterDate            *string    `json:"dataRegistroTituloCobranca,omitempty"`
	DueDate                 *string    `json:"dataVencimentoTituloCobranca,omitempty"`
	OriginalValue           *float64   `json:"valorOriginalTituloCobranca,omitempty"`
	CurrentValue            *float64   `json:"valorAtualTituloCobranca,omitempty"`
	PaidValue               *float64   `json:"valorPagoSacado,omitempty"`
	DiscountValue           *float64   `json:"valorDescontoTituloCobranca,omitempty"`
	RebateValue             *float64   `json:"valorAbatimentoTituloCobranca,omitempty"`
	PayerName               *string    `json:"nomeSacadoCobranca,omitempty"`
	PayerAddress            *string    `json:"textoEnderecoSacadoCobranca,omitempty"`
	PayerZIPCode            *int64     `json:"numeroCepSacadoCobranca,omitempty"`
	PayerCity               *string    `json:"nomeMunicipioSacadoCobranca,omitempty"`
	PayerDistrict           *string    `json:"nomeBairroSacadoCobranca,omitempty"`
	PayerState              *string    `json:"siglaUnidadeFederacaoSacadoCobranca,omitempty"`
	QRCode                  *PixQRCode `json:"qrCode,omitempty"`
}

// ChangeNominalValue defines the new nominal billing value.
type ChangeNominalValue struct {
	NewNominalValue *float64 `json:"novoValorNominal,omitempty"`
}

// Rebate defines the rebate block used by billing change operations.
type Rebate struct {
	Value *float64 `json:"valorAbatimento,omitempty"`
}

// RebateChange defines a rebate change block for billing updates.
type RebateChange struct {
	NewValue *float64 `json:"novoValorAbatimento,omitempty"`
}

// DueDateChange defines a new due date for billing updates.
type DueDateChange struct {
	NewDueDate *string `json:"novaDataVencimento,omitempty"`
}

// DiscountDateChange defines new discount limit dates.
type DiscountDateChange struct {
	NewFirstDiscountDate  *string `json:"novaDataLimitePrimeiroDesconto,omitempty"`
	NewSecondDiscountDate *string `json:"novaDataLimiteSegundoDesconto,omitempty"`
	NewThirdDiscountDate  *string `json:"novaDataLimiteTerceiroDesconto,omitempty"`
}

// DiscountChange defines discount update rules for billing changes.
type DiscountChange struct {
	FirstDiscountType        *int64   `json:"tipoPrimeiroDesconto,omitempty"`
	NewFirstDiscountValue    *float64 `json:"novoValorPrimeiroDesconto,omitempty"`
	NewFirstDiscountPercent  *float64 `json:"novoPercentualPrimeiroDesconto,omitempty"`
	NewFirstDiscountDate     *string  `json:"novaDataLimitePrimeiroDesconto,omitempty"`
	SecondDiscountType       *int64   `json:"tipoSegundoDesconto,omitempty"`
	NewSecondDiscountValue   *float64 `json:"novoValorSegundoDesconto,omitempty"`
	NewSecondDiscountPercent *float64 `json:"novoPercentualSegundoDesconto,omitempty"`
	NewSecondDiscountDate    *string  `json:"novaDataLimiteSegundoDesconto,omitempty"`
	ThirdDiscountType        *int64   `json:"tipoTerceiroDesconto,omitempty"`
	NewThirdDiscountValue    *float64 `json:"novoValorTerceiroDesconto,omitempty"`
	NewThirdDiscountPercent  *float64 `json:"novoPercentualTerceiroDesconto,omitempty"`
	NewThirdDiscountDate     *string  `json:"novaDataLimiteTerceiroDesconto,omitempty"`
}

// PayerAddressChange defines payer address changes.
type PayerAddressChange struct {
	Address      *string `json:"enderecoPagador,omitempty"`
	Neighborhood *string `json:"bairroPagador,omitempty"`
	City         *string `json:"cidadePagador,omitempty"`
	State        *string `json:"UFPagador,omitempty"`
	ZIPCode      *int64  `json:"CEPPagador,omitempty"`
}

// AcceptanceWindowChange defines the acceptance window after expiration.
type AcceptanceWindowChange struct {
	AcceptanceDays *int64 `json:"quantidadeDiasAceite,omitempty"`
}

// BeneficiaryNumberChange defines the beneficiary title number change.
type BeneficiaryNumberChange struct {
	YourNumber *string `json:"codigoSeuNumero,omitempty"`
}

// UpdateBillingRequest is the request body for PATCH /boletos/{id}.
type UpdateBillingRequest struct {
	AgreementNumber                  int64                    `json:"numeroConvenio"`
	ChangeDueDateIndicator           *string                  `json:"indicadorNovaDataVencimento,omitempty"`
	DueDateChange                    *DueDateChange           `json:"alteracaoData,omitempty"`
	ChangeProtestIndicator           *string                  `json:"indicadorProtestar,omitempty"`
	Protest                          *Protest                 `json:"protesto,omitempty"`
	ChangeNegativeReportingIndicator *string                  `json:"indicadorNegativar,omitempty"`
	NegativeReporting                *NegativeReporting       `json:"negativacao,omitempty"`
	ChangeRebateIndicator            *string                  `json:"indicadorAlterarAbatimento,omitempty"`
	RebateChange                     *RebateChange            `json:"alteracaoAbatimento,omitempty"`
	ChangeDiscountIndicator          *string                  `json:"indicadorAlterarDesconto,omitempty"`
	DiscountChange                   *DiscountChange          `json:"alteracaoDesconto,omitempty"`
	ChangeDiscountDateIndicator      *string                  `json:"indicadorAlterarDataDesconto,omitempty"`
	DiscountDateChange               *DiscountDateChange      `json:"alteracaoDataDesconto,omitempty"`
	ChangeNominalValueIndicator      *string                  `json:"indicadorAlterarValorNominal,omitempty"`
	NominalValueChange               *ChangeNominalValue      `json:"AlterarValorNominal,omitempty"`
	ChangeBeneficiaryNumberIndicator *string                  `json:"indicadorAlterarSeuNumero,omitempty"`
	BeneficiaryNumberChange          *BeneficiaryNumberChange `json:"alteracaoSeuNumero,omitempty"`
	ChangePayerAddressIndicator      *string                  `json:"indicadorAlterarEnderecoPagador,omitempty"`
	PayerAddressChange               *PayerAddressChange      `json:"alteracaoEndereco,omitempty"`
	ChangeExpiredAcceptanceIndicator *string                  `json:"indicadorAlterarPrazoBoletoVencido,omitempty"`
	AcceptanceWindowChange           *AcceptanceWindowChange  `json:"alteracaoPrazo,omitempty"`
}

// Protest defines the protest block used by billing updates.
type Protest struct {
	ProtestDays *float64 `json:"quantidadeDiasProtesto,omitempty"`
}

// NegativeReporting defines the negative reporting block used by billing updates.
type NegativeReporting struct {
	NegativeReportingDays *int64 `json:"quantidadeDiasNegativacao,omitempty"`
	NegativeReportingType *int64 `json:"tipoNegativacao,omitempty"`
	NegativeReportingOrg  *int64 `json:"orgaoNegativador,omitempty"`
}

// UpdateBillingResponse is the response body for PATCH /boletos/{id}.
type UpdateBillingResponse struct {
	BillingContractNumber *int64  `json:"numeroContratoCobranca,omitempty"`
	UpdateDate            *string `json:"dataAtualizacao,omitempty"`
	UpdateTime            *string `json:"horarioAtualizacao,omitempty"`
}

// CancelBillingRequest is the request body for POST /boletos/{id}/baixar.
type CancelBillingRequest struct {
	AgreementNumber int64 `json:"numeroConvenio"`
}

// CancelBillingResponse is the response body for POST /boletos/{id}/baixar.
type CancelBillingResponse struct {
	BillingContractNumber *string `json:"numeroContratoCobranca,omitempty"`
	CancelDate            *string `json:"dataBaixa,omitempty"`
	CancelTime            *string `json:"horarioBaixa,omitempty"`
}

// BillingPixRequest is the request body for Pix generation/cancellation.
type BillingPixRequest struct {
	AgreementNumber int64 `json:"numeroConvenio"`
}

// BillingPixResponse is the response body for Pix generation/cancellation.
type BillingPixResponse struct {
	Pix *struct {
		Key *string `json:"chave,omitempty"`
	} `json:"pix,omitempty"`
	QRCode *PixQRCode `json:"qrCode,omitempty"`
}

// GetBillingPixParams holds query parameters for GET /boletos/{id}/pix.
type GetBillingPixParams struct {
	AgreementNumber int64
}

// GetBillingPixResponse is the response body for GET /boletos/{id}/pix.
type GetBillingPixResponse struct {
	ID                 string      `json:"id,omitempty"`
	RegisterDate       *string     `json:"dataRegistroTituloCobranca,omitempty"`
	BeneficiaryAgency  *int64      `json:"agenciaBeneficiario,omitempty"`
	BeneficiaryAccount *int64      `json:"contaBeneficiario,omitempty"`
	OriginalValue      *float64    `json:"valorOriginalTituloCobranca,omitempty"`
	ValidityDate       *string     `json:"validadeTituloCobranca,omitempty"`
	Pix                *PixPayload `json:"pix,omitempty"`
	QRCode             *PixQRCode  `json:"qrCode,omitempty"`
}

// OperationalWriteOffListParams holds query parameters for GET /boletos-baixa-operacional.
type OperationalWriteOffListParams struct {
	Agency                   int64
	Account                  int64
	Wallet                   int64
	Variation                int64
	OperationalWriteOffState *int64
	TitleModality            *int64
	StartDueDate             *string
	EndDueDate               *string
	StartRegisterDate        *string
	EndRegisterDate          *string
	StartScheduleDate        string
	EndScheduleDate          string
	StartScheduleTime        *string
	EndScheduleTime          *string
	NextTitleID              *string
}

// OperationalWriteOffSchedule holds scheduling data for an operational write-off.
type OperationalWriteOffSchedule struct {
	Moment               *string  `json:"momento,omitempty"`
	FinancialInstitution *int64   `json:"instituicaoFinanceira,omitempty"`
	Channel              *int64   `json:"canal,omitempty"`
	CIPValue             *float64 `json:"valorCIP,omitempty"`
}

// OperationalWriteOffTitle describes a title in the operational write-off response.
type OperationalWriteOffTitle struct {
	ID              *string                      `json:"id,omitempty"`
	StateCode       *int64                       `json:"estadoBaixaOperacional,omitempty"`
	ModalityCode    *int64                       `json:"modalidade,omitempty"`
	RegisterDate    *string                      `json:"dataRegistro,omitempty"`
	DueDate         *string                      `json:"dataVencimento,omitempty"`
	OriginalValue   *float64                     `json:"valorOriginal,omitempty"`
	PaymentSchedule *OperationalWriteOffSchedule `json:"agendamentoPagamento,omitempty"`
}

// OperationalWriteOffItem describes each record from the operational write-off list.
type OperationalWriteOffItem struct {
	Wallet    *int64                    `json:"carteira,omitempty"`
	Variation *int64                    `json:"variacao,omitempty"`
	Agreement *int64                    `json:"convenio,omitempty"`
	Title     *OperationalWriteOffTitle `json:"titulo,omitempty"`
}

// OperationalWriteOffListResponse is the response body for GET /boletos-baixa-operacional.
type OperationalWriteOffListResponse struct {
	HasMoreTitles *string                   `json:"possuiMaisTitulos,omitempty"`
	NextTitle     *string                   `json:"proximoTitulo,omitempty"`
	Items         []OperationalWriteOffItem `json:"lista,omitempty"`
}

// ListBillings lists billings using the filter parameters supported by GET /boletos.
func (c *Client) ListBillings(
	ctx context.Context,
	params *ListBillingsParams,
) (*ListBillingsResponse, error) {
	query := url.Values{}

	if params != nil {
		query.Set("indicadorSituacao", params.StatusIndicator)
		query.Set("agenciaBeneficiario", strconv.FormatInt(params.BeneficiaryAgency, 10))
		query.Set("contaBeneficiario", strconv.FormatInt(params.BeneficiaryAccount, 10))

		bbapi.SetInt64(query, "contaCaucao", params.CollateralAccount)
		bbapi.SetInt64(query, "carteiraConvenio", params.AgreementWallet)
		bbapi.SetInt64(query, "variacaoCarteiraConvenio", params.AgreementWalletVariation)
		bbapi.SetInt64(query, "modalidadeCobranca", params.ChargingMode)
		bbapi.SetInt64(query, "cnpjPagador", params.PayerCNPJ)
		bbapi.SetInt64(query, "digitoCNPJPagador", params.PayerCNPJDigit)
		bbapi.SetInt64(query, "cpfPagador", params.PayerCPF)
		bbapi.SetInt64(query, "digitoCPFPagador", params.PayerCPFDigit)
		bbapi.SetString(query, "dataInicioVencimento", params.StartDueDate)
		bbapi.SetString(query, "dataFimVencimento", params.EndDueDate)
		bbapi.SetString(query, "dataInicioRegistro", params.StartRegisterDate)
		bbapi.SetString(query, "dataFimRegistro", params.EndRegisterDate)
		bbapi.SetString(query, "dataInicioMovimento", params.StartMovementDate)
		bbapi.SetString(query, "dataFimMovimento", params.EndMovementDate)
		bbapi.SetInt64(query, "codigoEstadoTituloCobranca", params.BillingStateCode)
		bbapi.SetString(query, "boletoVencido", params.OverdueBilletIndicator)
		bbapi.SetInt64(query, "indice", params.Index)
	}

	return bbapi.Get[*ListBillingsResponse](ctx, c.Client, bbapi.BuildPath(endpointBillets, query))
}

// CreateBilling registers a billing.
func (c *Client) CreateBilling(
	ctx context.Context,
	req *CreateBillingRequest,
) (*CreateBillingResponse, error) {
	return bbapi.Post[*CreateBillingResponse](ctx, c.Client, endpointBillets, req)
}

// GetBilling returns a single billing.
func (c *Client) GetBilling(
	ctx context.Context,
	id string,
	params *GetBillingParams,
) (*BillingDetailResponse, error) {
	query := url.Values{}
	if params != nil {
		query.Set("numeroConvenio", strconv.FormatInt(params.AgreementNumber, 10))
	}

	return bbapi.Get[*BillingDetailResponse](
		ctx,
		c.Client,
		bbapi.BuildPath(fmt.Sprintf(endpointBillet, id), query),
	)
}

// UpdateBilling changes a billing.
func (c *Client) UpdateBilling(
	ctx context.Context,
	id string,
	req *UpdateBillingRequest,
) (*UpdateBillingResponse, error) {
	return bbapi.Patch[*UpdateBillingResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointBillet, id),
		req,
	)
}

// CancelBilling requests the cancellation/write-off of a billing.
func (c *Client) CancelBilling(
	ctx context.Context,
	id string,
	req *CancelBillingRequest,
) (*CancelBillingResponse, error) {
	return bbapi.Post[*CancelBillingResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointBilletCancel, id),
		req,
	)
}

// GenerateBillingPix generates a Pix QR code linked to the billing.
func (c *Client) GenerateBillingPix(
	ctx context.Context,
	id string,
	req *BillingPixRequest,
) (*BillingPixResponse, error) {
	return bbapi.Post[*BillingPixResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointBilletGeneratePix, id),
		req,
	)
}

// CancelBillingPix cancels the Pix linked to the billing.
func (c *Client) CancelBillingPix(
	ctx context.Context,
	id string,
	req *BillingPixRequest,
) (*BillingPixResponse, error) {
	return bbapi.Post[*BillingPixResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointBilletCancelPix, id),
		req,
	)
}

// GetBillingPix returns Pix information linked to a billing.
func (c *Client) GetBillingPix(
	ctx context.Context,
	id string,
	params *GetBillingPixParams,
) (*GetBillingPixResponse, error) {
	query := url.Values{}
	if params != nil {
		query.Set("numeroConvenio", strconv.FormatInt(params.AgreementNumber, 10))
	}

	return bbapi.Get[*GetBillingPixResponse](
		ctx,
		c.Client,
		bbapi.BuildPath(fmt.Sprintf(endpointBilletPix, id), query),
	)
}

// ListOperationalWriteOffs lists same-day operational write-offs for billings.
func (c *Client) ListOperationalWriteOffs(
	ctx context.Context,
	params *OperationalWriteOffListParams,
) (*OperationalWriteOffListResponse, error) {
	query := url.Values{}

	if params != nil {
		query.Set("agencia", strconv.FormatInt(params.Agency, 10))
		query.Set("conta", strconv.FormatInt(params.Account, 10))
		query.Set("carteira", strconv.FormatInt(params.Wallet, 10))
		query.Set("variacao", strconv.FormatInt(params.Variation, 10))
		query.Set("dataInicioAgendamentoTitulo", params.StartScheduleDate)
		query.Set("dataFimAgendamentoTitulo", params.EndScheduleDate)

		bbapi.SetInt64(query, "estadoBaixaOperacional", params.OperationalWriteOffState)
		bbapi.SetInt64(query, "modalidadeTitulo", params.TitleModality)
		bbapi.SetString(query, "dataInicioVencimentoTitulo", params.StartDueDate)
		bbapi.SetString(query, "dataFimVencimentoTitulo", params.EndDueDate)
		bbapi.SetString(query, "dataInicioRegistroTitulo", params.StartRegisterDate)
		bbapi.SetString(query, "dataFimRegistroTitulo", params.EndRegisterDate)
		bbapi.SetString(query, "horarioInicioAgendamentoTitulo", params.StartScheduleTime)
		bbapi.SetString(query, "horarioFimAgendamentoTitulo", params.EndScheduleTime)
		bbapi.SetString(query, "idProximoTitulo", params.NextTitleID)
	}

	return bbapi.Get[*OperationalWriteOffListResponse](
		ctx,
		c.Client,
		bbapi.BuildPath(endpointOperationalWriteOffs, query),
	)
}
