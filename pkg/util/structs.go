package util

import "time"

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

type FreeriderLocation struct {
	// We only care about the name for now
	Name string `json:"name"`
}

type Trip struct {
	RideID        int    `json:"rideId"`
	From          string `json:"from"`
	To            string `json:"to"`
	AvailableFrom string `json:"availableFrom"`
	Expires       string `json:"expires"`
}
type Watcher struct {
	ID          string    `json:"id"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	MinDate     time.Time `json:"min_date"`
	MaxDate     time.Time `json:"max_date"`
}
