package models

import "time"

type Transaction struct {
	Pembayaran     string     `json:"pembayaran" binding:"required"`
	JenisLayanan   string     `json:"jenis_layanan" binding:"required"`
	Provider       string     `json:"provider" binding:"required"`
	Nominal        int        `json:"nominal" binding:"required"`
	Nomor_Hp       int        `json:"nomor_hp" binding:"required"`
	Status         string     `json:"status"`
	Buyer_Sku_Code string     `json:"buyer_sku_code" binding:"required"`
	CreatedAt      *time.Time `json:"created_at"`
}

type TransactionPLN struct {
	Transaction
	ID_Pelanggan string `json:"id_pelanggan"`
}
