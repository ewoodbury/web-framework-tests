package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type reading struct {
	Timestamp int64   `json:"timestamp"`
	Voltage   float64 `json:"voltage"`
	Current   float64 `json:"current"`
}

type readingData []reading

var readings = readingData{}

var Conn *pgx.Conn

func main() {
	connectToDb()
	e := echo.New()
	e.POST("/data/local", postHandlerLocal)
	e.Logger.Fatal(e.Start("127.0.0.1:8000"))
}

func connectToDb() error {
	err := godotenv.Load()
	url := fmt.Sprintf("postgres://%s:%s@localhost:5432/postgres", os.Getenv("PG_USER"), os.Getenv("PG_PASS"))
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Could not connect to database: %s", err)
	}
	Conn = conn
	return nil
}

type rawData []struct {
	Voltage   float64 `json:"voltage"`
	Current   float64 `json:"current"`
}

func postHandlerLocal(c echo.Context) error {
    data := new(rawData)
	if err := c.Bind(data); err != nil {
		fmt.Println(err)
        return err
    }
	readings = append(readings, reading{Timestamp: 1000 * time.Now().Unix(), Voltage: (*data)[0].Voltage, Current: (*data)[0].Current})
	return c.JSON(http.StatusOK, data)
}

func postHandlerDb(c echo.Context) {
    data := new(rawData)
	if err := c.Bind(data); err != nil {
		fmt.Println(err)
        return
    }
	
	var query strings.Builder
	query.WriteString(`INSERT INTO cell_signals (test_id, measured_at, cell_voltage, cell_current) VALUES`)
	var vals []interface{}

	for i, s := range (*data) {
		query.WriteString(fmt.Sprintf(`(1, current_timestamp, $%d, $%d)`, (i*2+1), (i*2+2)))
		vals = append(vals, s.Voltage, s.Current)
		if i < len(*data) - 1 {query.WriteString(", ")} else {query.WriteString(";")}
	}
	_, execErr := Conn.Exec(context.Background(), query.String(), vals...)
	if execErr != nil {
		fmt.Printf("Insert execution err: %v", execErr)
	}
}