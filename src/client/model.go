package client

type Lang struct {
	ID            uint64   `json:"id"`
	Code          string   `json:"code"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	AncestorCodes []string `json:"ancestors"`
}

type Tag struct {
	ID          uint64 `json:"id"`
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Description string `json:"desc"`
}

type Word struct {
	ID            uint64   `json:"id"`
	Orthography   string   `json:"orthography"`
	Romanisation  string   `json:"romanisation"`
	Pronunciation string   `json:"ipa"`
	Meanings      []string `json:"meanings"`
	Tags          []*Tag   `json:"tags"`
	Etymology     string   `json:"etymology"`
	RootWord      uint64   `json:"root"`
	Notes         string   `json:"notes"`
	Language      *Lang    `json:"language"`
}
