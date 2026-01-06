package mobile

// Module is intended for using the API in gomobile

import (
	"encoding/json"
	"freerider-rest-api/pkg/client"
	"freerider-rest-api/pkg/service"
	"strings"
	"time"
)

type WatcherCallback interface {
	OnRideFound(tripJSON string)
	OnError(message string)
}

// Replaces trips.getTrips
func GetFilteredTrips(origin, destination, startDate, endDate string) (string, error) {
	trips, err := client.FetchTrips()
	if err != nil {
		return "", err
	}

	// TODO: split origin and destination
	origins := strings.Split(origin, "|")
	destinations := strings.Split(destination, "|")

	filtered, err := service.FilterTrips(
		trips,
		origins,
		destinations,
		startDate,
		endDate,
	)
	if err != nil {
		return "", err
	}

	result, _ := json.Marshal(filtered)
	return string(result), nil
}

// Replaces watch.startWatcher
func StartWatch(origin, destination, minDate, maxDate string, cb WatcherCallback) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		seenRides := make(map[int]bool)

		for {
			allTrips, err := client.FetchTrips()
			if err != nil {
				cb.OnError(err.Error())
				return
			}

			origins := strings.Split(origin, "|")
			destinations := strings.Split(destination, "|")

			filtered, err := service.FilterTrips(
				allTrips,
				origins,
				destinations,
				minDate,
				maxDate,
			)

			if err != nil {
				cb.OnError(err.Error())
				return
			}

			for _, trip := range filtered {
				if !seenRides[trip.RideID] {
					data, _ := json.Marshal(trip)
					cb.OnRideFound(string(data))
					seenRides[trip.RideID] = true
				}
			}

			<-ticker.C // Wait for next tick
		}
	}()
}
