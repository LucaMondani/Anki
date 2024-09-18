package main

import (
	"luca/ankii/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", home)
	e.GET("/new", new)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func home(c echo.Context) error {
	return HTML(c, templates.Home())
}
func new(c echo.Context) error {
	return HTML(c, templates.New())
}
func HTML(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
