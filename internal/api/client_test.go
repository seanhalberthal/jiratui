package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_BasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			t.Errorf("expected Basic auth, got %q", auth)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept: application/json, got %q", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"displayName":"test"}`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Username: "user", Token: "token", Auth: AuthBasic})
	resp, err := c.Get(context.Background(), V2("/myself"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := DecodeResponse[MeResponse](resp)
	if err != nil {
		t.Fatal(err)
	}
	if result.DisplayName != "test" {
		t.Errorf("got %q, want %q", result.DisplayName, "test")
	}
}

func TestClient_BearerAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer my-pat" {
			t.Errorf("expected 'Bearer my-pat', got %q", auth)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"displayName":"bearer-user"}`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Token: "my-pat", Auth: AuthBearer})
	resp, err := c.Get(context.Background(), V2("/myself"))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
}

func TestClient_PostSetsContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %q", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Username: "u", Token: "t", Auth: AuthBasic})
	resp, err := c.Post(context.Background(), V2("/issue/TEST-1/comment"), map[string]string{"body": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
}

func TestClient_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid credentials"}`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Username: "u", Token: "bad", Auth: AuthBasic})
	resp, err := c.Get(context.Background(), V2("/myself"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodeResponse[MeResponse](resp)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should contain status code, got: %v", err)
	}
}

func TestCheckResponse_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Username: "u", Token: "t", Auth: AuthBasic})
	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatal(err)
	}
	if err := CheckResponse(resp); err != nil {
		t.Errorf("expected no error for 204, got: %v", err)
	}
}

func TestCheckResponse_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`forbidden`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Username: "u", Token: "t", Auth: AuthBasic})
	resp, err := c.Get(context.Background(), "/test")
	if err != nil {
		t.Fatal(err)
	}
	err = CheckResponse(resp)
	if err == nil {
		t.Fatal("expected error for 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should contain status code, got: %v", err)
	}
}

func TestPathConstructors(t *testing.T) {
	if got := V1("/board/1"); got != "/rest/agile/1.0/board/1" {
		t.Errorf("V1 = %q", got)
	}
	if got := V2("/myself"); got != "/rest/api/2/myself" {
		t.Errorf("V2 = %q", got)
	}
	if got := V3("/search/jql"); got != "/rest/api/3/search/jql" {
		t.Errorf("V3 = %q", got)
	}
}
