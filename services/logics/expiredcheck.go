package logics

import (
	"errors"
	"time"
	"user/domains/entities"
	"user/utils"

	"github.com/sirupsen/logrus"
)

func (service *TransactionService) CheckExpiredStatus(expiredTime time.Time, TransactionToDBValid *entities.Transactions) (string, error) {
	Trans := &entities.Transactions{}
	// Memeriksa apakah waktu kadaluwarsa sudah terlewati dan statusnya masih Tunggu
	if time.Now().After(expiredTime) {
		TransactionToDBSecond, err := service.transactionRepository.GetTransactionByInvoice(TransactionToDBValid.Invoice_Number)
		if err != nil {
			utils.PrintLog("error [services][logics][transaction][CheckExpiredStatus] ", err)
			logrus.Error("error [services][logics][transaction][CheckExpiredStatus] ", err)
			return "", errors.New("data invoice number tidak ditemukan")
		}
		// Jika pengecekan data Status masih Tunggu Maka Update Transaction Invoice Number Menjadi Expired
		if TransactionToDBSecond.Status == "Tunggu" {
			TransactionToDBSecond.Status = "Expired"
			Trans.Status = TransactionToDBSecond.Status
			Trans.Invoice_Number = TransactionToDBSecond.Invoice_Number
			err := service.transactionRepository.UpdateTransactionByInvoiceNumber(Trans)
			if err != nil {
				utils.PrintLog("error [services][logics][transaction][CheckExpiredStatus] ", err)
				logrus.Error("error [services][logics][transaction][CheckExpiredStatus] ", err)
				return "", errors.New("data invoice number tidak ditemukan")
			}
			return "Expired", nil
		}

	}
	return "Tunggu", nil
}
