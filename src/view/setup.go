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
}

func RenderTemplate(c echo.Context, tmplPath string, d Data) error {
	err := g_tempEnv.Execute(tmplPath, c.Response(), d)

	if err != nil {
		return c.HTML(
			http.StatusInternalServerError,
			fmt.Sprintf("<h2 style='color:red'>Error: failed to render template '%s'</h2><h3>Error:</h3><p>%v</p>", tmplPath, err),
		)
	}

	return c.NoContent(http.StatusOK)
}
