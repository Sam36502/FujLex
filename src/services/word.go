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
	PARAM_WORD_ID = "word_id"
)

func PageSetWord(c echo.Context) error {
	langID_s := c.Param(PARAM_LANG_ID)
	wordID_s := c.Param(PARAM_WORD_ID)
	isCreate := wordID_s == ""
	redirectUrl := fmt.Sprintf("/lang/%s", langID_s)

	langID, err := strconv.ParseUint(langID_s, 10, 64)
	if err != nil {
		return view.FailRequestWithError(c, "Invalid Lang-ID provided", err, "/")
	}

	var lang client.Lang
	var word client.Word
	tagstring := ""
	if !isCreate {
		redirectUrl = fmt.Sprintf("%s/word/set/%s", redirectUrl, wordID_s)
		wordID, err := strconv.ParseUint(wordID_s, 10, 64)
		if err != nil {
			return view.FailRequestWithError(c, "Invalid Word-ID provided", err, redirectUrl)
		}
		word, err = client.GetWordByID(langID, wordID)
		if err != nil {
			return view.FailRequestWithError(c, "No word with that ID found", err, redirectUrl)
		}
		lang = *word.Language

		// Get tag string
		for _, t := range word.Tags {
			tagstring = fmt.Sprintf("%s, %s", tagstring, t.Tag)
		}
	} else {
		lang, err = client.GetLangByID(langID)
		if err != nil {
			return view.FailRequestWithError(c, fmt.Sprintf("Failed to get info for language with ID %d", langID), err, redirectUrl)
		}
	}

	return view.RenderTemplate(
		c, view.TMP_WORD_UPDATE,
		view.Data{
			"is_create": isCreate,
			"word":      word,
			"plang":     lang,
			"tags":      tagstring,
			"show_IPA":  true,
		},
	)
}

func HandleSetWord(c echo.Context) error {
	langID_s := c.Param(PARAM_LANG_ID)
	wordID_s := c.Param(PARAM_WORD_ID)
	isCreate := wordID_s == ""
	redirectUrl := fmt.Sprintf("/lang/%s/word/set", langID_s)

	// Check lang-ID
	langID, err := strconv.ParseUint(langID_s, 10, 64)
	if err != nil {
		return view.FailRequestWithError(c, "Invalid Lang-ID provided", err, "/")
	}

	// Set up word object
	var word client.Word
	word.Language = &client.Lang{}
	word.Language.ID = &langID
	word.ID = nil
	if !isCreate {
		wordID, err := strconv.ParseUint(wordID_s, 10, 64)
		if err != nil {
			return view.FailRequestWithError(c, "Invalid Word-ID provided", err, redirectUrl)
		}
		redirectUrl += "/" + wordID_s
		word.ID = &wordID
	}

	// Validate fields
	values, err := c.FormParams()
	if err != nil {
		return view.FailRequestWithError(c, "Failed to retrieve form params", err, redirectUrl)
	}

	if !values.Has("roman") || values.Get("roman") == "" {
		redirectUrl = fmt.Sprintf("%s?%s",
			redirectUrl,
			view.ErrorNotice("Field invalid", "The \"Romanisation\" field is required."))
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	word.Romanisation = values.Get("roman")

	if !values.Has("ipa") || values.Get("ipa") == "" {
		redirectUrl = fmt.Sprintf("%s?%s",
			redirectUrl,
			view.ErrorNotice("Field invalid", "The \"Pronunciation\" field is required."))
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	word.Pronunciation = values.Get("ipa")

	if !values.Has("gloss") || values.Get("gloss") == "" {
		redirectUrl = fmt.Sprintf("%s?%s",
			redirectUrl,
			view.ErrorNotice("Field invalid", "The \"Meanings\" field is required."))
		return c.Redirect(http.StatusFound, redirectUrl)
	}
	word.Meanings = strings.Split(values.Get("gloss"), ",")
	for i, m := range word.Meanings {
		word.Meanings[i] = strings.TrimSpace(m)
	}
	word.Orthography = values.Get("ortho")

	// Get Tags
	word.Tags = []*client.Tag{}
	for _, tcode := range strings.Split(values.Get("tags"), ",") {
		tcode = strings.TrimSpace(tcode)
		if tcode == "" {
			continue
		}
		tag, err := client.GetTagByCode(tcode)
		if err != nil {
			redirectUrl = fmt.Sprintf("%s?%s",
				redirectUrl,
				view.ErrorNotice("Unknown Tag", fmt.Sprintf("The tag '%s' is unrecognised.", tcode)))
			return c.Redirect(http.StatusFound, redirectUrl)
		}
		word.Tags = append(word.Tags, &tag)
	}

	word.Etymology = values.Get("etym")
	word.Notes = values.Get("notes")

	// Send updated/new lang info
	_, err = client.PutWord(word)
	if err != nil {
		return view.FailRequestWithError(c, "Failed to send new word info", err, redirectUrl)
	}

	// Succeeded, redirect to created language page
	msg := "The word was successfully updated"
	if isCreate {
		msg = "The word was successfully added"
	}
	redirectUrl = fmt.Sprintf("/lang/%d?%s=%s&%s",
		langID,
		PARAM_SEARCH_QUERY,
		word.Romanisation,
		view.SuccessNotice("Success!", msg))
	return c.Redirect(http.StatusFound, redirectUrl)
}
