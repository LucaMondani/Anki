package main

import (
	//"database/sql"
	"fmt"
	"luca/ankii/templates"
	"net/http"
	"strconv"
	"math/rand"
	"log"
	"time"
	"github.com/a-h/templ"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type created_Deck struct {
	Title       string   `form:"title"`
	Description string   `form:"description"`
	Vokabel1    []string `form:"voc1[]"`
	Vokabel2    []string `form:"voc2[]"`
}

type edited_Deck struct {
	Title       string   `form:"title"`
	Description string   `form:"description"`
	Vokabel1    []string `form:"voc1[]"`
	Vokabel2    []string `form:"voc2[]"`
	Id          string   `form:"id"`
	Vokabel_ids []int `form:"vokabel_id[]"`
}

type DB_Deck struct {
	Id          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}
type DB_Card struct {
	Id    int    `db:"id"`
	Front string `db:"front"`
	Back  string `db:"back"`
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
	e.POST("/new-input", newInput)
	e.DELETE("/delete-input", deleteInput)
	e.POST("/create", create)
	e.GET("/edit/:id", edit)
	e.POST("/edited", edited)
	e.DELETE("/delete/:id", delete)
	e.GET("/next-card/:deck_id", nextCard)
	e.GET("/set-card-status/:card_id", setCardStatus)
	e.GET("/get-back/:card-id", getBack)
	e.GET("/get-front/:card-id", getFront)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func home(c echo.Context) error {
	db_decks := []DB_Deck{}
	err := db.Select(&db_decks, "SELECT * FROM decks ORDER BY id")
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	decks := []templates.Deck_info{}
	for i := 0; i < len(db_decks); i++ {
		decks = append(decks, templates.Deck_info{Id: db_decks[i].Id, Title: db_decks[i].Title, Description: db_decks[i].Description})
	}
	time_today:= time.Now()
	date_today := time_today.Format("2006-01-02")
	time_tomorrow := time_today.AddDate(0,0,1)
	date_tomorrow := time_tomorrow.Format("2006-01-02")
	fmt.Println(date_today)
	fmt.Println(date_tomorrow)
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
	send_deck := templates.Deck{Info: templates.Deck_info{Id: deck_info[0].Id, Title: deck_info[0].Title, Description: deck_info[0].Description}, Cards: []templates.Card{}}
	for i := 0; i < len(cards); i++ {
		send_deck.Cards = append(send_deck.Cards, templates.Card{Id: cards[i].Id, Front: cards[i].Front, Back: cards[i].Back})
	}
	fmt.Println(send_deck)
	return HTML(c, templates.ViewDeck(send_deck))
}

func newInput(c echo.Context) error {
	return HTML(c, templates.Vokabel_input_feld())
}

func deleteInput(c echo.Context) error {
	return nil
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
	time_today:= time.Now()
	date_today := time_today.Format("2006-01-02")
	for i := 0; i < len(Deck_created.Vokabel1); i++ {
		tx.MustExec("INSERT INTO cards (Front, Back, Deck_Id, Learn_Date) VALUES ($1, $2, $3, $4)", Deck_created.Vokabel1[i], Deck_created.Vokabel2[i], id, date_today)
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
	send_deck := templates.Deck{Info: templates.Deck_info{Id: deck_info[0].Id, Title: deck_info[0].Title, Description: deck_info[0].Description}, Cards: []templates.Card{}}
	for i := 0; i < len(cards); i++ {
		send_deck.Cards = append(send_deck.Cards, templates.Card{Id: cards[i].Id, Front: cards[i].Front, Back: cards[i].Back})
	}
	fmt.Println(send_deck)
	return HTML(c, templates.Edit(send_deck))
}

func edited(c echo.Context) error {
	var Deck_edited edited_Deck
	err := c.Bind(&Deck_edited)
	fmt.Printf("Deck id is: %s\n", Deck_edited.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(Deck_edited.Vokabel1) != len(Deck_edited.Vokabel2) {
		fmt.Println("Error while editing Deck. Vocab1 and Vocab2 not the same length.")
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	} else if len(Deck_edited.Vokabel1) == 0 {
		fmt.Println("Error while editing Deck. Deck is empty.")
		return echo.NewHTTPError(http.StatusBadRequest, "Empty Deck")
	}

	
	for i := 0; i < len(Deck_edited.Vokabel1); i++ {
		fmt.Printf("%s : %s \n", Deck_edited.Vokabel1[i], Deck_edited.Vokabel2[i])
	}
	tx := db.MustBegin()
	tx.NamedExec("UPDATE decks SET title = :title, description = :description WHERE id = :id", &Deck_edited)

	if len(Deck_edited.Vokabel_ids) == 0 {
		tx.MustExec("DELETE FROM cards WHERE deck_id = $1", Deck_edited.Id)
	} else {
		query, args, err := sqlx.In("DELETE FROM cards WHERE deck_id = (?) AND id NOT IN (?)", Deck_edited.Id, Deck_edited.Vokabel_ids)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		query = db.Rebind(query)
		_, err = tx.Exec(query, args...)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	i := 0
	time_today:= time.Now()
	date_today := time_today.Format("2006-01-02")
	for ; i < len(Deck_edited.Vokabel_ids); i++ {
		tx.MustExec("UPDATE cards SET front = $1, back = $2, Learn_Date = $3 WHERE id = $4 AND deck_id = $5", Deck_edited.Vokabel1[i], Deck_edited.Vokabel2[i], date_today , Deck_edited.Vokabel_ids[i], Deck_edited.Id)
	}
	for ; i < len(Deck_edited.Vokabel1); i++ {
		tx.MustExec("INSERT INTO cards (Front, Back, Deck_Id, Learn_Date) VALUES ($1, $2, $3, $4)", Deck_edited.Vokabel1[i], Deck_edited.Vokabel2[i], Deck_edited.Id, date_today)
	}
	
	tx.Commit()
	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusCreated, "")
}

func delete(c echo.Context) error{
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	fmt.Println(id)
	tx := db.MustBegin()
	tx.MustExec("DELETE FROM decks WHERE id = $1", id)
	tx.MustExec("DELETE FROM cards WHERE deck_id = $1", id)
	tx.Commit()
	c.Response().Header().Set("HX-Redirect", "/")
	return c.String(http.StatusContinue, "")
}

func nextCard(c echo.Context) error{
	deck_id, err := strconv.Atoi(c.Param("deck_id"))
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	cards := []DB_Card{}
	db.Select(&cards, "SELECT id, front, back FROM cards WHERE deck_id = $1", deck_id)
	random_db_card := cards[rand.Intn(len(cards))];
	random_card := templates.Card{Id: random_db_card.Id, Front: random_db_card.Front, Back: random_db_card.Back}
	db_deck := []DB_Deck{}
	error := db.Select(&db_deck, "SELECT * FROM decks WHERE Id = $1", deck_id)
	if error != nil {
		fmt.Println(error)
		return echo.NewHTTPError(http.StatusInternalServerError, error.Error())
	}
	deck_info := templates.Deck_info{Id: db_deck[0].Id, Title: db_deck[0].Title, Description: db_deck[0].Description}
	return HTML(c, templates.Learn(deck_info, random_card))
}

func setCardStatus(c echo.Context) error {
	// status = 0 : einfach
	// status = 1 : ok
	// status = 2 : nochmal
	status := c.QueryParam("status")
	fmt.Printf("status: %s \n", status)
	card_id, err := strconv.Atoi(c.Param("card_id"))
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	fmt.Printf("card_id: %d \n", card_id)
	deck_id := []int{}
	db.Select(&deck_id, "SELECT deck_id FROM cards WHERE id = $1", card_id)
	// Save status
	time_today:= time.Now()
	date_today := time_today.Format("2006-01-02")
	tx := db.MustBegin()
	if status == "0" {
		time_4_days := time_today.AddDate(0,0,4)
		date_4_days := time_4_days.Format("2006-01-02")
		tx.MustExec("UPDATE cards SET Learn_Date = $1 WHERE id = $2", date_4_days, card_id)
	} else if status == "1" {
		time_tomorrow := time_today.AddDate(0,0,1)
		date_tomorrow := time_tomorrow.Format("2006-01-02")
		tx.MustExec("UPDATE cards SET Learn_Date = $1 WHERE id = $2", date_tomorrow, card_id)
	} else if status == "2" {
		tx.MustExec("UPDATE cards SET Learn_Date = $1 WHERE id = $2", date_today, card_id)
	}
	tx.Commit()

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/next-card/%d", deck_id[0]))
	return c.String(http.StatusCreated, "")
}

func getBack(c echo.Context) error {
	card_id, err := strconv.Atoi(c.Param("card-id"))
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	fmt.Println(card_id)
	db_card := []DB_Card{}
	db.Select(&db_card, "SELECT id, front, back FROM cards WHERE id = $1", card_id)
	card := templates.Card{Id: db_card[0].Id, Front: db_card[0].Front, Back: db_card[0].Back}
	return HTML(c, templates.Vokabel_back(card))
}

func getFront(c echo.Context) error {
	card_id, err := strconv.Atoi(c.Param("card-id"))
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Id.")
	}
	fmt.Println(card_id)
	db_card := []DB_Card{}
	db.Select(&db_card, "SELECT id, front, back FROM cards WHERE id = $1", card_id)
	card := templates.Card{Id: db_card[0].Id, Front: db_card[0].Front, Back: db_card[0].Back}
	return HTML(c, templates.Vokabel_front(card))
}

func HTML(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
