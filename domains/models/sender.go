package models

type Sender struct {
	Chat         string `json:"chat"`
	Sender_Phone string `json:"sender_phone"`
	Message_Type string `json:"message_type"`
	Message_Text string `json:"message_text"`
}

type APISenderWithButton struct {
	Recipient_type string `json:"recipient_type"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Interactive    struct {
		Type   string `json:"type"`
		Header struct {
			Text string `json:"text"`
		} `json:"header"`
		Body struct {
			Text string `json:"text"`
		} `json:"body"`
		Footer struct {
			Text string `json:"text"`
		} `json:"footer"`
		Action struct {
			Buttons []Button `json:"buttons"`
		}
	}
}

type Button struct {
	Type  string `json:"type"`
	Reply struct {
		ID    string `json:"id"`
		Title string `json:"Title"`
	} `json:"reply"`
}
