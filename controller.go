package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// getDeliveries responds with the list of all deliveries as JSON.
func getDeliveries(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, deliveries)
}

func getDeliveryById(context *gin.Context) {
	id := context.Param("id")

	for _, delivery := range deliveries {
		if id == delivery.ID {
			context.IndentedJSON(http.StatusOK, delivery)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}

func postDelivery(context *gin.Context) {

	// todo validate request please here or early
	// base data you can take from database, rested depended on cases add in request validator

	var newModel delivery

	//err := context.BindJSON(&newModel) // this write in header Content-Type: text/plain!
	err := context.ShouldBindJSON(&newModel)

	if err != nil {
		//context.Errors.Errors()
		//context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Assign error?"})
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, newModel)
}

func deleteDeliveryById(context *gin.Context) {

	id := context.Param("id")

	for _, delivery := range deliveries {
		if id == delivery.ID {
			context.IndentedJSON(http.StatusOK, delivery.ID)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}
