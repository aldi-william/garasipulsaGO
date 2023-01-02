package models

type MootaCallback struct {
	Amount           string `json:"amount"`
	Type_Transaction string `json:"type"`
	Mutation_ID      string `json:"mutation_id"`
	Token            string `json:"token"`
	Balance          int    `json:"balance"`
	Bank_ID          string `json:"bank_id"`
	Description      string `json:"description"`
	Account_Number   int    `json:"account_number"`
}

type ResultToBuyer struct {
	Invoice_Number string `json:"invoice_number"`
	Unique_Code    int    `json:"unique_code"`
	Total          int    `json:"total"`
}
