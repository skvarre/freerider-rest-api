package client

import (
	"encoding/json"
	"net/http"
	"time"

	"freerider-rest-api/internal/util"
)

type Trip = util.Trip
type FreeriderRoute = util.FreeriderRoute

// FetchTrips Fetches all trips from Freerider API
func FetchTrips() ([]Trip, error) {

	resp, err := http.Get(util.FreeriderURL)
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
				RideID:        trip.ID,
				From:          route.PickupLocationName,
				To:            route.ReturnLocationName,
				AvailableFrom: formatDate(trip.AvailableAt),
				Expires:       formatDate(trip.ExpireTime),
			})
		}
	}

	return allTrips, nil
}

func formatDate(date string) string {
	location, _ := time.LoadLocation(util.TimeZone)

	t, _ := time.Parse(util.TimeLayout, date)
	return t.In(location).Format(util.TimeLayout)
}
