package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

type reading struct {
	Timestamp int64   `json:"timestamp"`
	Voltage   float64 `json:"voltage"`
	Current   float64 `json:"current"`
}

type readingData []reading

var readings = readingData{
	{
		Timestamp: 1635130015000,
		Voltage:   2.910,
		Current:   1.12,
	},
}

var id int = 1

var Conn *pgx.Conn

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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(Conn)
	io.WriteString(w, "Welcome to this Go benchmark! Try navigating to /data...\n")
}

func postHandlerLocal(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		voltage, _ := strconv.ParseFloat(r.FormValue("voltage"), 64)
		current, _ := strconv.ParseFloat(r.FormValue("current"), 64)
		readings = append(readings, reading{Timestamp: 1000 * time.Now().Unix(), Voltage: voltage, Current: current})
	} else if r.Method == "GET" {
		json.NewEncoder(w).Encode(readings)
	} else {
		fmt.Fprintf(w, "Only POST and GET methods allowed.")
	}
}

func postHandlerDb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		voltage := r.FormValue("voltage")
		current := r.FormValue("current")
		// Execute insert statement against Postgres db:
		insertStatement := `INSERT INTO cell_signals (id, test_id, measured_at, cell_voltage, cell_current)
		VALUES ($1, 1, current_timestamp, $2, $3)`
		fmt.Println(insertStatement, voltage, current)
		tag, err := Conn.Exec(context.Background(), insertStatement, id, voltage, current)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(tag)
		id++
	} else {
		fmt.Fprintf(w, "Only POST method allowed.")
	}
}

func main() {
	err := connectToDb()
	if err != nil {
		log.Println(err)
		return
	}
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/data/local", postHandlerLocal)
	http.HandleFunc("/data/db", postHandlerDb)
	log.Println("Running server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
