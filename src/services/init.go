package services

import "github.com/labstack/echo/v4"

func Initialise(e *echo.Echo) {

	// Static dir
	e.Static("/static", "static")

	// Root page
	e.GET("/", PageRoot)

}
