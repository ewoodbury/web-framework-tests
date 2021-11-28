package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
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
	r := gin.Default()
	connectToDb()
	pingDb()
	r.GET("/", homePage)
	r.GET("/time", timePage)
	r.POST("/data/local", postHandlerLocal)
	r.POST("/data/db", postHandlerDb)
	r.Run("127.0.0.1:8000")
}

func connectToDb() error {
	err := godotenv.Load()
	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASS")
	url := fmt.Sprintf("postgres://%s:%s@localhost:5432/postgres", pgUser, pgPass)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Could not connect to database: %s", err)
	}
	Conn = conn
	return nil
}

func pingDb() {
	ctx := context.Background()
	if err := Conn.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Database pinged successfully")
}

func homePage(c *gin.Context) {
	c.String(200, "Hello world!")
}

func timePage(c *gin.Context) {
	t := time.Now()
	c.String(200, fmt.Sprintf("The current time is %s", t))
}

type rawData []struct {
	Voltage   float64 `json:"voltage"`
	Current   float64 `json:"current"`
}

func postHandlerLocal(c *gin.Context) {
    var data rawData
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err)
        return
    }
	fmt.Println(data[0])
	// fmt.Println(readings)
	readings = append(readings, reading{Timestamp: 1000 * time.Now().Unix(), Voltage: data[0].Voltage, Current: data[0].Current})
}

func postHandlerDb(c *gin.Context) {
    var data rawData
	if err := c.BindJSON(&data); err != nil {
		fmt.Println(err)
        return
    }
	
	var query strings.Builder
	query.WriteString(`INSERT INTO cell_signals (test_id, measured_at, cell_voltage, cell_current) VALUES`)
	var vals []interface{}
	// Loop through data, and append to query string using incrementing Postgres args (prevent SQL injection)
	for i, s := range data {
		query.WriteString(fmt.Sprintf(`(1, current_timestamp, $%d, $%d)`, (i*2+1), (i*2+2)))
		vals = append(vals, s.Voltage, s.Current)
		if i < len(data) - 1 {query.WriteString(", ")} else {query.WriteString(";")}
	}
	_, execErr := Conn.Exec(context.Background(), query.String(), vals...)
	if execErr != nil {
		fmt.Printf("Insert execution err: %v", execErr)
	}
}