package main

import (
	"io/ioutil"
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
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusNotFound, "Failed")
	}
	return c.String(http.StatusOK, string(body))
}
