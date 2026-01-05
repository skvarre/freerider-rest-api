package server

import (
	"freerider-rest-api/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	r.GET("/trips", service.GetTrips)
	//r.GET("/locations")
	//r.POST("/watchlist")

	r.Run(":8080")
	log.Println("Server started on port 8080")
}
