package scryfall

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func FetchCards(query string) ([]Card, error) {
	encodedQuery := url.QueryEscape(query)

	url := fmt.Sprintf("https://api.scryfall.com/cards/search?q=%s", encodedQuery)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "GoMTGDeckBuilder/0.1")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result ScryfallResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for i := range result.Data {
		result.Data[i].ParseTypeLine()
	}

	return result.Data, nil
}
