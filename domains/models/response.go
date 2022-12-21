package models

type ResponseOnboarding struct {
	UserID int `json:"user_id"`
}

type BodyResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
