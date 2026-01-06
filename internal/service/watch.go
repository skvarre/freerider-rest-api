package service

import (
	"fmt"
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
	ticker := time.NewTicker(10 * time.Minute)
	seenRides := make(map[int]bool)

	for range ticker.C {
		log.Println("Checking for watched rides")
		allTrips, err := client.FetchTrips()
		if err != nil {
			log.Println("Error fetching trips: ", err)
			continue
		}

		for _, watcher := range watchList {
			filtered, err := FilterTrips(
				allTrips,
				[]string{watcher.Origin},
				[]string{watcher.Destination},
				watcher.MinDate.Format(util.TimeLayout),
				watcher.MaxDate.Format(util.TimeLayout),
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
}

// TODO: Send response
func sendNotification(t util.Trip) {
	fmt.Println("Found ride", t.From, "to", t.To)
}
