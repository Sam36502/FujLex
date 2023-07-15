package view

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tyler-sommer/stick"
	"github.com/tyler-sommer/stick/twig"
)

const (
	TEMPLATE_DIR = "tmpl/"
)

type Data map[string]stick.Value

var g_tempEnv *stick.Env

func InitTemplates() {
	ldr := stick.NewFilesystemLoader(TEMPLATE_DIR)
	g_tempEnv = twig.New(ldr)

	// Add twig tests
	g_tempEnv.Tests["defined"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) bool {
		return (val != nil) && (val != "")
	}

	// Add filters
	g_tempEnv.Filters["idp"] = func(ctx stick.Context, val stick.Value, args ...stick.Value) stick.Value {
		if p, ok := val.(*uint64); ok {
			return *p
		}

		return val
	}
}

func RenderTemplate(c echo.Context, tmplPath string, d Data) error {
	ntc := GetNotice(c)
	if ntc != nil {
		d[NTC_PARAM] = *ntc
	}

	err := g_tempEnv.Execute(tmplPath, c.Response(), d)

	if err != nil {
		return c.HTML(
			http.StatusInternalServerError,
			fmt.Sprintf("<h2 style='color:red'>Error: failed to render template '%s'</h2><h3>Error:</h3><p>%v</p>", tmplPath, err),
		)
	}

	return c.NoContent(http.StatusOK)
}
