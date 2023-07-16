package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
			return view.FailRequestWithError(c, fmt.Sprintf("Query '%s' failed", query), err, fmt.Sprint("/lang/", langID))
		}
	}

	return view.RenderTemplate(
		c, view.TMP_LANG_DETAIL,
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
	var langAncestorCode string // Only using one in the front-end
	if !isCreate {
		langID, err := strconv.ParseUint(langID_s, 10, 64)
		if err != nil {
			return view.FailRequestWithError(c, "Invalid Lang-ID provided", err, "/")
		}
		lang, err = client.GetLangByID(langID)
		if err != nil {
			return view.FailRequestWithError(c, fmt.Sprintf("Failed to get info for language with ID %d", langID), err, "/")
		}
		if len(lang.AncestorCodes) > 0 {
			langAncestorCode = lang.AncestorCodes[0]
		}
	}

	// Get list of languages for dropdown
	langlist, err := client.GetAllLangs()
	if err != nil {
		return view.FailRequestWithError(c, "Failed to get language list", err, "/")
	}

	return view.RenderTemplate(
		c, view.TMP_LANG_UPDATE,
		view.Data{
			"is_create":     isCreate,
			"lang":          lang,
			"lang_ancestor": langAncestorCode,
			"all_langs":     langlist,
		},
	)
}

func HandleSetLang(c echo.Context) error {
	langID_s := c.Param(PARAM_LANG_ID)
	isCreate := langID_s == ""
	redirectUrl := "/lang/set"

	// Set up lang object
	var lang client.Lang
	lang.ID = nil
	if !isCreate {
		redirectUrl += "/" + langID_s
		id, err := strconv.ParseUint(langID_s, 10, 64)
		if err != nil {
			return view.FailRequestWithError(c, "Invalid Lang-ID provided", err, redirectUrl)
		}
		lang.ID = &id
	}

	// Validate fields
	values, err := c.FormParams()
	if err != nil {
		return view.FailRequestWithError(c, "Failed to retrieve form params", err, redirectUrl)
	}

	langlist, err := client.GetAllLangs()
	if err != nil {
		return view.FailRequestWithError(c, "Failed to get language list", err, "/")
	}

	if !values.Has("name") || values.Get("name") == "" {
		redirectUrl = fmt.Sprintf("%s?%s",
			redirectUrl,
			view.ErrorNotice("Field invalid", "The \"Name\" field is required."))
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	lang.Name = values.Get("name")

	if !values.Has("code") || values.Get("code") == "" {
		redirectUrl = fmt.Sprintf("%s?%s",
			redirectUrl,
			view.ErrorNotice("Field invalid", "The \"Code\" field is required."))
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	lang.Code = values.Get("code")
	lang.Desc = values.Get("desc")

	if values.Has("ancestor") && values.Get("ancestor") != "" {
		if len(lang.AncestorCodes) > 0 {
			lang.AncestorCodes[0] = values.Get("ancestor")
		} else {
			lang.AncestorCodes = append(lang.AncestorCodes, values.Get("ancestor"))
		}
	}

	// Check for unique fields
	if isCreate {
		for _, v := range langlist {
			if strings.TrimSpace(v.Name) == strings.TrimSpace(lang.Name) {
				redirectUrl = fmt.Sprintf("%s?%s",
					redirectUrl,
					view.ErrorNotice("Field invalid", "The language name should be unique. It's already in use by another language."))
				return c.Redirect(http.StatusFound, redirectUrl)
			}
			if strings.TrimSpace(v.Code) == strings.TrimSpace(lang.Code) {
				redirectUrl = fmt.Sprintf("%s?%s",
					redirectUrl,
					view.ErrorNotice("Field invalid", "The language code should be unique. It's already in use by '"+v.Name+"'."))
				return c.Redirect(http.StatusFound, redirectUrl)
			}
		}
	}

	// Send updated/new lang info
	ID, err := client.PutLang(lang)
	if err != nil {
		return view.FailRequestWithError(c, "Failed to send new language info", err, redirectUrl)
	}

	// Succeeded, redirect to created language page
	msg := "The language was successfully updated"
	if isCreate {
		msg = "The language was successfully created"
	}
	redirectUrl = fmt.Sprintf("/lang/%d?%s",
		ID,
		view.SuccessNotice("Success!", msg))
	return c.Redirect(http.StatusFound, redirectUrl)
}
