package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Trip struct {
	From string `json:"pickupLocationName"`
	To   string `json:"returnLocationName"`
	//Date string `json:"date"`
}

func Start() {
	r := gin.Default()
	r.GET("/trips", func(ctx *gin.Context) {
		originQuery := strings.ToLower(ctx.Query("origin"))
		destQuery := strings.ToLower(ctx.Query("destination"))

		allTrips, err := fetchTrips()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if originQuery == "" && destQuery == "" {
			ctx.JSON(http.StatusOK, allTrips)
			return
		}

		var filtered []Trip
		for _, t := range allTrips {

			matchOrigin := originQuery == "" || strings.Contains(strings.ToLower(t.From), originQuery)
			matchDest := destQuery == "" || strings.Contains(strings.ToLower(t.To), destQuery)

			if matchOrigin && matchDest {
				filtered = append(filtered, t)
			}
		}

		ctx.JSON(http.StatusOK, filtered)
	})

	r.Run(":8080")
	log.Println("Server started on port 8080")
}

func fetchTrips() ([]Trip, error) {
	resp, err := http.Get("https://www.hertzfreerider.se/api/transport-routes/?country=SWEDEN")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var trips []Trip
	if err := json.NewDecoder(resp.Body).Decode(&trips); err != nil {
		return nil, err
	}

	return trips, nil
}
