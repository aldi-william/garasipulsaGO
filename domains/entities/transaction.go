package entities

import "time"

type Transactions struct {
	ID               uint       `json:"id"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	Pembayaran       string     `json:"pembayaran"`
	JenisLayanan     string     `json:"jenis_layanan"`
	Provider         string     `json:"provider"`
	Nominal          int        `json:"nominal"`
	Nomor_Hp         int        `json:"nomor_hp"`
	Status           string     `json:"status"`
	Invoice_Number   string     `json:"invoice_number"`
	Buyer_Sku_Code   string     `json:"buyer_sku_code"`
	Total            int        `json:"total"`
	Status_Pengisian string     `json:"status_pengisian"`
	Kode_Unik        int        `json:"kode_unik"`
	Id_Pelanggan     string     `json:"id_pelanggan"`
}
