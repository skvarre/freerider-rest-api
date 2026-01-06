package service

import (
	"freerider-rest-api/pkg/client"
	"freerider-rest-api/pkg/util"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WatchTrips holds logic for POST /watch endpoint
func WatchTrips(ctx *gin.Context) {
	var watcher util.Watcher
	if err := ctx.ShouldBindJSON(&watcher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set headers for SSE
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	// Create a channel for this specific request
	rideChan := make(chan util.Trip)

	// Start the background checker
	go startWatcher(ctx, watcher, rideChan)

	// Stream the results to the caller
	ctx.Stream(func(w io.Writer) bool {
		if trip, ok := <-rideChan; ok {
			ctx.SSEvent("ride-found", trip)
			return true
		}
		return false
	})
}

func startWatcher(ctx *gin.Context, watcher util.Watcher, rideChan chan util.Trip) {
	log.Println("Starting background worker")
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	seenRides := make(map[int]bool)

	// Actual search logic
	performSearch := func() {
		log.Println("Performing search")
		allTrips, _ := client.FetchTrips()
		minDateStr := ""
		maxDateStr := ""

		if !watcher.MinDate.IsZero() {
			minDateStr = watcher.MinDate.Format(util.TimeLayout)
		}

		if !watcher.MaxDate.IsZero() {
			maxDateStr = watcher.MaxDate.Format(util.TimeLayout)
		}

		filtered, _ := FilterTrips(
			allTrips,
			[]string{watcher.Origin},
			[]string{watcher.Destination},
			minDateStr,
			maxDateStr,
		)

		for _, trip := range filtered {
			if !seenRides[trip.RideID] {
				rideChan <- trip
				seenRides[trip.RideID] = true
			}
		}
	}

	performSearch()

	for {
		select {
		case <-ctx.Request.Context().Done():
			return // Client disconnected
		case <-ticker.C:
			performSearch()
		}
	}
}
