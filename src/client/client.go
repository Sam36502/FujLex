package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func getFromURL(path []string, result interface{}, params map[string]string) error {
	// Make request
	currURI := g_clexClient.baseURL.JoinPath(path...)
	if params != nil {
		q := currURI.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		currURI.RawQuery = q.Encode()
	}

	resp, err := http.Get(currURI.String())
	if err != nil {
		return err
	}

	// Parse reponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

func GetAllLangs() ([]Lang, error) {
	var langs []Lang
	err := getFromURL(
		[]string{"lang"},
		&langs,
		nil,
	)
	return langs, err
}

func GetLangByID(langID uint64) (Lang, error) {
	var lang Lang
	err := getFromURL(
		[]string{"lang", strconv.FormatInt(int64(langID), 10)},
		&lang,
		nil,
	)
	return lang, err
}

func SearchWords(langID uint64, query string) ([]Word, error) {
	var words []Word
	err := getFromURL(
		[]string{"lang", strconv.FormatInt(int64(langID), 10), "search"},
		&words,
		map[string]string{
			"q": query,
		},
	)
	return words, err
}
