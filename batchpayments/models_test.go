package bbapi_test

import (
	"encoding/json"
	"testing"

	bbapi "github.com/raykavin/bbapi-go"
)

func TestCreatePixTransferBatchRequestMarshal(t *testing.T) {
	description := "Payroll"
	email := "employee@example.com"

	payload, err := json.Marshal(bbapi.CreatePixTransferBatchRequest{
		RequestNumber: 77,
		PaymentType:   bbapi.PaymentTypeSalary,
		Transfers: []bbapi.PixTransfer{
			{
				Date:               15042026,
				Value:              199.90,
				IdentificationMode: 2,
				Email:              &email,
				PaymentDescription: &description,
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if decoded["numeroRequisicao"] != float64(77) {
		t.Fatalf("unexpected numeroRequisicao: %#v", decoded["numeroRequisicao"])
	}
	if decoded["tipoPagamento"] != float64(bbapi.PaymentTypeSalary) {
		t.Fatalf("unexpected tipoPagamento: %#v", decoded["tipoPagamento"])
	}

	items, ok := decoded["listaTransferencias"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("unexpected listaTransferencias: %#v", decoded["listaTransferencias"])
	}
}

func TestCreateTransferBatchResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"estadoRequisicao": 1,
		"quantidadeTransferencias": 2,
		"valorTransferencias": 350.5,
		"quantidadeTransferenciasValidas": 2,
		"valorTransferenciasValidas": 350.5,
		"transferencias": [
			{
				"identificadorTransferencia": 11,
				"valorTransferencia": 150.25
			}
		]
	}`)

	var response bbapi.CreateTransferBatchResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response.RequestState != 1 || response.TransferCount != 2 || len(response.Transfers) != 1 {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.Transfers[0].TransferIdentifier != "11" {
		t.Fatalf("unexpected transfer identifier: %+v", response.Transfers[0])
	}
}

func TestGetTransferPaymentResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"id": 15,
		"listaPagamentos": [
			{
				"numeroCOMPE": 1,
				"agenciaCredito": 1234,
				"nomeBeneficiario": "Alice",
				"documentoCredito": 9988
			}
		],
		"listaDevolucao": [
			{
				"codigoMotivo": 321,
				"dataDevolucao": "20260401",
				"valorDevolucao": 10.5
			}
		]
	}`)

	var response bbapi.GetTransferPaymentResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(response.PaymentItems) != 1 {
		t.Fatalf("unexpected payment items: %+v", response.PaymentItems)
	}
	if response.PaymentItems[0].BeneficiaryName != "Alice" {
		t.Fatalf("unexpected payment item: %+v", response.PaymentItems[0])
	}
	if len(response.ReturnItems) != 1 || response.ReturnItems[0].ReasonCode != 321 {
		t.Fatalf("unexpected return items: %+v", response.ReturnItems)
	}
}

func TestGetGRUBatchRequestResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"id": 88,
		"lancamentos": [
			{
				"nomeConvenente": "Tesouro",
				"textoCodigoBarras": "123",
				"pagamento": [
					{
						"id": 7,
						"valor": 45.67,
						"cpfCnpjContribuinte": 12345678901
					}
				],
				"erros": [0]
			}
		]
	}`)

	var response bbapi.GetGRUBatchRequestResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(response.Entries) != 1 {
		t.Fatalf("unexpected entries: %+v", response.Entries)
	}
	entry := response.Entries[0]
	if entry.ConventionName != "Tesouro" || len(entry.Payments) != 1 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	if entry.Payments[0].TaxpayerDocument != 12345678901 {
		t.Fatalf("unexpected nested payment: %+v", entry.Payments[0])
	}
}

func TestGetDARFPaymentResponseUnmarshal(t *testing.T) {
	raw := []byte(`{
		"id": 9,
		"listaPagamentos": [
			{
				"codigo": 5952,
				"tipoContribuinte": 2,
				"identificacaoContribuinte": 12345678901,
				"identificacaoTributo": "IRPF",
				"textoLivre": "quota"
			}
		],
		"listaDevolucao": [
			{
				"codigoMotivo": 99
			}
		]
	}`)

	var response bbapi.GetDARFPaymentResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if len(response.PaymentItems) != 1 {
		t.Fatalf("unexpected payment items: %+v", response.PaymentItems)
	}
	if response.PaymentItems[0].TaxIdentifier != "IRPF" {
		t.Fatalf("unexpected DARF payment item: %+v", response.PaymentItems[0])
	}
	if len(response.ReturnItems) != 1 || response.ReturnItems[0].ReasonCode != 99 {
		t.Fatalf("unexpected DARF return items: %+v", response.ReturnItems)
	}
}
