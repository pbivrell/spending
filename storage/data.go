package storage

import "time"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Page struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}

type Estimate struct {
	BaseTransaction
	Occurance *int       `json:"occurance"`
	Period    *int       `json:"period"`
	Date      *time.Time `json:"date"`
}

type BaseTransaction struct {
	ID      int    `json:"id"`
	PageID  int    `json:"pageID"`
	Amount  int    `json:"amount"`
	Decimal int    `json:"decimal"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

type Transaction struct {
	BaseTransaction
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

type Role struct {
	ID         int `json:"id"`
	UserID     int `json:"userID"`
	Permission int `json:"permission"`
	PageID     int `json:"pageID"`
}
