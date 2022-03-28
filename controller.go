package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	//"strings"

	"reflect"

	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/mattn/go-sqlite3"
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

	db, err := sql.Open(ConfigDbDriver, ConfigDataSourceName)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	//row, err := db.QueryRow("select * from delivery where id = ?", id)
	row := db.QueryRow("select * from delivery where id = ?", id)

	// if err != nil {
	// 	context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// defer row.Close()

	d := delivery{}
	err2 := row.Scan(&d.ID, &d.UserName, &d.UserPhone, &d.UserAddress, &d.UserComment, &d.DateTimeFrom, &d.DateTimeTo, &d.Status, &d.CreatedAt, &d.UpdatedAt)
	if err2 != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, d)

	/*
		rows, err := db.Query("select * from delivery where id = " + id)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // "message": "Not Found"
			return
		}
		defer rows.Close()


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
	*/

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

	context.IndentedJSON(http.StatusCreated, gin.H{
		"Count of rows added": rowsAffected,
		"Last insert id":      lastInsertId,
		"Saved model":         newModel,
	})

}

/*
func getField(v *delivery, field string) bool {
	metaValue := reflect.ValueOf(v).Elem()

	//r := reflect.ValueOf(v)
	//f := reflect.Indirect(r).FieldByName(field)

	return bool(metaValue.FieldByName(field))
}
*/

type deliveryUpdateRequestMap map[string]interface{}

// https://ru.stackoverflow.com/questions/773326/%D0%9F%D0%BE%D0%BB%D1%83%D1%87%D0%B5%D0%BD%D0%B8%D0%B5-%D1%82%D0%B5%D0%BB%D0%B0-%D0%B7%D0%B0%D0%BF%D1%80%D0%BE%D1%81%D0%B0-%D0%B8-%D0%BE%D0%B1%D0%BD%D0%BE%D0%B2%D0%BB%D0%B5%D0%BD%D0%B8%D0%B5-%D1%81%D0%BE%D0%BE%D1%82%D0%B2%D0%B5%D1%82%D1%81%D1%82%D0%B2%D1%83%D1%8E%D1%89%D0%B5%D0%B3%D0%BE-%D0%BF%D0%BE%D0%BB%D1%8F-%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D1%8B-golang
func updateDeliveryById(context *gin.Context) {

	// Vaidate request
	var deliveryUpdateRequestValidator deliveryUpdateRequest
	//context.ShouldBindJSON(&deliveryUpdateRequestValidator)
	requestErr := context.ShouldBindBodyWith(&deliveryUpdateRequestValidator, binding.JSON)
	if requestErr != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": requestErr.Error()})
		return
	}

	// Box for values
	var deliveryMap deliveryUpdateRequestMap
	//jsonDataBytes2, errR := ioutil.ReadAll(context.Request.Body)
	errR := context.ShouldBindBodyWith(&deliveryMap, binding.JSON)
	if errR != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"_message_": errR.Error()})
		return
	}

	var deliveryModel delivery

	element := reflect.ValueOf(&deliveryModel).Elem()

	fieldNames := make(map[string]interface{}) // read wht this does exactlly?

	for i := 0; i < element.NumField(); i++ {
		fieldName := element.Type().Field(i).Name
		//varType := element.Type().Field(i).Type
		//varValue := element.Field(i).Interface()

		//fmt.Printf("%v %v %v\n", varName, varType, varValue)

		_, ok := deliveryMap[fieldName]
		if ok {

			fieldNames[fieldName] = deliveryMap[fieldName]
			// todo sql string
			//context.IndentedJSON(http.StatusOK, gin.H{"found": fieldName})
			// fieldNames[fieldName] = deliveryMap[fieldName]

		} else {
			//context.IndentedJSON(http.StatusOK, gin.H{"not found": fieldName})
		}

	}

	/*
			id := context.Param("id")
		fieldsToUpdate := implode(",", context.ParamsExcludingId())
		res, err := db.Exec("update delivery set " + fieldsToUpdate + " where id = ?", id)

		// $sql = 'UPDATE `'.static::tableName().'` SET '.$this->getUpdatingAvailableValues().' WHERE `'.static::$primaryKeyName.'`="'.$this->getPrimaryKey().'"';
		res, err := db.Exec("update delivery set user_address = 'zxczczczxcxc' where id = ?", id)
		if err != nil {
			//panic(err)
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": rowsAffectedErr.Error()})
			return
		}

		rowsAffectedStr := strconv.Itoa(int(rowsAffected))

		context.IndentedJSON(http.StatusOK, gin.H{"message": "Updated " + rowsAffectedStr + " row"})
	*/

	//todo
	//strings.Join(fieldNames, ",")

	context.IndentedJSON(http.StatusOK, gin.H{"message": fieldNames})

	return

	//var jsonRaw = `{"ip":"0.0.0.0","tag":["something tag"],"apps":["foo","bar"],"active":true,"params":{"key1":"val1","key2":"val2"}} `
	jsonDataBytes, errReadAll := ioutil.ReadAll(context.Request.Body)
	if errReadAll != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message2": errReadAll.Error()})
		return
	}

	// Get json fields from request body
	//var deliveryMap deliveryUpdateRequestMap

	//var modelUpdateRequest deliveryUpdateRequest
	//err := json.Unmarshal([]byte(jsonRaw), &model)
	err := json.Unmarshal(jsonDataBytes, &deliveryMap)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message3": err.Error()})
		return
	}
	//model["ip"] = "1.2.3.4"
	context.IndentedJSON(http.StatusOK, gin.H{"message": deliveryMap})
	return

	id := context.Param("id")

	db, err := sql.Open(ConfigDbDriver, ConfigDataSourceName)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer db.Close()

	var exist bool
	row := db.QueryRow("select exists (select * from delivery where id = ?) as `exist`", id)
	err2 := row.Scan(&exist)
	if err2 != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error2": err2.Error()})
		return
	}

	if exist == true {

		jsonDataBytes, errRead := ioutil.ReadAll(context.Request.Body)
		if errRead != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": errRead.Error()})
			return
		}

		// xType := fmt.Sprintf("%T", jsonDataBytes)
		// context.JSON(http.StatusOK, xType)

		// return

		//var var_CreateExamType = &CreateExamType{}
		//json_parse(string(jsonDataBytes), var_CreateExamType)
		//log.Println(var_CreateExamType)

		var obj interface{}
		err := json.Unmarshal(jsonDataBytes, obj)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		context.IndentedJSON(http.StatusOK, obj)
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{"message": "Not found"})
}

func json_parse(Data string, obj interface{}) {

	var b_Data = []byte(Data)
	err := json.Unmarshal(b_Data, obj)
	if err != nil {
		panic(err)
		//log.Println("error:", err)
	}

	//return obj
}

func deleteDeliveryById(context *gin.Context) {

	id := context.Param("id")

	// _, err := database.Exec("delete from productdb.Products where id = ?", id)

	db, err := sql.Open(ConfigDbDriver, ConfigDataSourceName)
	if err != nil {
		//panic(err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	res, err := db.Exec("delete from delivery where id = ?", id)
	if err != nil {
		//panic(err)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, rowsAffectedErr := res.RowsAffected()
	if rowsAffectedErr != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": rowsAffectedErr.Error()})
		return
	}

	rowsAffectedStr := strconv.Itoa(int(rowsAffected))

	context.IndentedJSON(http.StatusOK, gin.H{"message": "Deleted " + rowsAffectedStr + " row"})

}
