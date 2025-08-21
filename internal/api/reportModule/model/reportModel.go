package reportModel

type ProductsReportPayload struct {
	FromDate         string `json:"fromDate"`
	ToDate           string `json:"toDate"`
	SearchField      string `json:"searchField"`
	PurchaseOrderId  int    `json:"purchaseOrderIdDropDown"`
	SupplierId       int    `json:"supplierIdDropDown"`
	PaginationOffset int    `json:"paginationOffset"`
	PaginationLimit  int    `json:"paginationLimit"`
}

type PurchaseOrderResponse struct {
	DummyProductsIdPK  int    `json:"dummyProductsId" gorm:"column:dummyProductsId;primaryKey;autoIncrement"`
	ProductName        string `json:"productName" gorm:"column:productName"`
	RefCategoryId      int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	RefSubCategoryId   int    `json:"refSubCategoryId" gorm:"column:refSubCategoryId"`
	HSNCode            string `json:"HSNCode" gorm:"column:HSNCode"`
	Price              string `json:"price" gorm:"column:price"`
	DiscountPercentage string `json:"discountPercentage" gorm:"column:discountPercentage"`
	DiscountAmount     string `json:"discountAmount" gorm:"column:discountAmount"`
	IsReceived         string `json:"isReceived" gorm:"column:isReceived"`
	AcceptanceStatus   string `json:"acceptanceStatus" gorm:"column:acceptanceStatus"`
	CreatedAt          string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy          string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt          string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy          string `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete           string `json:"isDelete" gorm:"column:isDelete"`
}
