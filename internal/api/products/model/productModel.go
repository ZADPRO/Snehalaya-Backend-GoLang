package productModel

type POProduct struct {
	POId          int    `json:"poId" gorm:"column:poId;primaryKey;autoIncrement"`
	PoName        string `json:"poName" gorm:"column:poName"`
	PoDescription string `json:"poDescription" gorm:"column:poDescription"`
	PoSKU         string `json:"poSKU" gorm:"column:poSKU"`
	PoHSN         string `json:"poHSN" gorm:"column:poHSN"`
	PoQuantity    string `json:"poQuantity" gorm:"column:poQuantity"`
	PoPrice       string `json:"poPrice" gorm:"column:poPrice"`
	PoDiscPercent string `json:"poDiscPercent" gorm:"column:poDiscPercent"`
	PoDisc        string `json:"poDisc" gorm:"column:poDisc"`
	PoTotalPrice  string `json:"poTotalPrice" gorm:"column:poTotalPrice"`
	CreatedAt     string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy     string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt     string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy     string `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete      bool   `json:"isDelete" gorm:"column:isDelete"`
}
