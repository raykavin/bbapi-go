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
// DARF (Federal Revenue Collection Document) payments:
//
//	Revenue  Taxpayer                   Id.Code  Assessment  Reference  Principal    Fine   Interest  Total      Due
//	6106     75.224.842/0001-26         18       14/04/2026  1          128.01                      128.01     15/04/2026
//	5952     26.707.621/0001-01         16       14/04/2026  112021     1,116.00     7.36           1,123.36   15/04/2026
//	1708     93.809.477/0001-01         16       14/04/2026             300.00       1.98           301.98     15/04/2026
func main() {
	bbClient, err := bbapi.NewClient(bbapi.Config{
		ClientID:     os.Getenv("BB_CLIENT_ID"),
		ClientSecret: os.Getenv("BB_CLIENT_SECRET"),
		AppKey:       os.Getenv("BB_APP_KEY"),
		Sandbox:      true,
		MTLSCertFile: os.Getenv("BB_CERT_FILE"),
		MTLSKeyFile:  os.Getenv("BB_KEY_FILE"),
		Scopes: []bbapi.Scope{
			bbapi.ScopeManualGuidePaymentsRequest,
			bbapi.ScopeManualGuidePaymentsInfo,
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
	paymentDate := int64(15042026)    // 15/04/2026 in ddmmaaaa format.
	assessmentDate := int64(14042026) // 14/04/2026 in ddmmaaaa format.

	batch, err := client.CreateDARFBatch(ctx, &batchpayments.CreateDARFBatchRequest{
		RequestID:              examples.RandomReqNumber(),
		DebitAgencyNumber:      examples.Ptr[int64](1607),
		DebitAccountNumber:     examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		Entries: []batchpayments.DARFEntry{
			{
				// Revenue 6106 | CNPJ 75.224.842/0001-26 | id code 18
				// Assessment 14/04/2026 | ref. 1 | principal 128.01 | total 128.01 | due 15/04/2026
				PaymentDate:            paymentDate,
				PaymentValue:           128.01,
				TaxRevenueCode:         examples.Ptr[int64](6106),
				TaxpayerTypeCode:       examples.Ptr[int64](18),
				TaxpayerIdentification: examples.Ptr[int64](75224842000126),
				AssessmentDate:         examples.Ptr(assessmentDate),
				ReferenceNumber:        examples.Ptr[int64](1),
				PrincipalValue:         examples.Ptr(128.01),
				DueDate:                examples.Ptr(paymentDate),
			},
			{
				// Revenue 5952 | CNPJ 26.707.621/0001-01 | id code 16
				// Assessment 14/04/2026 | ref. 112021 | principal 1,116.00 | fine 7.36 | total 1,123.36 | due 15/04/2026
				PaymentDate:            paymentDate,
				PaymentValue:           1123.36,
				TaxRevenueCode:         examples.Ptr[int64](5952),
				TaxpayerTypeCode:       examples.Ptr[int64](16),
				TaxpayerIdentification: examples.Ptr[int64](26707621000101),
				AssessmentDate:         examples.Ptr(assessmentDate),
				ReferenceNumber:        examples.Ptr[int64](112021),
				PrincipalValue:         examples.Ptr(1116.00),
				FineValue:              examples.Ptr(7.36),
				DueDate:                examples.Ptr(paymentDate),
			},
			{
				// Revenue 1708 | CNPJ 93.809.477/0001-01 | id code 16
				// Assessment 14/04/2026 | principal 300.00 | fine 1.98 | total 301.98 | due 15/04/2026
				PaymentDate:            paymentDate,
				PaymentValue:           301.98,
				TaxRevenueCode:         examples.Ptr[int64](1708),
				TaxpayerTypeCode:       examples.Ptr[int64](16),
				TaxpayerIdentification: examples.Ptr[int64](93809477000101),
				AssessmentDate:         examples.Ptr(assessmentDate),
				PrincipalValue:         examples.Ptr(300.00),
				FineValue:              examples.Ptr(1.98),
				DueDate:                examples.Ptr(paymentDate),
			},
		},
	})
	if err != nil {
		log.Fatalf("creating DARF batch: %v", err)
	}

	fmt.Printf(" DARF Batch \n")
	fmt.Printf("ID: %d\n", batch.ID)
	fmt.Printf("State: %d | Valid entries: %d | Valid value: %.2f\n",
		batch.StateCode,
		batch.ValidEntryCount,
		batch.ValidEntryValue,
	)
	for i, e := range batch.Entries {
		fmt.Printf("  [%d] id=%d revenue=%d taxpayer=%d accepted=%q errors=%v\n",
			i+1, e.PaymentIdentifier, e.TaxRevenueCode, e.TaxpayerIdentification,
			e.AcceptanceIndicator, e.Errors)
	}
}
