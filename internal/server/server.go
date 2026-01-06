package server

import (
	"freerider-rest-api/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	r.GET("/trips", service.GetTrips)
	r.POST("/watch", service.WatchTrips)
	//r.GET("/locations")

	log.Println("Server started on port 8080")
	r.Run(":8080")
}
