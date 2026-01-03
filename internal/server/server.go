package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type Trip struct {
	From string `json:"from"`
	To   string `json:"to"`
	Date string `json:"date"`
}

func Start() {
	r := gin.Default()
	r.GET("/trips", func(c *gin.Context) {
		trips := scrapeTrips()
		c.JSON(http.StatusOK, trips)
	})
	r.Run(":8080")
	log.Println("Server started on port 8080")
}

func scrapeTrips() []Trip {
	var trips []Trip
	c := colly.NewCollector(
		colly.AllowedDomains("www.hertzfreerider.se", "hertzfreerider.se"),
	)

	// For example, if trips are in a table row:
	c.OnHTML("tr.trip-row", func(e *colly.HTMLElement) {
		trip := Trip{
			From: strings.TrimSpace(e.ChildText(".origin")),
			To:   strings.TrimSpace(e.ChildText(".destination")),
			Date: strings.TrimSpace(e.ChildText(".date")),
		}
		trips = append(trips, trip)
	})

	err := c.Visit("https://www.hertzfreerider.se/sv-se/")

	if err != nil {
		log.Fatal(err)
	}

	return trips
}
