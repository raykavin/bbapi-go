package main

import (
	"context"
	"fmt"
	"log"
	"os"

	bbapi "github.com/raykavin/bbapi-go"
	"github.com/raykavin/bbapi-go/examples"
)

// Test data BB homologation environment.
//
// Pix transfer via key:
//
//	Type 1 (phone)  : (11) 985732102        → CNPJ 95.127.446/0001-98
//	Type 2 (email)  : hmtestes2@bb.com.br   → CNPJ 95.127.446/0001-98
//	Type 3 (CNPJ)   : 92037500000116        → CNPJ 92.037.500/0001-16
//	Type 4 (random) : 9e881f18-cc66-4fc7-8f2c-a795dbb2bfc1
//	Type 2 (email)  : testqrcode01@bb.com.br → CPF 287.792.958-27
//	Type 3 (CPF)    : 28779295827            → CPF 287.792.958-27
//	Type 4 (random) : d14d32de-b3b9-4c31-9f89-8df2cec92c50
//
// Pix transfer via account data:
//
//	Type 5 | COMPE 1 | account type 1 | branch 4267 | account 1704959-8 | CPF  287.792.958-27
//	Type 5 | COMPE 1 | account type 1 | branch  551 | account   43814-6 | CNPJ 95.127.446/0001-98
func main() {
	client, err := bbapi.NewClient(bbapi.Config{
		ClientID:     os.Getenv("BB_CLIENT_ID"),
		ClientSecret: os.Getenv("BB_CLIENT_SECRET"),
		AppKey:       os.Getenv("BB_APP_KEY"),
		Sandbox:      true,
		MTLSCertFile: os.Getenv("BB_CERT_FILE"),
		MTLSKeyFile:  os.Getenv("BB_KEY_FILE"),
		Scopes: []bbapi.Scope{
			bbapi.ScopePixTransfersRequest,
			bbapi.ScopePixInfo,
		},
	})
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}

	ctx := context.Background()
	scheduledDate := int64(15042026) // 15/04/2026 in ddmmaaaa format.

	batch := &bbapi.CreatePixTransferBatchRequest{
		RequestNumber:          examples.RandomReqNumber(),
		ContractNumber:         examples.Ptr[int64](731030),
		DebitAgency:            examples.Ptr[int64](1607),
		DebitAccount:           examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		PaymentType:            bbapi.PaymentTypeSuppliers,
		Transfers: []bbapi.PixTransfer{
			{
				// Type 1 (phone): (11) 985732102 → CNPJ 95.127.446/0001-98
				Date:               scheduledDate,
				Value:              123.45,
				IdentificationMode: 1,
				PhoneAreaCode:      examples.Ptr[int64](11),
				PhoneNumber:        examples.Ptr[int64](985732102),
				BeneficiaryCNPJ:    examples.Ptr[int64](95127446000198),
			},
			{
				// Type 2 (email): hmtestes2@bb.com.br → CNPJ 95.127.446/0001-98
				Date:               scheduledDate,
				Value:              234.56,
				IdentificationMode: 2,
				Email:              examples.Ptr("hmtestes2@bb.com.br"),
				BeneficiaryCNPJ:    examples.Ptr[int64](95127446000198),
			},
			{
				// Type 3 (CNPJ): 92037500000116 → CNPJ 92.037.500/0001-16
				Date:               scheduledDate,
				Value:              345.67,
				IdentificationMode: 3,
				BeneficiaryCNPJ:    examples.Ptr[int64](92037500000116),
			},
			{
				// Type 4 (random): 9e881f18-cc66-4fc7-8f2c-a795dbb2bfc1
				Date:               scheduledDate,
				Value:              456.78,
				IdentificationMode: 4,
				RandomIdentifier:   examples.Ptr("9e881f18-cc66-4fc7-8f2c-a795dbb2bfc1"),
			},
			{
				// Type 2 (email): testqrcode01@bb.com.br → CPF 287.792.958-27
				Date:               scheduledDate,
				Value:              567.89,
				IdentificationMode: 2,
				Email:              examples.Ptr("testqrcode01@bb.com.br"),
				BeneficiaryCPF:     examples.Ptr[int64](28779295827),
			},
			{
				// Type 3 (CPF): 28779295827 → CPF 287.792.958-27
				Date:               scheduledDate,
				Value:              678.90,
				IdentificationMode: 3,
				BeneficiaryCPF:     examples.Ptr[int64](28779295827),
			},
			{
				// Type 4 (random): d14d32de-b3b9-4c31-9f89-8df2cec92c50
				Date:               scheduledDate,
				Value:              789.01,
				IdentificationMode: 4,
				RandomIdentifier:   examples.Ptr("d14d32de-b3b9-4c31-9f89-8df2cec92c50"),
			},
		},
	}

	// Pix batch via key
	keyBatch, err := client.CreatePixTransferBatch(ctx, batch)
	if err != nil {
		log.Fatalf("creating Pix key batch: %v", err)
	}

	fmt.Printf(" Pix Batch Key \n")
	fmt.Printf("State: %d | Valid transfers: %d | Valid value: %.2f\n",
		keyBatch.RequestState,
		keyBatch.ValidTransferCount,
		keyBatch.ValidTransferValue,
	)
	for i, t := range keyBatch.Transfers {
		fmt.Printf("  [%d] id=%d accepted=%q errors=%v\n",
			i+1, t.PaymentIdentifier, t.AcceptanceIndicator, t.Errors)
	}

	// Pix batch via account data (type 5)
	accountBatch, err := client.CreatePixTransferBatch(ctx, &bbapi.CreatePixTransferBatchRequest{
		RequestNumber:          examples.RandomReqNumber(),
		DebitAgency:            examples.Ptr[int64](1607),
		DebitAccount:           examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		PaymentType:            bbapi.PaymentTypeMiscellaneous,
		Transfers: []bbapi.PixTransfer{
			{
				// COMPE 1 | account type 1 | branch 4267 | account 1704959-8 | CPF 287.792.958-27
				Date:               scheduledDate,
				Value:              650.00,
				IdentificationMode: 5,
				COMPENumber:        examples.Ptr[int64](1), // Banco do Brasil
				AccountType:        examples.Ptr[int64](1), // checking account
				Agency:             examples.Ptr[int64](4267),
				Account:            examples.Ptr[int64](1704959),
				AccountCheckDigit:  examples.Ptr("8"),
				BeneficiaryCPF:     examples.Ptr[int64](28779295827),
			},
			{
				// COMPE 1 | account type 1 | branch 551 | account 43814-6 | CNPJ 95.127.446/0001-98
				Date:               scheduledDate,
				Value:              1800.00,
				IdentificationMode: 5,
				COMPENumber:        examples.Ptr[int64](1),
				AccountType:        examples.Ptr[int64](1),
				Agency:             examples.Ptr[int64](551),
				Account:            examples.Ptr[int64](43814),
				AccountCheckDigit:  examples.Ptr("6"),
				BeneficiaryCNPJ:    examples.Ptr[int64](95127446000198),
			},
		},
	})
	if err != nil {
		log.Fatalf("creating Pix account batch: %v", err)
	}

	fmt.Printf("\n Pix Batch Account Data \n")
	fmt.Printf("State: %d | Valid transfers: %d | Valid value: %.2f\n",
		accountBatch.RequestState,
		accountBatch.ValidTransferCount,
		accountBatch.ValidTransferValue,
	)
	for i, t := range accountBatch.Transfers {
		fmt.Printf("  [%d] id=%d accepted=%q errors=%v\n",
			i+1, t.PaymentIdentifier, t.AcceptanceIndicator, t.Errors)
	}
}
