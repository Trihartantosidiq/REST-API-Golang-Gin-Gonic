package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	Email        string `json:"email"`
}

type Product struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	UserID int    `json:"user_id"`
}
