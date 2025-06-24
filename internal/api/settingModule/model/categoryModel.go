package model

type Category struct {
	RefCategoryId int    `json:"refCategoryId" gorm:"column:refCategoryid;primaryKey;autoIncrement"`
	CategoryName  string `json:"categoryName" gorm:"column:categoryName"`
	CategoryCode  string `json:"categoryCode" gorm:"column:categoryCode"`
	IsActive      bool   `json:"isActive" gorm:"column:isActive"`
	IsDelete      bool   `json:"isDelete" gorm:"column:isDelete"`
	CreatedAt     string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy     string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt     string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy     string `json:"updatedBy" gorm:"column:updatedBy"`
}

type SubCategory struct {
	RefSubCategoryId int    `json:"refSubCategoryId" gorm:"column:refSubCategoryId;primaryKey;autoIncrement"`
	SubCategoryName  string `json:"subCategoryName" gorm:"column:subCategoryName"`
	RefCategoryId    int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	SubCategoryCode  string `json:"subCategoryCode" gorm:"column:subCategoryCode"`
	IsActive         bool   `json:"isActive" gorm:"column:isActive"`
	CreatedAt        string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy        string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt        string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy        string `json:"updatedBy" gorm:"column:updatedBY"` // 🔧 Fix is here
	IsDelete         bool   `json:"isDelete" gorm:"column:isDelete"`
}
