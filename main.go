package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
)

func main() {

	router := gin.Default()
	router.GET("/deliveries", getDeliveries)
	router.GET("/deliveries/:id", getDeliveryById)
	router.POST("/deliveries", postDelivery)
	router.DELETE("/deliveries/:id", deleteDeliveryById)

	router.Run(ConfigHostPort)
}

//// delivery represents data about a record delivery.
//type delivery struct {
//	ID               string  `json:"id"`
//	UserName           string  `json:"userName"`
//	Address          string  `json:"address"`
//	DateTimeFrom     string  `json:"dateTimeFrom"`
//	DateTimeTo       string  `json:"dateTimeTo"`
//	Comment          string  `json:"comment"`
//}

//// deliveries slice to seed record delivery data.
//var deliveries = []delivery{
//	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
//	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
//	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
//}

//// getDeliveries responds with the list of all deliveries as JSON.
//func getDeliveries(context *gin.Context) {
//	context.IndentedJSON(http.StatusOK, deliveries)
//}
//
//func getDeliveryById(context *gin.Context) {
//	id := context.Param("id")
//
//	for _, album := range deliveries {
//		if id == album.ID {
//			context.IndentedJSON(http.StatusOK, album)
//			return
//		}
//	}
//
//	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found"})
//}
//
//func postDelivery(context *gin.Context) {
//
//	var newAlbum delivery
//
//	//err := context.BindJSON(&newAlbum) // this write in header Content-Type: text/plain!
//	err := context.ShouldBindJSON(&newAlbum)
//
//	if err != nil {
//		context.Errors.Errors()
//		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Assign error?"})
//		return
//	}
//
//	context.IndentedJSON(http.StatusOK, newAlbum)
//}
