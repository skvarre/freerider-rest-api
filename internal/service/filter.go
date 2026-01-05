package service

import (
	"strings"
	"time"

	"freerider-rest-api/internal/util"
)

func FilterTrips(
	trips []util.Trip,
	origins, destinations []string,
	startDate, endDate string,
) ([]util.Trip, error) {

	var filtered []util.Trip

	for _, t := range trips {
		start, _ := time.Parse(util.TimeLayout, t.StartDate)
		end, _ := time.Parse(util.TimeLayout, t.EndDate)

		matchOrigin := matchLocation(origins, t.From)
		matchDest := matchLocation(destinations, t.To)

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

// Check if provided query location matches an existing destination or origin.
// No provided location is seen as a match since it's considered as any location.
func matchLocation(loc []string, queryLoc string) bool {
	if !(len(loc) == 0) {
		for _, d := range loc {
			if strings.Contains(strings.ToLower(queryLoc), d) {
				return true
			}
		}
		return false
	}

	return true
}

// Allow flexible date formats based on time input or not.
func parseFlexTime(s string) (time.Time, error) {
	if len(s) <= 10 {
		return time.Parse("2006-01-02", s)
	}
	return time.Parse(util.TimeLayout, s)
}
