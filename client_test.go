package bbapi_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	bbapi "github.com/raykavin/bbapi-go"
)

// newTestClient creates a Client pre-configured to hit the given httptest server.
func newTestClient(t *testing.T, server *httptest.Server) *bbapi.Client {
	t.Helper()
	c, err := bbapi.NewClient(bbapi.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		AppKey:       "test-app-key",
		APIURL:       server.URL,
		AuthURL:      server.URL + "/oauth/token",
		AccessToken:  "test-bearer-token",
		MaxRetries:   0,
		RetryWaitMin: time.Millisecond,
		RetryWaitMax: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func readJSONBody(t *testing.T, body io.Reader, target any) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(target); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}

// TestNewClient_MissingFields verifies that NewClient returns errors for missing required fields.
func TestNew_MissingFields(t *testing.T) {
	cases := []struct {
		name string
		cfg  bbapi.Config
	}{
		{"missing ClientID", bbapi.Config{ClientSecret: "s", AppKey: "k"}},
		{"missing ClientSecret", bbapi.Config{ClientID: "i", AppKey: "k"}},
		{"missing AppKey", bbapi.Config{ClientID: "i", ClientSecret: "s"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := bbapi.NewClient(tc.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestNewClient_Success(t *testing.T) {
	client, err := bbapi.NewClient(bbapi.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		AppKey:       "key",
		AccessToken:  "prefilled-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := client.GetAccessToken(); got != "prefilled-token" {
		t.Fatalf("expected prefilled token, got %q", got)
	}
}

func TestSetGetAccessToken(t *testing.T) {
	c, _ := bbapi.NewClient(bbapi.Config{ClientID: "i", ClientSecret: "s", AppKey: "k"})
	c.SetAccessToken("tok123")
	if got := c.GetAccessToken(); got != "tok123" {
		t.Fatalf("want tok123, got %q", got)
	}
}

func TestAuthenticateClientCredentials(t *testing.T) {
	var (
		gotAuth string
		gotForm url.Values
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			http.NotFound(w, r)
			return
		}

		gotAuth = r.Header.Get("Authorization")
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		gotForm = r.Form

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bbapi.TokenResponse{
			AccessToken: "bearer-xyz",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			Scope:       "pagamentos-lote.transferencias-requisicao",
		})
	}))
	defer srv.Close()

	c, _ := bbapi.NewClient(bbapi.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		AppKey:       "key",
		APIURL:       srv.URL,
		AuthURL:      srv.URL + "/oauth/token",
		Scopes:       []bbapi.Scope{bbapi.ScopeTransfersRequest},
	})

	tr, err := c.Authenticate(context.Background())
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if tr.AccessToken != "bearer-xyz" {
		t.Fatalf("want bearer-xyz, got %q", tr.AccessToken)
	}
	if got := c.GetAccessToken(); got != "bearer-xyz" {
		t.Fatalf("token not cached: got %q", got)
	}
	if gotForm.Get("grant_type") != "client_credentials" {
		t.Fatalf("unexpected grant type: %q", gotForm.Get("grant_type"))
	}
	if gotForm.Get("scope") != string(bbapi.ScopeTransfersRequest) {
		t.Fatalf("unexpected scope: %q", gotForm.Get("scope"))
	}
	if gotAuth == "" {
		t.Fatal("expected basic authorization header")
	}
}

func TestAuthenticate_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"erros": []map[string]string{{"codigo": "401", "mensagem": "invalid_client"}},
		})
	}))
	defer srv.Close()

	c, _ := bbapi.NewClient(bbapi.Config{
		ClientID:     "bad",
		ClientSecret: "bad",
		AppKey:       "key",
		AuthURL:      srv.URL + "/token",
	})

	_, err := c.Authenticate(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !bbapi.IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized, got: %v", err)
	}
}

func TestReleasePayments(t *testing.T) {
	var (
		gotKey  string
		gotAuth string
		gotBody bbapi.ReleasePaymentsRequest
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/liberar-pagamentos" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		gotKey = r.URL.Query().Get("gw-dev-app-key")
		gotAuth = r.Header.Get("Authorization")
		readJSONBody(t, r.Body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bbapi.ReleasePaymentsResponse{ReturnMessage: "released"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	resp, err := c.ReleasePayments(context.Background(), &bbapi.ReleasePaymentsRequest{
		RequestNumber:  123,
		FloatIndicator: "S",
	})
	if err != nil {
		t.Fatalf("ReleasePayments: %v", err)
	}

	if gotKey != "test-app-key" {
		t.Fatalf("gw-dev-app-key not set; got %q", gotKey)
	}
	if gotAuth != "Bearer test-bearer-token" {
		t.Fatalf("Authorization header not set; got %q", gotAuth)
	}
	if gotBody.RequestNumber != 123 || gotBody.FloatIndicator != "S" {
		t.Fatalf("unexpected body: %+v", gotBody)
	}
	if resp.ReturnMessage != "released" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCreateTransferBatch(t *testing.T) {
	var gotBody bbapi.CreateTransferBatchRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/lotes-transferencias" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		readJSONBody(t, r.Body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"estadoRequisicao":                1,
			"quantidadeTransferencias":        1,
			"valorTransferencias":             1500.0,
			"quantidadeTransferenciasValidas": 1,
			"valorTransferenciasValidas":      1500.0,
			"transferencias": []map[string]any{
				{
					"identificadorTransferencia": 42,
					"valorTransferencia":         1500.0,
				},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	resp, err := c.CreateTransferBatch(context.Background(), &bbapi.CreateTransferBatchRequest{
		RequestNumber: 55,
		PaymentType:   bbapi.PaymentTypeMiscellaneous,
		Transfers: []bbapi.Transfer{
			{
				TransferDate:  15042026,
				TransferValue: 1500.0,
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateTransferBatch: %v", err)
	}
	if gotBody.RequestNumber != 55 || gotBody.PaymentType != bbapi.PaymentTypeMiscellaneous {
		t.Fatalf("unexpected request body: %+v", gotBody)
	}
	if resp.RequestState != 1 || resp.TransferCount != 1 || len(resp.Transfers) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestListTransferBatchesQueryParams(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bbapi.ListTransferBatchesResponse{Index: 0})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	agency := int64(1234)
	account := int64(567890)
	digit := "X"
	startDate := int64(1012026)
	endDate := int64(15012026)
	state := int64(4)
	index := int64(2)
	_, err := c.ListTransferBatches(context.Background(), &bbapi.ListTransferBatchesParams{
		DebitAgency:            &agency,
		DebitAccount:           &account,
		DebitAccountCheckDigit: &digit,
		StartDate:              &startDate,
		EndDate:                &endDate,
		RequestState:           &state,
		Index:                  &index,
	})
	if err != nil {
		t.Fatalf("ListTransferBatches: %v", err)
	}
	for key, expected := range map[string]string{
		"agenciaDebito":                  "1234",
		"contaCorrenteDebito":            "567890",
		"digitoVerificadorContaCorrente": "X",
		"dataInicio":                     "1012026",
		"dataFim":                        "15012026",
		"estadoRequisicao":               "4",
		"indice":                         "2",
		"gw-dev-app-key":                 "test-app-key",
	} {
		if gotQuery.Get(key) != expected {
			t.Fatalf("query %s: want %q, got %q", key, expected, gotQuery.Get(key))
		}
	}
}

func TestCreateGRUBatchUsesOpenAPINames(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/pagamentos-gru" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		readJSONBody(t, r.Body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"numeroRequisicao":      9,
			"estadoRequisicao":      1,
			"quantidadeTotal":       1,
			"valorTotal":            50.25,
			"quantidadeTotalValido": 1,
			"valorTotalValido":      50.25,
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	agency := int64(1234)
	account := int64(99999)
	digit := "1"
	resp, err := c.CreateGRUBatch(context.Background(), &bbapi.CreateGRUBatchRequest{
		RequestNumber:     9,
		ContractCode:      ptrInt64(88),
		Agency:            &agency,
		Account:           &account,
		AccountCheckDigit: &digit,
		Entries: []bbapi.GRUEntry{
			{
				Barcode:      "12345678901234567890123456789012345678901234",
				PaymentDate:  15042026,
				PaymentValue: 50.25,
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateGRUBatch: %v", err)
	}
	if _, ok := gotBody["codigoContrato"]; !ok {
		t.Fatalf("expected codigoContrato in request body: %+v", gotBody)
	}
	if _, ok := gotBody["listaRequisicao"]; !ok {
		t.Fatalf("expected listaRequisicao in request body: %+v", gotBody)
	}
	if resp.RequestNumber != 9 || resp.TotalCount != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestListReturnedPayments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"indice":                   0,
			"quantidadeTotalRegistros": 1,
			"quantidadeRegistros":      1,
			"pagamentos": []map[string]any{
				{
					"identificadorPagamento": 123,
					"valorPagamento":         10.50,
				},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	resp, err := c.ListReturnedPayments(context.Background(), &bbapi.ListReturnedPaymentsParams{
		StartDate: 1012026,
		EndDate:   15012026,
		Index:     0,
	})
	if err != nil {
		t.Fatalf("ListReturnedPayments: %v", err)
	}
	if resp.TotalRecordCount != 1 || len(resp.Payments) != 1 || resp.Payments[0].PaymentIdentifier != 123 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAPIError_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"erros":[{"codigo":"404","mensagem":"not found"}]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv)
	_, err := c.GetBatch(context.Background(), "9999")
	if !bbapi.IsNotFound(err) {
		t.Fatalf("expected IsNotFound, got: %v", err)
	}
}

func TestAPIError_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c, _ := bbapi.NewClient(bbapi.Config{
		ClientID:     "i",
		ClientSecret: "s",
		AppKey:       "k",
		APIURL:       srv.URL,
		AuthURL:      srv.URL + "/token",
		AccessToken:  "tok",
		MaxRetries:   -1,
		RetryWaitMin: time.Millisecond,
		RetryWaitMax: time.Millisecond,
	})
	_, err := c.GetBatch(context.Background(), "1")
	if !bbapi.IsServerError(err) {
		t.Fatalf("expected IsServerError, got: %v", err)
	}
}

func ptrInt64(value int64) *int64 {
	return &value
}
