package main

import (
	"fmt"

	"github.com/Sam36502/FujLex/src/client"
	"github.com/Sam36502/FujLex/src/services"
	"github.com/Sam36502/FujLex/src/view"
	"github.com/labstack/echo/v4"
)

func main() {

	view.InitTemplates()
	err := client.Initialise("http://lexapi.pearcenet.ch")
	if err != nil {
		fmt.Printf("Error: Failed to start the server:\n%v", err)
	}

	// Initialise REST framework
	e := echo.New()
	services.Initialise(e)
	err = e.Start(":1919")
	if err != nil {
		fmt.Printf("Error: Failed to start the server:\n%v", err)
	}

}
