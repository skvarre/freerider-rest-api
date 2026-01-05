package service

import (
	"strings"
	"time"

	"freerider-rest-api/internal/util"
)

// Filter trips based on provided query params
func FilterTrips(
	trips []util.Trip,
	origins, destinations []string,
	startDate, endDate string,
) ([]util.Trip, error) {

	var filtered []util.Trip

	for _, t := range trips {
		start, _ := time.Parse(util.TimeLayout, t.StartDate)
		end, _ := time.Parse(util.TimeLayout, t.EndDate)

		// Filter trips based on provided query params
		originMatch := matchLocation(origins, t.From)
		destMatch := matchLocation(destinations, t.To)
		dateMatch, err := matchDate(startDate, endDate, start, end)

		if err != nil {
			return nil, err
		}

		if originMatch && destMatch && dateMatch {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

// Check if provided query location matches an existing destination or origin.
// No provided location is seen as a match since it's considered as any location.
func matchLocation(queryLocs []string, loc string) bool {
	if !(len(queryLocs) == 0) {
		for _, d := range queryLocs {
			if strings.Contains(strings.ToLower(loc), strings.ToLower(d)) {
				return true
			}
		}
		return false
	}

	return true
}

// Check if provided query dates match the trip dates.
// No provided dates are seen as a match since it's considered as any date.
func matchDate(queryStartDate, queryEndDate string, start, end time.Time) (bool, error) {
	dateMatch := true

	if queryStartDate != "" {
		sd, err := parseFlexTime(queryStartDate)
		dateMatch = start.After(sd)

		if err != nil {
			return false, err
		}
	}

	if queryEndDate != "" {
		ed, err := parseFlexTime(queryEndDate)
		dateMatch = dateMatch && end.Before(ed)

		if err != nil {
			return false, err
		}
	}

	return dateMatch, nil
}

// Allow flexible date formats based on whether the time was provided as part of the query or not.
func parseFlexTime(s string) (time.Time, error) {
	if len(s) <= 10 {
		return time.Parse("2006-01-02", s)
	}
	return time.Parse(util.TimeLayout, s)
}
