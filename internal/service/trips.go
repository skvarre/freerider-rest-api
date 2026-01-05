package service

import (
	"freerider-rest-api/internal/client"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetTrips holds logic for GET /trips endpoint
func GetTrips(ctx *gin.Context) {
	if allowed, param := allowedParams(ctx); !allowed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameter: " + param})
		return
	}

	origin := strings.ToLower(ctx.Query("origin"))
	dest := strings.ToLower(ctx.Query("destination"))
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	trips, err := client.FetchTrips()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filtered, err := FilterTrips(trips, origin, dest, startDate, endDate)
	if err != nil {
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
