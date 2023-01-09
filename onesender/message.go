package onesender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ApiUrl string
	ApiKey string
)

type TextBody struct {
	Body string `json:"body"`
}

type ImageBody struct {
	Url     string `json:"link"`
	Caption string `json:"caption"`
}

type DocumentBody struct {
	Url string `json:"link"`
}

type MessageText struct {
	To            string   `json:"to"`
	RecipientType string   `json:"recipient_type"`
	Type          string   `json:"type"`
	Text          TextBody `json:"text"`
}

type MessageImage struct {
	To            string    `json:"to"`
	RecipientType string    `json:"recipient_type"`
	Type          string    `json:"type"`
	Image         ImageBody `json:"image"`
}

type MessageDocument struct {
	To            string       `json:"to"`
	RecipientType string       `json:"recipient_type"`
	Type          string       `json:"type"`
	Document      DocumentBody `json:"document"`
}

func (m *MessageText) SetMessage(message string) {
	tipe := "individual"
	if strings.Contains(m.To, "@g.us") {
		tipe = "group"
	}

	m.Type = "text"
	m.RecipientType = tipe
	m.Text = TextBody{
		Body: message,
	}
}

func (m *MessageImage) SetMessage(link, message string) {
	tipe := "individual"
	if strings.Contains(m.To, "@g.us") {
		tipe = "group"
	}

	m.Type = "image"
	m.RecipientType = tipe
	m.Image = ImageBody{
		Url:     link,
		Caption: message,
	}
}

func (m *MessageDocument) SetMessage(link string) {
	tipe := "individual"
	if strings.Contains(m.To, "@g.us") {
		tipe = "group"
	}

	m.Type = "document"
	m.RecipientType = tipe
	m.Document = DocumentBody{
		Url: link,
	}
}

func Transform(msg interface{}) ([]byte, error) {
	var output []byte
	var err error
	output, err = json.Marshal(msg)
	if err != nil {
		return output, err
	}

	return output, nil
}

func SendMessage(msg interface{}) (string, error) {

	msgByte, err := Transform(msg)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ApiUrl, bytes.NewBuffer(msgByte))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "Error post data", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func SendTextMessage(to, message string) (string, error) {
	textMessage := MessageText{
		To: to,
	}
	textMessage.SetMessage(message)

	res, err := SendMessage(textMessage)

	if err != nil {
		return "", err
	}

	return res, nil
}

func SendImageMessage(to, link, caption string) (string, error) {
	imgMessage := MessageImage{
		To: to,
	}
	imgMessage.SetMessage(link, caption)

	res, err := SendMessage(imgMessage)

	if err != nil {
		return "", err
	}

	return res, nil
}

func SendDocumentMessage(to, link string) (string, error) {
	docMessage := MessageDocument{
		To: to,
	}
	docMessage.SetMessage(link)

	res, err := SendMessage(docMessage)

	if err != nil {
		return "", err
	}

	return res, nil
}
