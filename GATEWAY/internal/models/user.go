package models

type RegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}
