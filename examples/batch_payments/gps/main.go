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
// GPS (Social Security Guide) payments:
//
//	Revenue  Taxpayer                   Id.Code  Period   INSS       Other  Adjustment  Total
//	1007     126.792.028-99             17       10/2022    700.00                     700.00
//	1007     517.960.368-46             17       10/2022    133.32                     133.32
//	1007     229.028.691-50             17       10/2022    380.00                     380.00
//	4308     74.910.037/0001-93         17       10/2022    706.90                     706.90
//	4308     98.959.112/0001-79         17       10/2022  1,443.75                   1,443.75
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
	paymentDate := int64(1102026) // 01/10/2026 in ddmmaaaa numeric format.

	batch, err := client.CreateGPSBatch(ctx, &batchpayments.CreateGPSBatchRequest{
		RequestNumber:          examples.RandomReqNumber(),
		DebitAgencyNumber:      examples.Ptr[int64](1607),
		DebitAccountNumber:     examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		Entries: []batchpayments.GPSEntry{
			{
				// Revenue 1007 | CPF 126.792.028-99 | id code 17 | period 10/2022
				PaymentDate:            paymentDate,
				PaymentValue:           700.00,
				TaxRevenueCode:         examples.Ptr[int64](1007),
				TaxpayerTypeCode:       examples.Ptr[int64](17),
				TaxpayerIdentification: examples.Ptr[int64](12679202899),
				CompetenceMonthYear:    examples.Ptr[int64](102022),
				INSSValue:              examples.Ptr(700.00),
			},
			{
				// Revenue 1007 | CPF 517.960.368-46 | id code 17 | period 10/2022
				PaymentDate:            paymentDate,
				PaymentValue:           133.32,
				TaxRevenueCode:         examples.Ptr[int64](1007),
				TaxpayerTypeCode:       examples.Ptr[int64](17),
				TaxpayerIdentification: examples.Ptr[int64](51796036846),
				CompetenceMonthYear:    examples.Ptr[int64](102022),
				INSSValue:              examples.Ptr(133.32),
			},
			{
				// Revenue 1007 | CPF 229.028.691-50 | id code 17 | period 10/2022
				PaymentDate:            paymentDate,
				PaymentValue:           380.00,
				TaxRevenueCode:         examples.Ptr[int64](1007),
				TaxpayerTypeCode:       examples.Ptr[int64](17),
				TaxpayerIdentification: examples.Ptr[int64](22902869150),
				CompetenceMonthYear:    examples.Ptr[int64](102022),
				INSSValue:              examples.Ptr(380.00),
			},
			{
				// Revenue 4308 | CNPJ 74.910.037/0001-93 | id code 17 | period 10/2022
				PaymentDate:            paymentDate,
				PaymentValue:           706.90,
				TaxRevenueCode:         examples.Ptr[int64](4308),
				TaxpayerTypeCode:       examples.Ptr[int64](17),
				TaxpayerIdentification: examples.Ptr[int64](74910037000193),
				CompetenceMonthYear:    examples.Ptr[int64](102022),
				INSSValue:              examples.Ptr(706.90),
			},
			{
				// Revenue 4308 | CNPJ 98.959.112/0001-79 | id code 17 | period 10/2022
				PaymentDate:            paymentDate,
				PaymentValue:           1443.75,
				TaxRevenueCode:         examples.Ptr[int64](4308),
				TaxpayerTypeCode:       examples.Ptr[int64](17),
				TaxpayerIdentification: examples.Ptr[int64](98959112000179),
				CompetenceMonthYear:    examples.Ptr[int64](102022),
				INSSValue:              examples.Ptr(1443.75),
			},
		},
	})
	if err != nil {
		log.Fatalf("creating GPS batch: %v", err)
	}

	fmt.Printf(" GPS Batch \n")
	fmt.Printf("Request number: %d\n", batch.RequestNumber)
	fmt.Printf("State: %d | Valid entries: %d | Valid value: %.2f\n",
		batch.RequestStateCode,
		batch.ValidTotalCount,
		batch.ValidEntryValue,
	)
	for i, e := range batch.Entries {
		fmt.Printf("  [%d] id=%d taxpayer=%d accepted=%q errors=%v\n",
			i+1, e.PaymentIdentifier, e.TaxpayerIdentification,
			e.AcceptanceIndicator, e.Errors)
	}
}
