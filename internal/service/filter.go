package service

import (
	"strings"
	"time"

	"freerider-rest-api/internal/util"
)

func FilterTrips(
	trips []util.Trip,
	origin, destination string,
	startDate, endDate string,
) ([]util.Trip, error) {

	var filtered []util.Trip

	for _, t := range trips {
		start, _ := time.Parse(util.TimeLayout, t.StartDate)
		end, _ := time.Parse(util.TimeLayout, t.EndDate)

		matchOrigin := origin == "" || strings.Contains(strings.ToLower(t.From), origin)
		matchDest := destination == "" || strings.Contains(strings.ToLower(t.To), destination)
		matchDate := true

		if startDate != "" {
			sd, _ := parseFlexTime(startDate)
			matchDate = start.After(sd)
		}

		if endDate != "" {
			ed, _ := parseFlexTime(endDate)
			matchDate = matchDate && end.Before(ed)
		}

		if matchOrigin && matchDest && matchDate {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

// Allow flexible date formats based on time input or not.
func parseFlexTime(s string) (time.Time, error) {
	if len(s) <= 10 {
		return time.Parse("2006-01-02", s)
	}
	return time.Parse(util.TimeLayout, s)
}
