package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GetClientInfoServeHttp(t *testing.T) {
	var getInfo GetClientInfo

	r := strings.NewReader(`{"id":1}`)
	req := httptest.NewRequest("POST", "http://localhost:8080/getclientinfo", r)
	w := httptest.NewRecorder()
	getInfo.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

	r = strings.NewReader(`{"i":1}`)
	req = httptest.NewRequest("POST", "http://localhost:8080/getclientinfo", r)
	w = httptest.NewRecorder()
	getInfo.ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

}
