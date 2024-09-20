package main

import (
	"fmt"
	"luca/ankii/templates"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Vocab struct {
	Vokabel1 []string `form:"voc1[]"`
	Vokabel2 []string `form:"voc2[]"`
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
	e.GET("/deck/:id", deck)
	e.GET("/learn/:id", learn)
	e.POST("/new-input", newInput)
	e.POST("/create", create)
	e.GET("/edit/:id", edit)
	e.POST("/edited", edited)
	//e.DELETE("/delete", delete)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func home(c echo.Context) error {
	decks := []templates.Deck{{Info: templates.Deck_info{Id: 1, Title: "Lernset 1", Description: "Das ist das erste Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Hallo", Back: "Hello"}, {Id: 2, Front: "Tschüss", Back: "Bye"}, {Id: 3, Front: "Apfel", Back: "apple"}}},
							  {Info: templates.Deck_info{Id: 2, Title: "Lernset 2", Description: "Das ist das zweite Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Ananas", Back: "pineapple"}}},
							  {Info: templates.Deck_info{Id: 3, Title: "Lernset 3", Description: "Das ist das dritte Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Kirsche", Back: "Cherry"}}},
							  {Info: templates.Deck_info{Id: 4, Title: "Lernset 4", Description: "Das ist das vierte Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Banane", Back: "banana"}}}}
	fmt.Println(decks)
	return HTML(c, templates.Home(decks))
}

func new(c echo.Context) error {
	return HTML(c, templates.New())
}

func deck(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	decks := []templates.Deck{{Info: templates.Deck_info{Id: 1, Title: "Lernset 1", Description: "Das ist das erste Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Hallo", Back: "Hello"}, {Id: 2, Front: "Tschüss", Back: "Bye"}, {Id: 3, Front: "Apfel", Back: "apple"}}},
							  {Info: templates.Deck_info{Id: 2, Title: "Lernset 2", Description: "Das ist das zweite Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Ananas", Back: "pineapple"}}},
							  {Info: templates.Deck_info{Id: 3, Title: "Lernset 3", Description: "Das ist das dritte Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Kirsche", Back: "Cherry"}}},
							  {Info: templates.Deck_info{Id: 4, Title: "Lernset 4", Description: "Das ist das vierte Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Banane", Back: "banana"}}}}
	deck := decks[id -1]
	fmt.Println(deck)
	return HTML(c, templates.ViewDeck(deck))
}

func learn(c echo.Context) error {
	return HTML(c, templates.Learn())
}

func newInput(c echo.Context) error {
	return HTML(c, templates.Vokabel_input_feld())
}

func create(c echo.Context) error {
	var vocabulary Vocab
	err := c.Bind(&vocabulary)
	if err != nil {
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

func edit(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	decks := []templates.Deck{{Info: templates.Deck_info{Id: 1, Title: "Lernset 1", Description: "Das ist das erste Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Hallo", Back: "Hello"}, {Id: 2, Front: "Tschüss", Back: "Bye"}, {Id: 3, Front: "Apfel", Back: "apple"}}},
							  {Info: templates.Deck_info{Id: 2, Title: "Lernset 2", Description: "Das ist das zweite Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Ananas", Back: "pineapple"}}},
							  {Info: templates.Deck_info{Id: 3, Title: "Lernset 3", Description: "Das ist das dritte Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Kirsche", Back: "Cherry"}}},
							  {Info: templates.Deck_info{Id: 4, Title: "Lernset 4", Description: "Das ist das vierte Lernset dieser App."}, Cards: []templates.Card{{Id: 1, Front: "Banane", Back: "banana"}}}}
	deck := decks[id-1]
	if len(deck.Cards) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Empty Deck")
	}
	return HTML(c, templates.Edit(deck))
}

func edited(c echo.Context) error {
	var vocabulary Vocab
	err := c.Bind(&vocabulary)
	if err != nil {
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


/*
	func delete(c echo.Context) error {
		return
	}
*/
func HTML(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
