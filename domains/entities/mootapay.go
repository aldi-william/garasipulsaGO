package entities

import "time"

type MooTaPayment struct {
	Invoice_Number    string
	Payment_Method_ID string
	Amount            int
	Start_Unique_Code int
	End_Unique_Code   int
	Expired_Date      time.Time
	CallBack_URL      string
	Items
}

type Items struct {
	Name      string
	QTY       int
	Price     int
	SKU       string
	Image_URL string
}
