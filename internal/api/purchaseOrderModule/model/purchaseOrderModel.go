package purchaseOrderModel

type SupplierDetails struct {
	SupplierID            int    `json:"supplierId"`
	SupplierName          string `json:"supplierName"`
	SupplierCompanyName   string `json:"supplierCompanyName"`
	SupplierGSTNumber     string `json:"supplierGSTNumber"`
	SupplierAddress       string `json:"supplierAddress"`
	SupplierPaymentTerms  string `json:"supplierPaymentTerms"`
	SupplierEmail         string `json:"supplierEmail"`
	SupplierContactNumber string `json:"supplierContactNumber"`
}

type BranchDetails struct {
	BranchID      int    `json:"branchId"`
	BranchName    string `json:"branchName"`
	BranchEmail   string `json:"branchEmail"`
	BranchAddress string `json:"branchAddress"`
}

type ProductDetails struct {
	ProductName      string `gorm:"column:productName" json:"productName"`
	RefCategoryID    int    `gorm:"column:refCategoryid" json:"refCategoryid"`
	RefSubCategoryID int    `gorm:"column:refSubCategoryId" json:"refSubCategoryId"`
	HSNCode          string `gorm:"column:HSNCode" json:"HSNCode"`
	PurchaseQuantity string `gorm:"column:purchaseQuantity" json:"purchaseQuantity"`
	PurchasePrice    string `gorm:"column:purchasePrice" json:"purchasePrice"`
	DiscountPrice    string `gorm:"column:discountPrice" json:"discountPrice"`
	DiscountAmount   string `gorm:"column:discountAmount" json:"discountAmount"`
	TotalAmount      string `gorm:"column:totalAmount" json:"totalAmount"`
	IsReceived       bool   `gorm:"column:isReceived" json:"isReceived"`
	AcceptanceStatus string `gorm:"column:acceptanceStatus" json:"acceptanceStatus"`
	CreatedAt        string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy        string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt        string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete         bool   `gorm:"column:isDelete" json:"isDelete"`
}

type TotalSummary struct {
	PONumber        string `json:"poNumber"`
	SupplierID      int    `json:"supplierId"`
	BranchID        int    `json:"branchId"`
	Status          int    `json:"status"`
	ExpectedDate    string `json:"expectedDate"`
	ModeOfTransport string `json:"modeOfTransport"`
	SubTotal        string `json:"subTotal"`
	DiscountOverall string `json:"discountOverall"`
	PayAmount       string `json:"payAmount"`
	IsTaxApplied    bool   `json:"isTaxApplied"`
	TaxPercentage   string `json:"taxPercentage"`
	TaxedAmount     string `json:"taxedAmount"`
	TotalAmount     string `json:"totalAmount"`
	TotalPaid       string `json:"totalPaid"`
	PaymentPending  string `json:"paymentPending"`
	CreatedAt       string `json:"createdAt"`
	CreatedBy       string `json:"createdBy"`
	UpdatedAt       string `json:"updatedAt"`
	UpdatedBy       string `json:"updatedBy"`
	IsDelete        bool   `json:"isDelete"`
}

type CreatePORequest struct {
	SupplierDetails SupplierDetails  `json:"supplierDetails"`
	BranchDetails   BranchDetails    `json:"branchDetails"`
	ProductDetails  []ProductDetails `json:"productDetails"`
	TotalSummary    TotalSummary     `json:"totalSummary"`
	PurchaseOrderID int              `json:"purchaseOrderId"`
}

type CreatePurchaseOrder struct {
	PurchaseOrderID int    `gorm:"column:purchaseOrderId;primaryKey;autoIncrement"`
	PONumber        string `gorm:"column:poNumber"`
	SupplierID      int    `gorm:"column:supplierId"`
	BranchID        int    `gorm:"column:branchId"`
	Status          int    `gorm:"column:status"`
	ExpectedDate    string `gorm:"column:expectedDate"`
	ModeOfTransport string `gorm:"column:modeOfTransport"`
	SubTotal        string `gorm:"column:subTotal"`
	DiscountOverall string `gorm:"column:discountOverall"`
	PayAmount       string `gorm:"column:payAmount"`
	IsTaxApplied    bool   `gorm:"column:isTaxApplied"`
	TaxPercentage   string `gorm:"column:taxPercentage"`
	TaxedAmount     string `gorm:"column:taxedAmount"`
	TotalAmount     string `gorm:"column:totalAmount"`
	TotalPaid       string `gorm:"column:totalPaid"`
	PaymentPending  string `gorm:"column:paymentPending"`
	CreatedAt       string `gorm:"column:createdAt"`
	CreatedBy       string `gorm:"column:createdBy"`
	UpdatedAt       string `gorm:"column:updatedAt"`
	UpdatedBy       string `gorm:"column:updatedBy"`
	IsDelete        string `gorm:"column:isDelete"`
}

func (CreatePurchaseOrder) TableName() string {
	return `"purchaseOrder"."CreatePurchaseOrder"`
}

type PurchaseOrderItem struct {
	ItemID           int    `gorm:"column:itemId;primaryKey;autoIncrement"`
	PurchaseOrderID  int    `gorm:"column:purchaseOrderId"`
	ProductName      string `gorm:"column:productName"`
	RefCategoryID    int    `gorm:"column:refCategoryid"`
	RefSubCategoryID int    `gorm:"column:refSubCategoryId"`
	HSNCode          string `gorm:"column:HSNCode"`
	PurchaseQuantity string `gorm:"column:purchaseQuantity"`
	PurchasePrice    string `gorm:"column:purchasePrice"`
	DiscountPrice    string `gorm:"column:discountPrice"`
	DiscountAmount   string `gorm:"column:discountAmount"`
	TotalAmount      string `gorm:"column:totalAmount"`
	IsReceived       bool   `gorm:"column:isReceived"`
	AcceptanceStatus string `gorm:"column:acceptanceStatus"`
	CreatedAt        string `gorm:"column:createdAt"`
	CreatedBy        string `gorm:"column:createdBy"`
	UpdatedAt        string `gorm:"column:updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy"`
	IsDelete         bool   `gorm:"column:isDelete"`
}

func (PurchaseOrderItem) TableName() string {
	return `"purchaseOrder"."PurchaseOrderItemsInitial"`
}

type OrderRow struct {
	PurchaseOrderID int    `gorm:"column:purchaseOrderId"`
	PONumber        string `gorm:"column:poNumber"`
	SupplierID      int    `gorm:"column:supplierId"`
	BranchID        int    `gorm:"column:branchId"`
	Status          int    `gorm:"column:status"`
	ExpectedDate    string `gorm:"column:expectedDate"`
	ModeOfTransport string `gorm:"column:modeOfTransport"`
	SubTotal        string `gorm:"column:subTotal"`
	DiscountOverall string `gorm:"column:discountOverall"`
	PayAmount       string `gorm:"column:payAmount"`
	IsTaxApplied    bool   `gorm:"column:isTaxApplied"`
	TaxPercentage   string `gorm:"column:taxPercentage"`
	TaxedAmount     string `gorm:"column:taxedAmount"`
	TotalAmount     string `gorm:"column:totalAmount"`
	TotalPaid       string `gorm:"column:totalPaid"`
	PaymentPending  string `gorm:"column:paymentPending"`
	CreatedAt       string `gorm:"column:createdAt"`
	CreatedBy       string `gorm:"column:createdBy"`
	UpdatedAt       string `gorm:"column:updatedAt"`
	UpdatedBy       string `gorm:"column:updatedBy"`
	IsDelete        bool   `gorm:"column:isDelete"`

	// Supplier
	SupplierName          string `gorm:"column:supplierName"`
	SupplierCompanyName   string `gorm:"column:supplierCompanyName"`
	SupplierGSTNumber     string `gorm:"column:supplierGSTNumber"`
	SupplierEmail         string `gorm:"column:supplierEmail"`
	SupplierContactNumber string `gorm:"column:supplierContactNumber"`
	SupplierPaymentTerms  string `gorm:"column:supplierPaymentTerms"`
	SupplierAddress       string `gorm:"column:supplierAddress"`

	// Branch
	BranchName    string `gorm:"column:branchName"`
	BranchEmail   string `gorm:"column:branchEmail"`
	BranchAddress string `gorm:"column:branchAddress"`
}

type ProductsDummyAcceptance struct {
	DummyProductsID  int    `gorm:"column:dummyProductsId;primaryKey;autoIncrement"`
	PurchaseOrderID  int    `gorm:"column:purchaseOrderId"`
	ProductName      string `gorm:"column:productName"`
	RefCategoryID    int    `gorm:"column:refCategoryId"`
	RefSubCategoryID int    `gorm:"column:refSubCategoryId"`
	HSNCode          string `gorm:"column:HSNCode"`
	DummySKU         string `gorm:"column:dummySKU"`
	Price            string `gorm:"column:price"`
	DiscountPercent  string `gorm:"column:discountPercentage"`
	DiscountAmount   string `gorm:"column:discountAmount"`
	IsReceived       string `gorm:"column:isReceived"`
	AcceptanceStatus string `gorm:"column:acceptanceStatus"`
	CreatedAt        string `gorm:"column:createdAt"`
	CreatedBy        string `gorm:"column:createdBy"`
	UpdatedAt        string `gorm:"column:updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy"`
	IsDelete         string `gorm:"column:isDelete"`
}

func (ProductsDummyAcceptance) TableName() string {
	return `"purchaseOrder"."ProductsDummyAcceptance"`
}

type Product struct {
	ProductID           int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                string `gorm:"column:name" json:"name"`
	SKU                 string `gorm:"column:sku" json:"sku"`
	GTIN                string `gorm:"column:gtin" json:"gtin"`
	CategoryID          int    `gorm:"column:category_id" json:"category"`
	SubCategoryID       int    `gorm:"column:subcategory_id" json:"subcategory"`
	Description         string `gorm:"column:description" json:"description"`
	DetailedDescription string `gorm:"column:detailed_description" json:"detailedDescription"`
	Price               string `gorm:"column:price" json:"price"`
	MRP                 string `gorm:"column:mrp" json:"mrp"`
	Cost                string `gorm:"column:cost" json:"cost"`
	SplPrice            string `gorm:"column:spl_price" json:"splPrice"`
	StartDate           string `gorm:"column:start_date" json:"startDate"`
	EndDate             string `gorm:"column:end_date" json:"endDate"`
	TaxClass            string `gorm:"column:tax_class" json:"taxClass"`
	ProductImage        string `gorm:"column:product_image" json:"productImage"`
	Featured            bool   `gorm:"column:featured" json:"featured"`
	CreatedAt           string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy           string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt           string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy           string `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete            string `gorm:"column:isDelete" json:"isDelete"`
}
