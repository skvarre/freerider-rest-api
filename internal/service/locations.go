package service

import (
	"freerider-rest-api/internal/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLocations(ctx *gin.Context) {
	locations, err := client.FetchLocations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, locations)
}
