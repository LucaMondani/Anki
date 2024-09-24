package main

import (
	//"database/sql"
	"fmt"
	"luca/ankii/templates"
	"net/http"
	"strconv"

	"log"

	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type created_Deck struct {
	Title string `form:"title"`
	Description string `form:"description"` 
	Vokabel1 []string `form:"voc1[]"`
	Vokabel2 []string `form:"voc2[]"`
}

type DB_Deck struct {
	Id int	`db:"id"`
	Title string	`db:"title"`
	Description string	`db:"description"`
}
type DB_Card struct {
	Id int	`db:"id"`
	Front string `db:"front"`
	Back string	`db:"back"`
}

var db *sqlx.DB

func main() {
	var err error
	db, err = sqlx.Connect("postgres", "user=dev dbname=app password=pwd sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }
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
	db_decks := []DB_Deck{}
	err := db.Select(&db_decks, "SELECT * FROM decks")
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	decks := []templates.Deck_info{}
	for i := 0; i < len(db_decks); i++ {
		decks = append(decks, templates.Deck_info{Id: db_decks[i].Id, Title: db_decks[i].Title, Description: db_decks[i].Description})
	}
	
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
	deck_info := []DB_Deck{}
	error := db.Select(&deck_info, "SELECT * FROM decks WHERE Id = $1", id)
	if error != nil {
		fmt.Println(error)
		return echo.NewHTTPError(http.StatusInternalServerError, error.Error())
	}
	cards := []DB_Card{}
	error2 := db.Select(&cards, "SELECT id, front, back FROM cards WHERE Deck_Id = $1", id)
	if error2 != nil {
		fmt.Println(error2)
		return echo.NewHTTPError(http.StatusInternalServerError, error2.Error())
	}
	send_deck := templates.Deck{Info: templates.Deck_info{Id:deck_info[0].Id, Title: deck_info[0].Title, Description: deck_info[0].Description}, Cards: []templates.Card{}}
	for i := 0; i < len(cards); i++ {
		send_deck.Cards = append(send_deck.Cards, templates.Card{Id: cards[i].Id, Front: cards[i].Front, Back: cards[i].Back})
	}
	fmt.Println(send_deck)
	return HTML(c, templates.ViewDeck(send_deck))
}

func learn(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	deck_info := []DB_Deck{}
	error := db.Select(&deck_info, "SELECT * FROM decks WHERE Id = $1", id)
	if error != nil {
		fmt.Println(error)
		return echo.NewHTTPError(http.StatusInternalServerError, error.Error())
	}
	cards := []DB_Card{}
	error2 := db.Select(&cards, "SELECT id, front, back FROM cards WHERE Deck_Id = $1", id)
	if error2 != nil {
		fmt.Println(error2)
		return echo.NewHTTPError(http.StatusInternalServerError, error2.Error())
	}
	send_deck := templates.Deck{Info: templates.Deck_info{Id:deck_info[0].Id, Title: deck_info[0].Title, Description: deck_info[0].Description}, Cards: []templates.Card{}}
	for i := 0; i < len(cards); i++ {
		send_deck.Cards = append(send_deck.Cards, templates.Card{Id: cards[i].Id, Front: cards[i].Front, Back: cards[i].Back})
	}
	fmt.Println(send_deck)
	return HTML(c, templates.Learn(send_deck))
}

func newInput(c echo.Context) error {
	return HTML(c, templates.Vokabel_input_feld())
}

func create(c echo.Context) error {
	var Deck_created created_Deck
	err := c.Bind(&Deck_created)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if len(Deck_created.Vokabel1) != len(Deck_created.Vokabel2) {
		fmt.Println("Error while creating new Deck. Vocab1 and Vocab2 not the same length.")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	} else if len(Deck_created.Vokabel1) == 0 {
		fmt.Println("Error while creating new Deck. Deck is empty.")
		return echo.NewHTTPError(http.StatusBadRequest, "Empty Deck")
	}
	fmt.Printf("title: %v, description: %v\n", Deck_created.Title, Deck_created.Description)
	for i := 0; i < len(Deck_created.Vokabel1); i++ {
		fmt.Printf("%v : %v \n", Deck_created.Vokabel1[i], Deck_created.Vokabel2[i])
	}
	tx := db.MustBegin()
	var id int
	rows, err := tx.NamedQuery("INSERT INTO decks (Title, Description) VALUES (:title, :description) RETURNING Id", Deck_created)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if rows.Next() {
		rows.Scan(&id)
	}
	rows.Close()
	fmt.Printf("id: %d \n", id)
	for i := 0; i < len(Deck_created.Vokabel1); i++ {
		tx.MustExec("INSERT INTO cards (Front, Back, Deck_Id) VALUES ($1, $2, $3)", Deck_created.Vokabel1[i], Deck_created.Vokabel2[i], id)
	}
	tx.Commit()
	
	
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
	var Deck_created created_Deck
	err := c.Bind(&Deck_created)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if len(Deck_created.Vokabel1) != len(Deck_created.Vokabel2) {
		fmt.Println("Error while creating new Deck. Vocab1 and Vocab2 not the same length.")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	} else if len(Deck_created.Vokabel1) == 0 {
		fmt.Println("Error while creating new Deck. Deck is empty.")
		return echo.NewHTTPError(http.StatusBadRequest, "Empty Deck")
	}
	for i := 0; i < len(Deck_created.Vokabel1); i++ {
		fmt.Printf("%v : %v \n", Deck_created.Vokabel1[i], Deck_created.Vokabel2[i])
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
