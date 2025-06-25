package model

type TransactionHistory struct {
	RefTransHisId   int    `gorm:"column:refTransHisId;primaryKey;autoIncrement" json:"refTransHisId"`
	RefTransTypeId  int    `gorm:"column:refTransTypeId" json:"refTransTypeId"`
	RefTransHisData string `gorm:"column:refTransHisData" json:"refTransHisData"`
	CreatedAt       string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy       string `gorm:"column:createdBy" json:"createdBy"`
	RefUserId       int    `gorm:"column:refUserId" json:"refUserId"`
}
