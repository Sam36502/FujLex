package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ClexiconClient struct {
	baseURL *url.URL
}

var g_clexClient ClexiconClient

func Initialise(baseUrl string) error {
	base, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}

	g_clexClient = ClexiconClient{
		baseURL: base,
	}
	return nil
}

func GetAllLangs() ([]Lang, error) {

	// Make request
	currURI := g_clexClient.baseURL.JoinPath("lang")
	resp, err := http.Get(currURI.String())
	if err != nil {
		return nil, err
	}

	// Parse reponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var langs []Lang
	err = json.Unmarshal(data, &langs)
	if err != nil {
		return nil, err
	}

	return langs, nil
}
