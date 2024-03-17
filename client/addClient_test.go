package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_AddClientServeHttp(t *testing.T) {
	var add AddClient

	r := strings.NewReader(`{"name":"Egor","balance":300}`)
	req := httptest.NewRequest("POST", "http://localhost:8080/addclient", r)
	w := httptest.NewRecorder()
	add.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

	r = strings.NewReader(`{"nam":"Egor","balanc":300}`)
	req = httptest.NewRequest("POST", "http://localhost:8080/addclient", r)
	w = httptest.NewRecorder()
	add.ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

}
