package service

import (
	"freerider-rest-api/internal/util"
	"testing"
)

func TestFilterTrips(t *testing.T) {
	trips := []util.Trip{
		{From: "Bollnäs & Self Service Kiosk", To: "Mora DT", StartDate: "2024-01-01T10:00:00", EndDate: "2024-01-01T16:00:00"},
		{From: "Stockholm Arlanda Flygplats / Self Service Kiosk", To: "Karlskrona Ahlberg Bil / Self Service Kiosk", StartDate: "2024-02-01T10:00:00", EndDate: "2024-02-01T16:00:00"},
	}

	tests := []struct {
		name         string
		origins      []string
		destinations []string
		startDate    string
		endDate      string
		expectedLen  int
	}{
		{
			name:        "Match by origin",
			origins:     []string{"Stockholm"},
			expectedLen: 1,
		},
		{
			name:         "Match by destination",
			destinations: []string{"Karlskrona"},
			expectedLen:  1,
		},
		{
			name:        "No match",
			origins:     []string{"Borlänge"},
			expectedLen: 0,
		},
		{
			name:        "Empty filters returns all",
			expectedLen: 2,
		},
		{
			name:        "Filter by date range",
			startDate:   "2023-12-31",
			endDate:     "2024-01-02",
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterTrips(trips, tt.origins, tt.destinations, tt.startDate, tt.endDate)
			if err != nil {
				t.Fatalf("FilterTrips failed: %v", err)
			}

			if len(got) != tt.expectedLen {
				t.Errorf("expected %d trips, got %d", tt.expectedLen, len(got))
			}
		})
	}
}
