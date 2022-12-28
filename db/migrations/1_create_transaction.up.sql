CREATE TABLE `transactions` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `jenis_layanan` varchar(50) NOT NULL,
  `provider` varchar(255) NOT NULL,
  `nomor_hp` varchar(20) NOT NULL,
  `nominal` bigint NOT NULL,
  `pembayaran` varchar(30) NOT NULL,
  `status` varchar(50) NOT NULL,
  `invoice_number` varchar(50) NOT NULL,
  `buyer_sku_code` varchar(150) NOT NULL,
  `id_pelanggan` varchar(150),
  `total` bigint NOT NULL,
  `kode_unik` tinyint NOT NULL,
  `status_pengisian` varchar(50),
  `created_at` timestamp,
  `updated_at` timestamp
);
