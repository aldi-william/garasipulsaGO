package models

type MootaCallback struct {
	Amount           int    `json:"amount"`
	Type_Transaction string `json:"type"`
	Mutation_ID      string `json:"mutation_id"`
	Token            string `json:"token"`
}

type ResultToBuyer struct {
	Invoice_Number string `json:"invoice_number"`
	Unique_Code    int    `json:"unique_code"`
	Total          int    `json:"total"`
}
