package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RideDetails struct {
	ID          int    "json:id"
	AvailableAt string "json:availableAt"
	ExpireTime  string "json:expireTime"
}

type FreeriderRoute struct {
	PickupLocationName string        `json:"pickupLocationName"`
	ReturnLocationName string        `json:"returnLocationName"`
	Routes             []RideDetails `json:"routes"`
}

type Trip struct {
	From      string `json:"from"`
	To        string `json:"to"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

const timeLayout = "2006-01-02T15:04:05"
const timeZone = "Europe/Stockholm"

func Start() {
	r := gin.Default()

	r.GET("/trips", func(ctx *gin.Context) {
		originQuery := strings.ToLower(ctx.Query("origin"))
		destQuery := strings.ToLower(ctx.Query("destination"))
		startDateQuery := ctx.Query("startDate")
		endDateQuery := ctx.Query("endDate")

		allTrips, err := fetchTrips()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// If no destination or origin is specified
		if originQuery == "" && destQuery == "" {
			ctx.JSON(http.StatusOK, allTrips)
			return
		}

		var filtered []Trip

		// Parse all trips
		for _, t := range allTrips {
			startDate, _ := time.Parse(timeLayout, t.StartDate)
			endDate, _ := time.Parse(timeLayout, t.EndDate)

			matchOrigin := originQuery == "" || strings.Contains(strings.ToLower(t.From), originQuery)
			matchDest := destQuery == "" || strings.Contains(strings.ToLower(t.To), destQuery)
			matchDate := true

			if startDateQuery != "" {
				startDateQueryTime, _ := parseFlexTime(startDateQuery)
				matchDate = startDate.After(startDateQueryTime)
			}

			if endDateQuery != "" {
				endDateQueryTime, _ := parseFlexTime(endDateQuery)
				matchDate = matchDate && endDate.Before(endDateQueryTime)
			}

			if matchOrigin && matchDest && matchDate {
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

	var freeriderRoute []FreeriderRoute
	if err := json.NewDecoder(resp.Body).Decode(&freeriderRoute); err != nil {
		return nil, err
	}

	var allTrips []Trip
	for _, route := range freeriderRoute {
		for _, trip := range route.Routes {

			allTrips = append(allTrips, Trip{
				From:      route.PickupLocationName,
				To:        route.ReturnLocationName,
				StartDate: formatDate(trip.AvailableAt),
				EndDate:   formatDate(trip.ExpireTime),
			})
		}
	}

	return allTrips, nil
}

func formatDate(date string) string {
	location, _ := time.LoadLocation(timeZone)

	t, _ := time.Parse(timeLayout, date)
	return t.In(location).Format(timeLayout)
}

// Allow flexible date formats based on time input or not.
func parseFlexTime(s string) (time.Time, error) {
	if len(s) <= 10 {
		return time.Parse("2006-01-02", s)
	}
	return time.Parse(timeLayout, s)
}
