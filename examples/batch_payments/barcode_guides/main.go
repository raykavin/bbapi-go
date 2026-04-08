package main

import (
	"context"
	"fmt"
	"log"
	"os"

	bbapi "github.com/raykavin/bbapi-go"
	"github.com/raykavin/bbapi-go/batchpayments"
	"github.com/raykavin/bbapi-go/examples"
)

// Test data BB homologation environment.
//
// Barcode guide payments:
//
//	83630000000641400052836100812355200812351310 →  64.14
//	83690000001057200052858120735518020735512003 → 105.72
//	83600000003021500052847119156147419156142102 → 302.15
//	84670000001800500470011027860709101194190210 → 180.05
//	89610000000250000010111707200000000000057461 →  25.00
//	89620000000658100010111838900000220203000022 →  65.81
//	84640000001498403132010955706087413535200100 → 149.84
//	82860000000781400181111071029270101202200003 →  78.14
//	84870000000449901602022012514009408900826123 →  44.99
//	85660000000876699122102222230173633469013581 →  87.66
func main() {
	bbClient, err := bbapi.NewClient(bbapi.Config{
		ClientID:     os.Getenv("BB_CLIENT_ID"),
		ClientSecret: os.Getenv("BB_CLIENT_SECRET"),
		AppKey:       os.Getenv("BB_APP_KEY"),
		MTLSCertFile: os.Getenv("BB_CERT_FILE"),
		MTLSKeyFile:  os.Getenv("BB_KEY_FILE"),
		Sandbox:      true,
		Scopes: []bbapi.Scope{
			bbapi.ScopeBarcodeGuidesRequest,
			bbapi.ScopeBarcodeGuidesInfo,
		},
	})
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}

	client, err := batchpayments.NewClient(bbClient)
	if err != nil {
		log.Fatalf("creating batch payments client: %v", err)
	}

	ctx := context.Background()
	scheduledDate := int64(15042026) // 15/04/2026 in ddmmaaaa format.

	batch, err := client.CreateBarcodeGuideBatch(ctx, &batchpayments.CreateBarcodeGuideBatchRequest{
		RequestNumber:          examples.RandomReqNumber(),
		DebitAgencyNumber:      examples.Ptr[int64](1607),
		DebitAccountNumber:     examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		Entries: []batchpayments.BarcodeGuideEntry{
			{
				Barcode:      "83630000000641400052836100812355200812351310",
				PaymentDate:  scheduledDate,
				PaymentValue: 64.14,
			},
			{
				Barcode:      "83690000001057200052858120735518020735512003",
				PaymentDate:  scheduledDate,
				PaymentValue: 105.72,
			},
			{
				Barcode:      "83600000003021500052847119156147419156142102",
				PaymentDate:  scheduledDate,
				PaymentValue: 302.15,
			},
			{
				Barcode:      "84670000001800500470011027860709101194190210",
				PaymentDate:  scheduledDate,
				PaymentValue: 180.05,
			},
			{
				Barcode:      "89610000000250000010111707200000000000057461",
				PaymentDate:  scheduledDate,
				PaymentValue: 25.00,
			},
			{
				Barcode:      "89620000000658100010111838900000220203000022",
				PaymentDate:  scheduledDate,
				PaymentValue: 65.81,
			},
			{
				Barcode:      "84640000001498403132010955706087413535200100",
				PaymentDate:  scheduledDate,
				PaymentValue: 149.84,
			},
			{
				Barcode:      "82860000000781400181111071029270101202200003",
				PaymentDate:  scheduledDate,
				PaymentValue: 78.14,
			},
			{
				Barcode:      "84870000000449901602022012514009408900826123",
				PaymentDate:  scheduledDate,
				PaymentValue: 44.99,
			},
			{
				Barcode:      "85660000000876699122102222230173633469013581",
				PaymentDate:  scheduledDate,
				PaymentValue: 87.66,
			},
		},
	})
	if err != nil {
		log.Fatalf("creating barcode guide batch: %v", err)
	}

	fmt.Printf(" Barcode Guide Batch \n")
	fmt.Printf("Request number: %d\n", batch.RequestNumber)
	fmt.Printf("State: %d | Valid entries: %d | Valid value: %.2f\n",
		batch.StateCode,
		batch.ValidEntryCount,
		batch.ValidEntryValue,
	)
	for i, e := range batch.Entries {
		fmt.Printf("  [%d] id=%s beneficiary=%q accepted=%q errors=%v\n",
			i+1, e.PaymentIdentifier, e.BeneficiaryName, e.AcceptanceIndicator, e.Errors)
	}
}
