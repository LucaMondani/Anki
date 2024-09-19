package main

import (
	"luca/ankii/templates"
	"net/http"
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Vocab struct {
	Vokabel1 [] string `form:"voc1[]"`
	Vokabel2 [] string `form:"voc2[]"`
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", home)
	e.GET("/new", new)
	e.GET("/deck", deck)
	e.GET("/learn", learn)
	e.POST("/new-input", newInput)
	e.POST("/create", create)
	//e.POST("/edit", edit)
	//e.DELETE("/delete", delete)

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
func deck(c echo.Context) error {
	return HTML(c, templates.Deck())
}
func learn(c echo.Context) error {
	return HTML(c, templates.Learn())
}
func newInput(c echo.Context) error {
	return HTML(c, templates.Vokabel_input_feld())
}
func create(c echo.Context) error {
	var vocabulary Vocab
	err := c.Bind(&vocabulary); if err!= nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if len(vocabulary.Vokabel1) != len(vocabulary.Vokabel2) {
		fmt.Println("Error while creating new Deck. Vocab1 and Vocab2 not the same length.")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	} else if len(vocabulary.Vokabel1) == 0 {
		fmt.Println("Error while creating new Deck. Deck is empty.")
		return echo.NewHTTPError(http.StatusBadRequest, "Empty Deck")
	} 
	for i := 0; i < len(vocabulary.Vokabel1); i++ {
		fmt.Printf("%v : %v \n", vocabulary.Vokabel1[i], vocabulary.Vokabel2[i])
	}
	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusCreated, "")
}
/*func edit(c echo.Context) error {
	return HTML(c, templates.Edit())
}
func delete(c echo.Context) error {
	return 
}
	*/
func HTML(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
