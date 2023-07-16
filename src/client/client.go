package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	CLEX_FIELD_ERRCODE = "code"
	CLEX_FIELD_ERRMSG  = "msg"
	CLEX_FIELD_ID      = "id"
)

type ClexiconClient struct {
	baseURL *url.URL
}

type Params map[string]string

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

func getFromURL(path []string, result interface{}, params Params) error {
	// Build request
	currURI := g_clexClient.baseURL.JoinPath(path...)
	if params != nil {
		q := currURI.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		currURI.RawQuery = q.Encode()
	}

	// Send request
	resp, err := http.Get(currURI.String())
	if err != nil {
		return err
	}

	// Read & Parse response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

func putToURL(path []string, obj interface{}) (uint64, error) {
	currURI := g_clexClient.baseURL.JoinPath(path...)

	// Encode data as JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}

	// Send object to server
	rdr := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPut, currURI.String(), rdr)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	// Parse response
	data, err = io.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}
	var errResp ClexError
	err = json.Unmarshal(data, &errResp)
	if err == nil && errResp.Code != "" {
		return 0, errResp
	}

	var idResp IDResponse
	err = json.Unmarshal(data, &idResp)
	if err != nil {
		return 0, err
	}

	return idResp.ID, nil
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
		[]string{"lang"},
		&lang,
		Params{
			"id": strconv.FormatUint(langID, 10),
		},
	)
	return lang, err
}

func PutLang(obj Lang) (uint64, error) {
	return putToURL(
		[]string{"lang"},
		obj,
	)
}

func GetWordByID(langID uint64, wordID uint64) (Word, error) {
	var word Word
	err := getFromURL(
		[]string{"lang", strconv.FormatUint(langID, 10), "word"},
		&word,
		Params{
			"id": strconv.FormatUint(wordID, 10),
		},
	)
	return word, err
}

func PutWord(obj Word) (uint64, error) {
	return putToURL(
		[]string{"lang", strconv.FormatUint(*obj.Language.ID, 10), "word"},
		obj,
	)
}

func SearchWords(langID uint64, query string) ([]Word, error) {
	var words []Word
	err := getFromURL(
		[]string{"lang", strconv.FormatUint(langID, 10), "search"},
		&words,
		Params{
			"q": query,
		},
	)
	for i, l := range words {
		if l.ID == nil {
			words = words[:i]
			break
		}
	}
	return words, err
}

func GetTagByCode(code string) (Tag, error) {
	var tag Tag
	err := getFromURL(
		[]string{"lang"},
		&tag,
		Params{
			"tag": code,
		},
	)
	return tag, err
}
