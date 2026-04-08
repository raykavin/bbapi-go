package billings

import (
	"context"
	"fmt"

	"github.com/raykavin/bbapi-go"
)

const (
	endpointAgreementReturnMovements = "/convenios/%s/listar-retorno-movimento"
	endpointActivateOperationalQuery = "/convenios/%s/ativar-consulta-baixa-operacional"
	endpointDisableOperationalQuery  = "/convenios/%s/desativar-consulta-baixa-operacional"
)

// AgreementReturnMovementListRequest is the request body for
// POST /convenios/{id}/listar-retorno-movimento.
type AgreementReturnMovementListRequest struct {
	StartReturnMovementDate *string `json:"dataMovimentoRetornoInicial,omitempty"`
	EndReturnMovementDate   *string `json:"dataMovimentoRetornoFinal,omitempty"`
	AgencyPrefix            *int64  `json:"codigoPrefixoAgencia,omitempty"`
	CheckingAccount         *int64  `json:"numeroContaCorrente,omitempty"`
	WalletNumber            *int64  `json:"numeroCarteiraCobranca,omitempty"`
	WalletVariationNumber   *int64  `json:"numeroVariacaoCarteiraCobranca,omitempty"`
	DesiredRecordNumber     *int64  `json:"numeroRegistroPretendido,omitempty"`
	DesiredRecordCount      *int64  `json:"quantidadeRegistroPretendido,omitempty"`
}

// AgreementReturnMovementItem describes a record returned from the agreement return movement API.
type AgreementReturnMovementItem struct {
	ReturnMovementDate           *string  `json:"dataMovimentoRetorno,omitempty"`
	AgreementNumber              *int64   `json:"numeroConvenio,omitempty"`
	BillingTitleNumber           *string  `json:"numeroTituloCobranca,omitempty"`
	ActionCommandCode            *int64   `json:"codigoComandoAcao,omitempty"`
	AgencyPrefix                 *int64   `json:"codigoPrefixoAgencia,omitempty"`
	CheckingAccount              *int64   `json:"numeroContaCorrente,omitempty"`
	WalletNumber                 *int64   `json:"numeroCarteiraCobranca,omitempty"`
	WalletVariationNumber        *int64   `json:"numeroVariacaoCarteiraCobranca,omitempty"`
	ChargingType                 *int64   `json:"tipoCobranca,omitempty"`
	ParticipantControlCode       *string  `json:"codigoControleParticipante,omitempty"`
	BillingSpeciesCode           *int64   `json:"codigoEspecieBoleto,omitempty"`
	DueDate                      *string  `json:"dataVencimentoBoleto,omitempty"`
	BillingValue                 *float64 `json:"valorBoleto,omitempty"`
	ReceivingBankCode            *int64   `json:"codigoBancoRecebedor,omitempty"`
	ReceivingAgencyPrefix        *int64   `json:"codigoPrefixoAgenciaRecebedora,omitempty"`
	PaymentCreditDate            *string  `json:"dataCreditoPagamentoBoleto,omitempty"`
	TariffValue                  *float64 `json:"valorTarifa,omitempty"`
	CalculatedOtherExpensesValue *float64 `json:"valorOutrasDespesasCalculadas,omitempty"`
	InterestDiscountValue        *float64 `json:"valorJurosDesconto,omitempty"`
	IOFDiscountValue             *float64 `json:"valorIofDesconto,omitempty"`
	RebateValue                  *float64 `json:"valorAbatimento,omitempty"`
	DiscountValue                *float64 `json:"valorDesconto,omitempty"`
	ReceivedValue                *float64 `json:"valorRecebido,omitempty"`
	LateInterestValue            *float64 `json:"valorJurosMora,omitempty"`
	OtherReceivedValues          *float64 `json:"valorOutrosValoresRecebidos,omitempty"`
	UnusedRebateValue            *float64 `json:"valorAbatimentoNaoUtilizado,omitempty"`
	PostingValue                 *float64 `json:"valorLancamento,omitempty"`
	PaymentMethodCode            *int64   `json:"codigoFormaPagamento,omitempty"`
	AdjustmentValueCode          *int64   `json:"codigoValorAjuste,omitempty"`
	AdjustmentValue              *float64 `json:"valorAjuste,omitempty"`
	PartialPaymentAuthorization  *int64   `json:"codigoAutorizacaoPagamentoParcial,omitempty"`
	PaymentChannelCode           *int64   `json:"codigoCanalPagamento,omitempty"`
	URL                          *string  `json:"URL,omitempty"`
	QRCodeIdentifierText         *string  `json:"textoIdentificadorQRCode,omitempty"`
	CalculationDays              *int64   `json:"quantidadeDiasCalculo,omitempty"`
	DiscountRateValue            *float64 `json:"valorTaxaDesconto,omitempty"`
	IOFRateValue                 *float64 `json:"valorTaxaIOF,omitempty"`
	ReceivingNature              *int64   `json:"naturezaRecebimento,omitempty"`
	CommandChargingTypeCode      *int64   `json:"codigoTipoCobrancaComando,omitempty"`
	LiquidationDate              *string  `json:"dataLiquidacaoBoleto,omitempty"`
}

// AgreementReturnMovementListResponse is the response body for
// POST /convenios/{id}/listar-retorno-movimento.
type AgreementReturnMovementListResponse struct {
	ContinuationIndicator *string                       `json:"indicadorContinuidade,omitempty"`
	LastRecordNumber      *int64                        `json:"numeroUltimoRegistro,omitempty"`
	Records               []AgreementReturnMovementItem `json:"listaRegistro,omitempty"`
}

// AgreementCustomizationStateResponse is the response body for the activation
// and deactivation PATCH endpoints.
type AgreementCustomizationStateResponse struct {
	CustomizationState *string `json:"estadoPersonalizacao,omitempty"`
	StateDateTime      *string `json:"dataHoraEstado,omitempty"`
}

// ListAgreementReturnMovements lists return movement data for an agreement.
func (c *Client) ListAgreementReturnMovements(
	ctx context.Context,
	id string,
	req *AgreementReturnMovementListRequest,
) (*AgreementReturnMovementListResponse, error) {
	return bbapi.Post[*AgreementReturnMovementListResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointAgreementReturnMovements, id),
		req,
	)
}

// ActivateOperationalWriteOffConsultation enables same-day operational write-off consultation.
func (c *Client) ActivateOperationalWriteOffConsultation(
	ctx context.Context,
	id string,
) (*AgreementCustomizationStateResponse, error) {
	return bbapi.Patch[*AgreementCustomizationStateResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointActivateOperationalQuery, id),
		nil,
	)
}

// DeactivateOperationalWriteOffConsultation disables same-day operational write-off consultation.
func (c *Client) DeactivateOperationalWriteOffConsultation(
	ctx context.Context,
	id string,
) (*AgreementCustomizationStateResponse, error) {
	return bbapi.Patch[*AgreementCustomizationStateResponse](
		ctx,
		c.Client,
		fmt.Sprintf(endpointDisableOperationalQuery, id),
		nil,
	)
}
