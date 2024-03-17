package quest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_CompleteQuestServeHttp(t *testing.T) {
	var complete CompleteQuest

	r := strings.NewReader(`{"client_id":"1","quest_id":"1"}`)
	req := httptest.NewRequest("POST", "http://localhost:8080/completequest", r)
	w := httptest.NewRecorder()
	complete.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

	r = strings.NewReader(`{"client":"1","quest":"1"}`)
	req = httptest.NewRequest("POST", "http://localhost:8080/completequest", r)
	w = httptest.NewRecorder()
	complete.ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code. Actual status code %v", resp.StatusCode)
	}

}
