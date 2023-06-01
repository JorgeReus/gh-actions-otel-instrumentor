package discord

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadEvent(t *testing.T, path string) *DiscordWebhook {
	validJson, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err.Error())
	}
	var validEvent DiscordWebhook
	err = json.Unmarshal(validJson, &validEvent)
	if err != nil {
		t.Fatal(err.Error())
	}
	return &validEvent
}

func TestSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	}))
	defer server.Close()

	req := loadEvent(t, "test_data/webhook.json")

  err := Notify(server.URL, *req)

  assert.NoError(t, err)
}

func TestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}))
	defer server.Close()

	req := loadEvent(t, "test_data/webhook.json")

  err := Notify(server.URL, *req)

  assert.Error(t, err)
}
