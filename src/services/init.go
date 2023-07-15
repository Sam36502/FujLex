package services

import (
	"github.com/Sam36502/FujLex/src/view"
	"github.com/labstack/echo/v4"
)

func Initialise(e *echo.Echo) {

	// Static dir
	e.Static("/static", "static")

	// Error Pages
	e.HTTPErrorHandler = view.CustomErrorHandler

	// Root page
	e.GET("/", PageRoot)

	// Lang Pages
	e.GET("/lang/set", PageSetLang)
	e.POST("/lang/set", HandleSetLang)
	e.GET("/lang/set/:"+PARAM_LANG_ID, PageSetLang)
	e.POST("/lang/set/:"+PARAM_LANG_ID, HandleSetLang)
	e.GET("/lang/:"+PARAM_LANG_ID, PageSearch)

}
