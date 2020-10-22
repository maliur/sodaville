package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type YtItem struct {
	VideoID string
	Title   string
}

var apiBaseUri = "https://googleapis.com/youtube/v3/search"

func SearchYoutube(query, apiKey string) (YtItem, error) {
	if len(query) == 0 {
		return YtItem{}, fmt.Errorf("query cant be empty")
	}

	if len(apiKey) == 0 {
		return YtItem{}, fmt.Errorf("no apiKey provided")
	}

	client := &http.Client{
		Timeout: time.Second + 10,
	}
	params := url.Values{
		"q":          {query},
		"key":        {apiKey},
		"type":       {"video"},
		"part":       {"snippet"},
		"maxResults": {"10"},
		"sort":       {"relevance"},
	}
	apiBaseUri += "?" + params.Encode()

	res, err := client.Get(apiBaseUri)
	if err != nil {
		return YtItem{}, err
	}
	defer res.Body.Close()

	type yResponse struct {
		Items []struct {
			ID struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title string `json:"title"`
			} `json:"snippet"`
		} `json:"items"`
	}

	var y yResponse
	err = json.NewDecoder(res.Body).Decode(&y)
	if err != nil {
		return YtItem{}, fmt.Errorf("could not decode youtube response: %v", err)
	}

	if len(y.Items) == 0 {
		return YtItem{}, fmt.Errorf("no songs found for query %s", query)
	}

	return YtItem{
		VideoID: y.Items[0].ID.VideoId,
		Title:   y.Items[0].Snippet.Title,
	}, nil
}
