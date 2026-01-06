package service

import (
	"freerider-rest-api/internal/client"
	"freerider-rest-api/internal/util"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TODO: Use some kind of sqlite db or similar
var watchList []util.Watcher

func WatchTrips(ctx *gin.Context) {
	var newWatch util.Watcher
	if err := ctx.ShouldBindJSON(&newWatch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Handle ID better
	newWatch.ID = time.Now().Format("150405")
	watchList = append(watchList, newWatch)
	ctx.JSON(http.StatusCreated, newWatch)

	// Keep watch
	go runBackgroundWorker()
}

func runBackgroundWorker() {
	// TODO: Every 10 min for now, but should be irregular
	log.Println("Starting background worker")
	ticker := time.NewTicker(10 * time.Minute)
	seenRides := make(map[int]bool)
	defer ticker.Stop()

	search(seenRides)

	for range ticker.C {
		search(seenRides)
	}
}

func search(seenRides map[int]bool) {
	log.Println("Checking for watched rides")
	allTrips, err := client.FetchTrips()
	if err != nil {
		log.Println("Error fetching trips: ", err)
		return
	}

	for _, watcher := range watchList {
		minDateStr := ""
		maxDateStr := ""

		if !watcher.MinDate.IsZero() {
			minDateStr = watcher.MinDate.Format(util.TimeLayout)
		}

		if !watcher.MaxDate.IsZero() {
			maxDateStr = watcher.MaxDate.Format(util.TimeLayout)
		}

		filtered, err := FilterTrips(
			allTrips,
			[]string{watcher.Origin},
			[]string{watcher.Destination},
			minDateStr,
			maxDateStr,
		)

		if err != nil {
			log.Println("Error filtering trips: ", err)
			continue
		}

		for _, trip := range filtered {
			if !seenRides[trip.RideID] {
				sendNotification(trip)
				seenRides[trip.RideID] = true
			}
		}
	}
}

// TODO: Send response
func sendNotification(t util.Trip) {
	log.Println("Found ride!", t.From, "to", t.To)
}
