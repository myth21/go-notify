package main

type deliveryUpdateRequest struct {
	UserName     string `json:"userName" binding:"required"`
	UserPhone    int    `json:"userPhone" binding:"required"`
	UserAddress  string `json:"userAddress" binding:"required"`
	UserComment  string `json:"comment"`
	DateTimeFrom string `json:"dateTimeFrom"`
	DateTimeTo   string `json:"dateTimeTo" binding:"required"`
	Status       string `json:"status"`
	UpdatedAt    string `json:"updated_at"`
}
