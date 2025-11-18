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

type StockTransferRequest struct {
	ReceivedBranchDetails struct {
		SupplierId          int    `json:"supplierId"`
		SupplierName        string `json:"supplierName"`
		SupplierCompanyName string `json:"supplierCompanyName"`
		SupplierGSTNumber   string `json:"supplierGSTNumber"`
		SupplierCode        string `json:"supplierCode"`
	} `json:"receivedBranchDetails"`

	BranchDetails struct {
		BranchId      int    `json:"branchId"`
		BranchName    string `json:"branchName"`
		BranchEmail   string `json:"branchEmail"`
		BranchAddress string `json:"branchAddress"`
		BranchCode    string `json:"branchCode"`
	} `json:"branchDetails"`

	ProductDetails []struct {
		ProductName      string `json:"productName"`
		RefCategoryId    int    `json:"refCategoryid"`
		RefSubCategoryId int    `json:"refSubCategoryId"`
		HSNCode          string `json:"HSNCode"`
		SKU              string `json:"SKU"`
		PurchaseQuantity string `json:"purchaseQuantity"`
		PurchasePrice    string `json:"purchasePrice"`
		DiscountPrice    string `json:"discountPrice"`
		DiscountAmount   string `json:"discountAmount"`
		TotalAmount      string `json:"totalAmount"`
		IsReceived       bool   `json:"isReceived"`
		AcceptanceStatus string `json:"acceptanceStatus"`
		CreatedAt        string `json:"createdAt"`
		CreatedBy        string `json:"createdBy"`
		UpdatedAt        string `json:"updatedAt"`
		UpdatedBy        string `json:"updatedBy"`
		IsDelete         bool   `json:"isDelete"`
	} `json:"productDetails"`

	TotalSummary struct {
		PoNumber        string `json:"poNumber"`
		BranchId        int    `json:"branchId"`
		Status          int    `json:"status"`
		ModeOfTransport string `json:"modeOfTransport"`
		SubTotal        string `json:"subTotal"`
		DiscountOverall string `json:"discountOverall"`
		PayAmount       string `json:"payAmount"`
		TotalAmount     string `json:"totalAmount"`
		PaymentPending  string `json:"paymentPending"`
		CreatedAt       string `json:"createdAt"`
		CreatedBy       string `json:"createdBy"`
		UpdatedAt       string `json:"updatedAt"`
		UpdatedBy       string `json:"updatedBy"`
		IsDelete        bool   `json:"isDelete"`
	} `json:"totalSummary"`
}

type StockTransferItem struct {
	StockTransferItemID int    `gorm:"column:stock_transfer_item_id;primaryKey" json:"stockTransferItemId"`
	StockTransferID     int    `gorm:"column:stock_transfer_id" json:"stockTransferId"`
	ProductInstanceID   int    `gorm:"column:product_instance_id" json:"productInstanceId"`
	ProductName         string `gorm:"column:product_name" json:"productName"`
	SKU                 string `gorm:"column:sku" json:"sku"`
	IsReceived          bool   `gorm:"column:is_received" json:"isReceived"`
	AcceptanceStatus    string `gorm:"column:acceptance_status" json:"acceptanceStatus"`
}

func (StockTransferItem) TableName() string {
	return `"purchaseOrderMgmt"."Inventory_StockTransferItems"`
}

type StockTransfer struct {
	StockTransferID   int                 `gorm:"column:stock_transfer_id;primaryKey" json:"stockTransferId"`
	FromBranchID      int                 `gorm:"column:from_branch_id" json:"fromBranchId"`
	FromBranchName    string              `gorm:"column:from_branch_name" json:"fromBranchName"`
	FromBranchEmail   string              `gorm:"column:from_branch_email" json:"fromBranchEmail"`
	FromBranchAddress string              `gorm:"column:from_branch_address" json:"fromBranchAddress"`
	ToBranchID        int                 `gorm:"column:to_branch_id" json:"toBranchId"`
	ToBranchName      string              `gorm:"column:to_branch_name" json:"toBranchName"`
	ToBranchEmail     string              `gorm:"column:to_branch_email" json:"toBranchEmail"`
	ToBranchAddress   string              `gorm:"column:to_branch_address" json:"toBranchAddress"`
	ModeOfTransport   string              `gorm:"column:mode_of_transport" json:"modeOfTransport"`
	SubTotal          string              `gorm:"column:sub_total" json:"subTotal"`
	DiscountOverall   string              `gorm:"column:discount_overall" json:"discountOverall"`
	TotalAmount       string              `gorm:"column:total_amount" json:"totalAmount"`
	PaymentPending    string              `gorm:"column:payment_pending" json:"paymentPending"`
	PoNumber          string              `gorm:"column:po_number" json:"poNumber"`
	Status            int                 `gorm:"column:status" json:"status"`
	CreatedAt         string              `gorm:"column:created_at" json:"createdAt"`
	CreatedBy         string              `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt         string              `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedBy         string              `gorm:"column:updated_by" json:"updatedBy"`
	IsDelete          bool                `gorm:"column:is_delete" json:"isDelete"`
	Items             []StockTransferItem `json:"items" gorm:"-"`
}

func (StockTransfer) TableName() string {
	return `"purchaseOrderMgmt"."Inventory_StockTransfers"`
}
