package quest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_AddQuestServeHttp(t *testing.T) {
	var add AddQuest

	r := strings.NewReader(`{"name":"Buy eggs","cost":100}`)
	req := httptest.NewRequest("POST", "http://localhost:8080/addquest", r)
	w := httptest.NewRecorder()
	add.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

	r = strings.NewReader(`{"name":"Buy eggs","cos":100}`)
	req = httptest.NewRequest("POST", "http://localhost:8080/addquest", r)
	w = httptest.NewRecorder()
	add.ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

}
