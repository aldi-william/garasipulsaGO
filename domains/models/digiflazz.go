package models

type TransactionDIGIFLAZZ struct {
	Command        string `json:"command"`
	Username       string `json:"username"`
	Buyer_SKU_Code string `json:"buyer_sku_code"`
	Customer_NO    string `json:"customer_no"`
	Ref_ID         string `json:"ref_id"`
	Sign           string `json:"sign"`
	Testing        bool   `json:"testing"`
}

type ResultDigiFlazzData struct {
	Data ResultDigiFlazz `json:"data"`
}

type ResultDigiFlazz struct {
	Ref_ID           string      `json:"ref_id"`
	Customer_No      string      `json:"customer_no"`
	Customer_Name    string      `json:"customer_name"`
	Buyer_SKU_Code   string      `json:"buyer_sku_code"`
	Admin            int         `json:"admin"`
	Message          string      `json:"message"`
	Status           string      `json:"status"`
	Response_Code    string      `json:"rc"`
	Serial_Number    string      `json:"sn"`
	Buyer_Last_Saldo int         `json:"buyer_last_saldo"`
	Price            int         `json:"price"`
	Selling_Price    int         `json:"selling_price"`
	Desc             interface{} `json:"desc"`
}
