package youtube

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchYoutube(t *testing.T) {
	fakeApiKey := "secret"
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fakeApiKey, r.URL.Query().Get("key"))

		jsonString := `
		{
		  "items": [
			{
			  "id": {
				"videoId": "Yw6u6YkTgQ4"
			  },
			  "snippet": {
				"title": "hello world"
			  }
			}
		  ]
		}`
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, jsonString)
	}
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	apiBaseUri = testServer.URL
	item, err := SearchYoutube("hello world", fakeApiKey)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello world", item.Title)
	assert.Equal(t, "Yw6u6YkTgQ4", item.VideoID)
}
