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
		start, _ := time.Parse(util.TimeLayout, t.AvailableFrom)
		end, _ := time.Parse(util.TimeLayout, t.Expires)

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

// Check if provided query date window matches the trip dates.
// No provided dates are seen as a match since it's considered as any date.
func matchDate(queryStartDate, queryEndDate string, start, end time.Time) (bool, error) {
	// Default to extreme values
	searchStart := time.Time{}                 // Beginning of time
	searchEnd := time.Now().AddDate(100, 0, 0) // Far future

	// Parse query dates
	if queryStartDate != "" {
		var err error
		searchStart, err = parseFlexTime(queryStartDate)
		if err != nil {
			return false, err
		}
	}

	if queryEndDate != "" {
		var err error
		searchEnd, err = parseFlexTime(queryEndDate)
		if err != nil {
			return false, err
		}
	}

	isAvailableBeforeWindowEnds := start.Before(searchEnd) || start.Equal(searchEnd)
	isNotExpiredBeforeWindowStarts := end.After(searchStart) || end.Equal(searchStart)

	return isAvailableBeforeWindowEnds && isNotExpiredBeforeWindowStarts, nil
}

// Allow flexible date formats based on whether the time was provided as part of the query or not.
func parseFlexTime(s string) (time.Time, error) {
	if len(s) <= 10 {
		return time.Parse("2006-01-02", s)
	}
	return time.Parse(util.TimeLayout, s)
}
