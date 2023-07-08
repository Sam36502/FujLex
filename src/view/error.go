package view

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	PARAM_ERR_CODE = "code"
	PARAM_ERR_MSG  = "msg"
)

var G_showDebug = false

func FailRequestWithError(c echo.Context, msg string, err error, redirUrl string) error {
	// Redirects to main page on fatal error
	if !G_showDebug {
		c.Redirect(http.StatusTemporaryRedirect, redirUrl)
	}

	errText := ""
	if err != nil {
		errText = err.Error()
	}
	return RenderTemplate(c, "errdebug.twig", Data{
		"msg":          msg,
		"errtext":      errText,
		"redirect_url": redirUrl,
	})
}

func CustomErrorHandler(err error, c echo.Context) {
	// Fetch error info
	code := 0
	msg := ""
	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		msg = fmt.Sprint(httpErr.Message)
	}

	errText := ""
	switch code {
	case 400:
		errText = "It seems you tried to enter something invalid."
	case 401:
		errText = "You are not authorised to view this page."
	case 403:
		errText = "You do not have permission to view this page."
	case 404:
		errText = "It seems you tried to access a page that doesn't exist."
	case 408:
		errText = "Unfortunately, the request took too long."
	case 413:
		errText = "It appears you tried to make coffee with a teapot."
	case 500:
		errText = "Something unknown's gone very wrong. Sorry."
	case 501:
		errText = "Very clever, but this page is ready yet."
	default:
		errText = "I'm not sure what's gone wrong...."
	}

	// Show error page
	RenderTemplate(c, "error.twig", Data{
		"code":         code,
		"msg":          msg,
		"errtext":      errText,
		"redirect_url": "/",
	})
}
