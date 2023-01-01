package models

import "time"

type Transaction struct {
	ID               int        `json:"id"`
	Pembayaran       string     `json:"pembayaran" binding:"required"`
	JenisLayanan     string     `json:"jenis_layanan" binding:"required"`
	Provider         string     `json:"provider" binding:"required"`
	Nominal          int        `json:"nominal" binding:"required"`
	Nomor_Hp         string     `json:"nomor_hp" binding:"required"`
	Status           string     `json:"status"`
	Buyer_Sku_Code   string     `json:"buyer_sku_code" binding:"required"`
	Invoice_Number   string     `json:"invoice_number"`
	Kode_Unik        int        `json:"kode_unik"`
	Total            int        `json:"total"`
	Status_Pengisian string     `json:"status_pengisian"`
	CreatedAt        *time.Time `json:"created_at"`
	ID_Pelanggan     string     `json:"id_pelanggan"`
	Serial_Number    string     `json:"sn"`
}
