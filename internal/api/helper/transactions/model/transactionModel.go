package model

type TransactionHistory struct {
	RefTransHisId   int    `gorm:"column:refTransHisId;primaryKey;autoIncrement"`
	RefTransTypeId  int    `gorm:"column:refTransTypeId"`
	RefTransHisData string `gorm:"column:refTransHisData"`
	CreatedAt       string `gorm:"column:createdAt"`
	CreatedBy       string `gorm:"column:createdBy"`
	RefUserId       int    `gorm:"column:refUserId"`
}
