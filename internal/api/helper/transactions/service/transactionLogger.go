package service

import (
	"time"

	"github.com/ZADPRO/Snehalaya-Backend-GoLang/internal/api/helper/transactions/model"
	"gorm.io/gorm"
)

func LogTransaction(db *gorm.DB, userId int, createdBy string, transTypeId int, message string) error {
	history := model.TransactionHistory{
		RefTransTypeId:  transTypeId,
		RefTransHisData: message,
		CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
		CreatedBy:       createdBy,
		RefUserId:       userId,
	}
	return db.Table(`"TransactionHistory"`).Create(&history).Error
}
