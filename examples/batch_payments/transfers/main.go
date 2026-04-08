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
// Debit account:
//
//	Branch: 1607 | Account: 99738672-X | Payment contract: 731030
//
// Beneficiaries (CPF salary payment):
//
//	CPF 993.919.161-80 → branch 0018 / account 3066-X
//	CPF 988.010.721-71 → branch 0018 / account 5745-2
//	CPF 342.441.521-99 → branch 0018 / account 5789-4
//	CPF 790.603.195-40 → branch 0018 / account 10841-3
//	CPF 934.966.031-86 → branch 0018 / account 18581-2
//
// Beneficiaries (CNPJ supplier payment):
//
//	CNPJ 84.526.081/0001-58 → branch 0551 / account 60.840-8
//	CNPJ 93.983.472/0001-92 → branch 0551 / account 31.771-3
//	CNPJ 96.059.151/0001-94 → branch 0551 / account 12.803-1
//	CNPJ 97.678.083/0001-04 → branch 0551 / account 14.669-2
//	CNPJ 93.809.477/0001-01 → branch 0551 / account 62.114-5
func main() {
	bbClient, err := bbapi.NewClient(bbapi.Config{
		ClientID:     os.Getenv("BB_CLIENT_ID"),
		ClientSecret: os.Getenv("BB_CLIENT_SECRET"),
		AppKey:       os.Getenv("BB_APP_KEY"),
		Sandbox:      true,
		MTLSCertFile: os.Getenv("BB_CERT_FILE"),
		MTLSKeyFile:  os.Getenv("BB_KEY_FILE"),
		Scopes: []bbapi.Scope{
			bbapi.ScopeTransfersRequest,
			bbapi.ScopeBatchesInfo,
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

	// Salary batch (CPF only)
	salaryBatch, err := client.CreateTransferBatch(ctx, &batchpayments.CreateTransferBatchRequest{
		RequestNumber:          examples.RandomReqNumber(),
		PaymentContractNumber:  examples.Ptr[int64](731030),
		DebitAgency:            examples.Ptr[int64](1607),
		DebitAccount:           examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		PaymentType:            bbapi.PaymentTypeSalary,
		Transfers: []batchpayments.Transfer{
			{
				// CPF 993.919.161-80 → branch 0018 / account 3066-X
				COMPENumber:             examples.Ptr[int64](1), // Banco do Brasil
				CreditAgency:            examples.Ptr[int64](18),
				CreditAccount:           examples.Ptr[int64](3066),
				CreditAccountCheckDigit: examples.Ptr("X"),
				BeneficiaryCPF:          examples.Ptr[int64](99391916180),
				TransferDate:            scheduledDate,
				TransferValue:           3500.00,
			},
			{
				// CPF 988.010.721-71 → branch 0018 / account 5745-2
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](18),
				CreditAccount:           examples.Ptr[int64](5745),
				CreditAccountCheckDigit: examples.Ptr("2"),
				BeneficiaryCPF:          examples.Ptr[int64](98801072171),
				TransferDate:            scheduledDate,
				TransferValue:           4200.00,
			},
			{
				// CPF 342.441.521-99 → branch 0018 / account 5789-4
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](18),
				CreditAccount:           examples.Ptr[int64](5789),
				CreditAccountCheckDigit: examples.Ptr("4"),
				BeneficiaryCPF:          examples.Ptr[int64](34244152199),
				TransferDate:            scheduledDate,
				TransferValue:           2800.00,
			},
			{
				// CPF 790.603.195-40 → branch 0018 / account 10841-3
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](18),
				CreditAccount:           examples.Ptr[int64](10841),
				CreditAccountCheckDigit: examples.Ptr("3"),
				BeneficiaryCPF:          examples.Ptr[int64](79060319540),
				TransferDate:            scheduledDate,
				TransferValue:           5100.00,
			},
			{
				// CPF 934.966.031-86 → branch 0018 / account 18581-2
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](18),
				CreditAccount:           examples.Ptr[int64](18581),
				CreditAccountCheckDigit: examples.Ptr("2"),
				BeneficiaryCPF:          examples.Ptr[int64](93496603186),
				TransferDate:            scheduledDate,
				TransferValue:           3900.00,
			},
		},
	})
	if err != nil {
		log.Fatalf("creating salary batch: %v", err)
	}

	fmt.Printf(" Salary Batch \n")
	fmt.Printf("State: %d | Valid transfers: %d | Valid value: %.2f\n",
		salaryBatch.RequestState,
		salaryBatch.ValidTransferCount,
		salaryBatch.ValidTransferValue,
	)
	for i, t := range salaryBatch.Transfers {
		fmt.Printf("  [%d] id=%s accepted=%q errors=%v\n",
			i+1, t.TransferIdentifier, t.AcceptanceIndicator, t.Errors)
	}

	// Supplier batch (CNPJ)
	supplierBatch, err := client.CreateTransferBatch(ctx, &batchpayments.CreateTransferBatchRequest{
		RequestNumber:          examples.RandomReqNumber(),
		PaymentContractNumber:  examples.Ptr[int64](731030),
		DebitAgency:            examples.Ptr[int64](1607),
		DebitAccount:           examples.Ptr[int64](99738672),
		DebitAccountCheckDigit: examples.Ptr("X"),
		PaymentType:            bbapi.PaymentTypeSuppliers,
		Transfers: []batchpayments.Transfer{
			{
				// CNPJ 84.526.081/0001-58 → branch 0551 / account 60.840-8
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](551),
				CreditAccount:           examples.Ptr[int64](60840),
				CreditAccountCheckDigit: examples.Ptr("8"),
				BeneficiaryCNPJ:         examples.Ptr[int64](84526081000158),
				TransferDate:            scheduledDate,
				TransferValue:           12000.00,
			},
			{
				// CNPJ 93.983.472/0001-92 → branch 0551 / account 31.771-3
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](551),
				CreditAccount:           examples.Ptr[int64](31771),
				CreditAccountCheckDigit: examples.Ptr("3"),
				BeneficiaryCNPJ:         examples.Ptr[int64](93983472000192),
				TransferDate:            scheduledDate,
				TransferValue:           8500.00,
			},
			{
				// CNPJ 96.059.151/0001-94 → branch 0551 / account 12.803-1
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](551),
				CreditAccount:           examples.Ptr[int64](12803),
				CreditAccountCheckDigit: examples.Ptr("1"),
				BeneficiaryCNPJ:         examples.Ptr[int64](96059151000194),
				TransferDate:            scheduledDate,
				TransferValue:           3200.00,
			},
			{
				// CNPJ 97.678.083/0001-04 → branch 0551 / account 14.669-2
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](551),
				CreditAccount:           examples.Ptr[int64](14669),
				CreditAccountCheckDigit: examples.Ptr("2"),
				BeneficiaryCNPJ:         examples.Ptr[int64](97678083000104),
				TransferDate:            scheduledDate,
				TransferValue:           6750.00,
			},
			{
				// CNPJ 93.809.477/0001-01 → branch 0551 / account 62.114-5
				COMPENumber:             examples.Ptr[int64](1),
				CreditAgency:            examples.Ptr[int64](551),
				CreditAccount:           examples.Ptr[int64](62114),
				CreditAccountCheckDigit: examples.Ptr("5"),
				BeneficiaryCNPJ:         examples.Ptr[int64](93809477000101),
				TransferDate:            scheduledDate,
				TransferValue:           9300.00,
			},
		},
	})
	if err != nil {
		log.Fatalf("creating supplier batch: %v", err)
	}

	fmt.Printf("\n Supplier Batch \n")
	fmt.Printf("State: %d | Valid transfers: %d | Valid value: %.2f\n",
		supplierBatch.RequestState,
		supplierBatch.ValidTransferCount,
		supplierBatch.ValidTransferValue,
	)
	for i, t := range supplierBatch.Transfers {
		fmt.Printf("  [%d] id=%s accepted=%q errors=%v\n",
			i+1, t.TransferIdentifier, t.AcceptanceIndicator, t.Errors)
	}
}
