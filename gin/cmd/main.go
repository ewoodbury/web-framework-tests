package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type reading struct {
	Timestamp int64   `json:"timestamp"`
	Voltage   float64 `json:"voltage"`
	Current   float64 `json:"current"`
}

type readingData []reading

var readings = readingData{}

func main() {
	r := gin.Default()
	r.GET("/", homePage)
	r.GET("/time", timePage)
	r.POST("/data/local", postHandlerLocal)
	r.Run("127.0.0.1:8000")
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
	fmt.Println(readings)
}