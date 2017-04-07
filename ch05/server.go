package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq" // sql database
	"github.com/pkg/errors"
)

// Queryer database/sql compatible query interface
type Queryer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

// Txer database/sql transaction interface
type Txer interface {
	Queryer
	Commit() error
	Rollback() error
}

// DBer database/sql
type DBer interface {
	Queryer
	Begin() (*sql.Tx, error)
	Close() error
	Ping() error
}

// DB database
type DB struct {
	*sql.DB
}

// DBConfig config
type DBConfig struct {
	Host     string
	User     string
	UserPass string
	Port     string
	DBName   string
	SSLMode  string
}

// App custom context
type App struct {
	DB     DBer
	Logger *log.Logger
}

// NewDB creates DB
func NewDB(c *DBConfig) (DBer, error) {
	conStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", c.User, c.DBName, c.SSLMode)
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create db")
	}
	return &DB{db}, nil
}

// NewApp new echo
func NewApp(dbCfg *DBConfig) (*App, error) {
	db, err := NewDB(dbCfg)
	if err != nil {
		return nil, err
	}
	e := &App{
		DB:     db,
		Logger: log.New(os.Stdout, "[echo app] ", log.LstdFlags),
	}
	return e, nil
}

// Message message
type Message struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
}

// Hello say hello
func (app *App) Hello(ctx echo.Context) error {
	var t time.Time
	err := app.DB.QueryRow(`select now()`).Scan(&t)
	if err != nil {
		return err
	}
	msg := Message{
		Text: "hello!",
		Time: t,
	}
	app.Logger.Printf("sending back message: %v", msg)
	l := ctx.Logger()
	l.Debugf("this is echo logger: %v", msg)
	ctx.JSON(http.StatusOK, msg)
	return nil
}

func main() {

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	cfg := &DBConfig{
		Host:    "localhost",
		DBName:  "gotodoit",
		User:    "gotodoit_api",
		Port:    "5432",
		SSLMode: "disable",
	}
	app, err := NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Route => handler
	e.GET("/hello", app.Hello)
	e.GET("/logger", func(ctx echo.Context) error {
		ctx.Logger().Debug("no log")
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
