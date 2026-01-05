package util

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
