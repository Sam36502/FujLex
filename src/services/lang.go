package services

import (
	"fmt"
	"strconv"

	"github.com/Sam36502/FujLex/src/client"
	"github.com/Sam36502/FujLex/src/view"
	"github.com/labstack/echo/v4"
)

const (
	PARAM_LANG_ID      = "lang_id"
	PARAM_SEARCH_QUERY = "q"
)

func PageSearch(c echo.Context) error {

	// Get Language info
	langID, err := strconv.ParseUint(c.Param(PARAM_LANG_ID), 10, 64)
	if err != nil {
		return view.FailRequestWithError(c, "Invalid Lang-ID provided", err, "/")
	}
	lang, err := client.GetLangByID(langID)
	if err != nil {
		return view.FailRequestWithError(c, fmt.Sprintf("Failed to get info for language with ID %d", langID), err, "/")
	}

	// Search dictionary if query provided
	var words []client.Word
	query := c.QueryParam(PARAM_SEARCH_QUERY)
	hasQuery := query != ""
	if hasQuery {
		words, err = client.SearchWords(langID, query)
		if err != nil {
			return view.FailRequestWithError(c, fmt.Sprintf("Query '%s' failed:", query), err, fmt.Sprint("/lang/", langID))
		}
	}

	return view.RenderTemplate(
		c, "lang/detail.twig",
		view.Data{
			"has_query": hasQuery,
			"query":     query,
			"words":     words,
			"lang":      lang,
		},
	)
}

func PageSetLang(c echo.Context) error {
	// Get Language info
	langID_s := c.Param(PARAM_LANG_ID)
	isCreate := langID_s == ""

	var lang client.Lang
	if !isCreate {
		langID, err := strconv.ParseUint(langID_s, 10, 64)
		if err != nil {
			return view.FailRequestWithError(c, "Invalid Lang-ID provided", err, "/")
		}
		lang, err = client.GetLangByID(langID)
		if err != nil {
			return view.FailRequestWithError(c, fmt.Sprintf("Failed to get info for language with ID %d", langID), err, "/")
		}
	}

	return view.RenderTemplate(
		c, "lang/setlang.twig",
		view.Data{
			"is_create": isCreate,
			"lang":      lang,
		},
	)
}
