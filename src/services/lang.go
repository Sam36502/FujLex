package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Sam36502/FujLex/src/client"
	"github.com/Sam36502/FujLex/src/view"
	"github.com/labstack/echo/v4"
)

const (
	PARAM_LANG_ID = "lang_id"
)

func PageSearch(c echo.Context) error {
	langID, err := strconv.ParseUint(c.Param(PARAM_LANG_ID), 10, 64)
	if err != nil {
		// TODO: EROR
		return c.String(http.StatusInternalServerError, fmt.Sprintf("FAIL: %v", err))
	}

	lang, err := client.GetLangByID(langID)
	if err != nil {
		// TODO: EROR
		return c.String(http.StatusInternalServerError, fmt.Sprintf("FAIL: %v", err))
	}

	words, err := client.SearchWords(langID, c.QueryParam("q"))
	if err != nil {
		// TODO: EROR
		return c.String(http.StatusInternalServerError, fmt.Sprintf("FAIL: %v", err))
	}

	return view.RenderTemplate(
		c, "search.twig",
		view.Data{
			"words": words,
			"lang":  lang,
		},
	)
}
