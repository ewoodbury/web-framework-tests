package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to this Go benchmark! Try navigating to /data...\n")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		voltage, _ := strconv.ParseFloat(r.FormValue("voltage"), 64)
		current, _ := strconv.ParseFloat(r.FormValue("current"), 64)
		readings = append(readings, reading{Timestamp: 1000 * time.Now().Unix(), Voltage: voltage, Current: current})
		fmt.Println(readings)
	} else if r.Method == "GET" {
		json.NewEncoder(w).Encode(readings)
	} else {
		fmt.Fprintf(w, "Only POST method allowed.")
	}
}

func main() {
	// Hello world, the web server

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/data", postHandler)
	log.Println("Running server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
