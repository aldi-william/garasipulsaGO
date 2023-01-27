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
			Buttons  []Button  `json:"buttons"`
			Button   string    `json:"button"`
			Sections []Section `json:"sections"`
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

type Section struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

type Row struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type APIReceiver struct {
	Code              int               `json:"code"`
	Receiver_Messages []Receive_Message `json:"messages"`
}

type Receive_Message struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	To             string `json:"to"`
	Receipent_Type string `json:"receipent_type"`
}

type SenderToWa struct {
	Recipient_Type string `json:"recipient_type"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Text           struct {
		Body string `json:"body"`
	} `json:"text"`
}
