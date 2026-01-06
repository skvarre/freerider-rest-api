package service

import (
	"errors"
	"freerider-rest-api/pkg/client"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTrips holds logic for GET /trips endpoint
func GetTrips(ctx *gin.Context) {
	if allowed, param := allowedParams(ctx); !allowed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameter: " + param})
		return
	}

	destinations := ctx.QueryArray("destination")
	origins := ctx.QueryArray("origin")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	for i, d := range destinations {
		destinations[i] = strings.ToLower(d)
	}

	for i, d := range origins {
		origins[i] = strings.ToLower(d)
	}

	// Fetch trips from Freerider API
	trips, err := client.FetchTrips()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filtered, err := FilterTrips(trips, origins, destinations, startDate, endDate)

	if err != nil {
		var parseErr *time.ParseError

		if errors.As(err, &parseErr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format: " + err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, filtered)
}

func allowedParams(ctx *gin.Context) (bool, string) {
	if len(ctx.Request.URL.Query()) == 0 {
		// We allow no query params
		return true, ""
	}

	allowedParams := map[string]bool{
		"origin":      true,
		"destination": true,
		"startDate":   true,
		"endDate":     true,
	}

	for key := range ctx.Request.URL.Query() {
		if !allowedParams[key] {
			return false, key
		}
	}

	return true, ""
}
