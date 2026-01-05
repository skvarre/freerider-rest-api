package service

import (
	"freerider-rest-api/internal/client"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetTrips holds logic for GET /trips endpoint
func GetTrips(ctx *gin.Context) {
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
