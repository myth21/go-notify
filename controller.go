package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

// getDeliveries responds with the list of all deliveries as JSON.
func getDeliveries(context *gin.Context) {

	db, err := sql.Open(ConfigDbDriver, ConfigDataSourceName)
	if err != nil {
		//panic(err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from delivery")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	//deliveries := []delivery{}
	var deliveries []delivery
	for rows.Next() {
		d := delivery{}
		err := rows.Scan(&d.ID, &d.UserName, &d.UserPhone, &d.UserAddress, &d.UserComment, &d.DateTimeFrom, &d.DateTimeTo, &d.Status, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		deliveries = append(deliveries, d)
	}

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

	db, err := sql.Open(ConfigDbDriver, ConfigDataSourceName)
	if err != nil {
		//panic(err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	sqlString := "insert into delivery (user_name, user_phone, user_address, user_comment, date_time_from, date_time_to) values ($1, $2, $3, $4, $5, $6)"
	result, err := db.Exec(sqlString, newModel.UserName, newModel.UserPhone, newModel.UserAddress, newModel.UserComment, newModel.DateTimeFrom, newModel.DateTimeTo)
	if err != nil {
		//panic(err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lastInsertId, LastInsertIdErr := result.LastInsertId()
	if LastInsertIdErr != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": LastInsertIdErr.Error()})
		return
	}

	rowsAffected, rowsAffectedErr := result.RowsAffected()
	if rowsAffectedErr != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": rowsAffectedErr.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"Count of rows added": rowsAffected,
		"Last insert id":      lastInsertId,
		"Saved model":         newModel,
	})

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
