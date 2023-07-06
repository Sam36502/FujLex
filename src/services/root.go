package services

import (
	"fmt"
	"net/http"

	"github.com/Sam36502/FujLex/src/client"
	"github.com/Sam36502/FujLex/src/view"
	"github.com/labstack/echo/v4"
)

func PageRoot(c echo.Context) error {

	langs, err := client.GetAllLangs()
	if err != nil {
		// TODO: EROR
		return c.String(http.StatusInternalServerError, fmt.Sprintf("FAIL: %v", err))
	}

	return view.RenderTemplate(
		c, "root.twig",
		view.Data{
			"languages": langs,
		},
	)
}
