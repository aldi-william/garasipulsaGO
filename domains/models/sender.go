package models

type Sender struct {
	Chat         string `json:"chat"`
	Sender_Phone string `json:"sender_phone"`
	Message_Type string `json:"message_type"`
	Message_Text string `json:"message_text"`
}
