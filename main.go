package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", fromAgent)

	e.Logger.Fatal(e.Start(":5678"))
}

func fromAgent(c echo.Context) error {
	return c.String(http.StatusOK, "From Agent")
}
