package main

import (
	docs "parser/docs"
	_http "parser/parser/delivery/http"
	"parser/parser/usecase"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	// Initialise echo context for routes
	e := echo.New()
	parserUsecase := usecase.NewParser()
	_http.NewParserHandler(e, parserUsecase)
	docs.NewDocumentation(e, e.Group(""))
	e.Logger.Fatal(e.Start(":" + viper.GetString("APPLICATION_PORT")))
}
