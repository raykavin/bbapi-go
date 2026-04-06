// Package bbapi provides a Go SDK for Banco do Brasil APIs.
//
// This first iteration implements the Pagamentos em Lote API while keeping the
// client structure ready for future Banco do Brasil products.
//
// Basic usage:
//
//	client, err := bbapi.NewClient(bbapi.Config{
//	    ClientID:     "your-client-id",
//	    ClientSecret: "your-client-secret",
//	    AppKey:       "your-app-key",
//	    Sandbox:      true,
//	    Scopes: []bbapi.Scope{
//	        bbapi.ScopeTransfersRequest,
//	        bbapi.ScopeBatchesRequest,
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	transfer, err := client.CreateTransferBatch(ctx, &bbapi.CreateTransferBatchRequest{
//	    RequestNumber: 123,
//	    PaymentType:   bbapi.PaymentTypeMiscellaneous,
//	    Transfers: []bbapi.Transfer{
//	        {
//	            TransferDate:  10042026,
//	            TransferValue: 150.75,
//	        },
//	    },
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	_ = transfer
package bbapi

// Scope represents an OAuth2 scope accepted by Banco do Brasil APIs.
type Scope string

const (
	ScopeBankSlipsInfo              Scope = "pagamentos-lote.boletos-info"
	ScopeBankSlipsRequest           Scope = "pagamentos-lote.boletos-requisicao"
	ScopeCancelRequest              Scope = "pagamentos-lote.cancelar-requisicao"
	ScopeReturnedPaymentsInfo       Scope = "pagamentos-lote.devolvidos-info"
	ScopeBarcodeGuidesInfo          Scope = "pagamentos-lote.guias-codigo-barras-info"
	ScopeBarcodeGuidesRequest       Scope = "pagamentos-lote.guias-codigo-barras-requisicao"
	ScopeBatchesInfo                Scope = "pagamentos-lote.lotes-info"
	ScopeBatchesRequest             Scope = "pagamentos-lote.lotes-requisicao"
	ScopeBarcodePaymentsInfo        Scope = "pagamentos-lote.pagamentos-codigo-barras-info"
	ScopeManualGuidePaymentsInfo    Scope = "pagamentos-lote.pagamentos-guias-sem-codigo-barras-info"
	ScopeManualGuidePaymentsRequest Scope = "pagamentos-lote.pagamentos-guias-sem-codigo-barras-requisicao"
	ScopePaymentsInfo               Scope = "pagamentos-lote.pagamentos-info"
	ScopePixInfo                    Scope = "pagamentos-lote.pix-info"
	ScopeTransfersInfo              Scope = "pagamentos-lote.transferencias-info"
	ScopePixTransfersInfo           Scope = "pagamentos-lote.transferencias-pix-info"
	ScopePixTransfersRequest        Scope = "pagamentos-lote.transferencias-pix-requisicao"
	ScopeTransfersRequest           Scope = "pagamentos-lote.transferencias-requisicao"
)

const (
	PaymentRequestStateConsistent          = 1
	PaymentRequestStateInconsistent        = 2
	PaymentRequestStateAllInconsistent     = 3
	PaymentRequestStatePending             = 4
	PaymentRequestStateProcessing          = 5
	PaymentRequestStateProcessed           = 6
	PaymentRequestStateRejected            = 7
	PaymentRequestStatePreparingUnreleased = 8
	PaymentRequestStateReleasedByAPI       = 9
	PaymentRequestStatePreparingReleased   = 10
)

const (
	PaymentTypeSuppliers     = 126
	PaymentTypeSalary        = 127
	PaymentTypeMiscellaneous = 128
)
