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
// GRU (Federal Government Collection Guide) payments:
//
// Entries without additional data (barcode and value only):
//
//	85880000001380003631130002185001233122022557 →    138.00
//	85850000000200003631130002185002174122025678 →     20.00
//	85800000002660004352882721486900675550002022 →    266.00
//	85830000002660004352882721486900431695002022 →    266.00
//	85800000002713002801874000096214200166000166 →    271.30
//	85860000010000002801874000100210557524000131 →  1,000.00
//	85890000000167402541111200216100039360992860 →     16.74
//	85880000000055802541111100216100023586755805 →      5.58
//
// Entries with reference, competence period, due date and taxpayer:
//
//	89970000000800000010109552316288320117811508 | ref. 50103006   | period 11/2022 | due 04/11/2022 | CPF 442.140.732-15 →  80.00
//	89900000001200000010109552316288320117811755 | ref. 2016021990 | period 11/2022 | due 04/11/2022 | CPF 435.529.512-53 → 120.00
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
			bbapi.ScopeBatchesRequest,
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
	gruDueDate := int64(4112022)     // 04/11/2022 in ddmmaaaa numeric format.

	batch, err := client.CreateGRUBatch(ctx, &batchpayments.CreateGRUBatchRequest{
		RequestNumber:     examples.RandomReqNumber(),
		Agency:            examples.Ptr[int64](1607),
		Account:           examples.Ptr[int64](99738672),
		AccountCheckDigit: examples.Ptr("X"),
		Entries: []batchpayments.GRUEntry{
			{
				Barcode:      "85880000001380003631130002185001233122022557",
				PaymentDate:  scheduledDate,
				PaymentValue: 138.00,
			},
			{
				Barcode:      "85850000000200003631130002185002174122025678",
				PaymentDate:  scheduledDate,
				PaymentValue: 20.00,
			},
			{
				Barcode:      "85800000002660004352882721486900675550002022",
				PaymentDate:  scheduledDate,
				PaymentValue: 266.00,
			},
			{
				Barcode:      "85830000002660004352882721486900431695002022",
				PaymentDate:  scheduledDate,
				PaymentValue: 266.00,
			},
			{
				Barcode:      "85800000002713002801874000096214200166000166",
				PaymentDate:  scheduledDate,
				PaymentValue: 271.30,
			},
			{
				Barcode:      "85860000010000002801874000100210557524000131",
				PaymentDate:  scheduledDate,
				PaymentValue: 1000.00,
			},
			{
				Barcode:      "85890000000167402541111200216100039360992860",
				PaymentDate:  scheduledDate,
				PaymentValue: 16.74,
			},
			{
				Barcode:      "85880000000055802541111100216100023586755805",
				PaymentDate:  scheduledDate,
				PaymentValue: 5.58,
			},
			{
				// ref. 50103006 | period 11/2022 | due 04/11/2022 | CPF 442.140.732-15
				Barcode:             "89970000000800000010109552316288320117811508",
				DueDate:             examples.Ptr[int64](gruDueDate),
				PaymentDate:         gruDueDate,
				PaymentValue:        80.00,
				ReferenceNumber:     examples.Ptr("50103006"),
				CompetenceMonthYear: examples.Ptr[int64](112022),
				TaxpayerID:          examples.Ptr[int64](44214073215),
			},
			{
				// ref. 2016021990 | period 11/2022 | due 04/11/2022 | CPF 435.529.512-53
				Barcode:             "89900000001200000010109552316288320117811755",
				DueDate:             examples.Ptr[int64](gruDueDate),
				PaymentDate:         gruDueDate,
				PaymentValue:        120.00,
				ReferenceNumber:     examples.Ptr("2016021990"),
				CompetenceMonthYear: examples.Ptr[int64](112022),
				TaxpayerID:          examples.Ptr[int64](43552951253),
			},
		},
	})
	if err != nil {
		log.Fatalf("creating GRU batch: %v", err)
	}

	fmt.Printf(" GRU Batch \n")
	fmt.Printf("Request number: %d\n", batch.RequestNumber)
	fmt.Printf("State: %d | Valid payments: %d | Valid value: %.2f\n",
		batch.RequestState,
		batch.ValidTotalCount,
		batch.ValidTotalValue,
	)
	for i, p := range batch.Payments {
		fmt.Printf("  [%d] id=%d receiver=%q accepted=%q errors=%v\n",
			i+1, p.PaymentID, p.ReceiverName, p.AcceptanceIndicator, p.Errors)
	}
}
