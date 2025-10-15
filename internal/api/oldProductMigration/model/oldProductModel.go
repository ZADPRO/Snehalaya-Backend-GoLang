package oldProductMigrationModel

type MigrateOldProductToDbModel struct {
	Id            int    `json:"id" gorm:"column:id;primaryKey; autoIncrement"`
	Unit          string `json:"unit" gorm:"column:unit"`
	ProductName   string `json:"productName" gorm:"column:productName"`
	SKU           string `json:"SKU" gorm:"column:SKU"`
	BrandId       int    `json:"brandId" gorm:"column:brandId"`
	Categoryid    int    `json:"categoryId" gorm:"column:categoryId"`
	SubCategoryId int    `json:"subCategoryid" gorm:"column:subCategoryid"`
	Quantity      string `json:"Quantity" gorm:"column:Quantity"`
	MRP           string `json:"MRP" gorm:"column:MRP"`
	Cost          string `json:"Cost" gorm:"column:Cost"`
	CreatedAt     string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy     string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt     string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy     string `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete      bool   `json:"isDelete" gorm:"column:isDelete"`
}
