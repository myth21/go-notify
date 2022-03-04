package main

type delivery struct {
	ID           string `json:"id"`
	UserName     string `json:"userName" binding:"required"`
	UserPhone    string `json:"userPhone" binding:"required"`
	UserAddress  string `json:"userAddress" binding:"required"`
	UserComment  string `json:"comment"`
	DateTimeFrom string `json:"dateTimeFrom"`
	DateTimeTo   string `json:"dateTimeTo" binding:"required"`
	Status       string `json:"status"`
	//DateTimeTo   time.Time `json:"dateTimeTo" binding:"required,ableDateTimeTo" time_format:"2022-02-23"`
}

//var ableDateTimeTo validator.Func = func(fl validator.FieldLevel) bool {
//	date, ok := fl.Field().Interface().(time.Time)
//	if ok {
//		today := time.Now()
//		if today.After(date) {
//			return false
//		}
//	}
//	return true
//}

// deliveries slice to seed record delivery data.
var deliveries = []delivery{
	{ID: "1", UserName: "Ivan", UserAddress: "John Street, 34", DateTimeFrom: "2022-02-23", DateTimeTo: "2022-02-24 12:12:12", UserComment: "A comment..."},
	{ID: "2", UserName: "Rita", UserAddress: "Coltrane Street, 22", DateTimeFrom: "2022-02-23", DateTimeTo: "2022-02-24 12:12:12", UserComment: "A long comment..."},
}
