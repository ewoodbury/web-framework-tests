package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

var readings = readingData{}

type rawData []struct {
	Voltage   float64 `json:"voltage"`
	Current   float64 `json:"current"`
}


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
	io.WriteString(w, "Welcome to this Go benchmark! Try navigating to /data...\n")
}

func postHandlerLocal(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Unpack request JSON data using io:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("io ReadAll err: %v", err)
		}
		// Parse unpacked data into new struct
		var data rawData
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Printf("Json Unmarshal err: %v", err)
			return
		}
		readings = append(readings, reading{Timestamp: 1000 * time.Now().Unix(), Voltage: data[0].Voltage, Current: data[0].Current})
	} else if r.Method == "GET" {
		json.NewEncoder(w).Encode(readings)
	} else {
		fmt.Fprintf(w, "Only POST and GET methods allowed.")
	}
}

func postHandlerDb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Unpack request JSON data using io:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("io ReadAll err: %v", err)
		}
		// Parse unpacked data into new struct
		var data rawData
		err = json.Unmarshal(b, &data)
		if err != nil {
			fmt.Printf("Json Unmarshal err: %v", err)
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
