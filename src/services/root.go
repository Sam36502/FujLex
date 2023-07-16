package services

import (
	"github.com/Sam36502/FujLex/src/client"
	"github.com/Sam36502/FujLex/src/view"
	"github.com/labstack/echo/v4"
)

func PageRoot(c echo.Context) error {

	langs, err := client.GetAllLangs()
	if err != nil {
		return view.FailRequestWithError(c, "Failed to get language list", err, "/")
	}

	return view.RenderTemplate(
		c, view.TMP_ROOT,
		view.Data{
			"languages": langs,
		},
	)
}
