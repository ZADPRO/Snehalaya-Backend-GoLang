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

type PurchaseOrderProduct struct {
	ProductInstanceID  int    `json:"productInstanceId"`
	PoProductID        int    `json:"poProductId"`
	LineNumber         string `json:"lineNumber"`
	ReferenceNumber    string `json:"referenceNumber"`
	ProductDescription string `json:"productDescription"`
	Discount           string `json:"discount"`
	UnitPrice          string `json:"unitPrice"`
	DiscountPrice      string `json:"discountPrice"`
	Margin             string `json:"margin"`
	TotalAmount        string `json:"totalAmount"`
	CategoryID         int    `json:"categoryId"`
	SubCategoryID      int    `json:"subCategoryId"`
	Status             string `json:"status"`
	CreatedAt          string `json:"createdAt"`
	CreatedBy          string `json:"createdBy"`
	UpdatedAt          string `json:"updatedAt"`
	UpdatedBy          string `json:"updatedBy"`
	IsDelete           bool   `json:"isDelete"`
	ProductName        string `json:"productName"`
	PurchaseOrderID    int    `json:"purchaseOrderId"`
	SKU                string `json:"SKU"`
	ProductBranchID    int    `json:"productBranchId"`
	Quantity           string `json:"quantity"`
	BranchName         string `json:"branchName" gorm:"-"` // not in table, for response
	ProductBranchid    int    `json:"productBranchId" gorm:"productBranchId"`
}
