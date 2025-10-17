package poModuleModel

type PurchaseOrder struct {
	PurchaseOrderID int                    `gorm:"primaryKey;column:purchase_order_id" json:"purchase_order_id"`
	SupplierID      int                    `gorm:"column:supplier_id" json:"supplier_id"`
	BranchID        int                    `gorm:"column:branch_id" json:"branch_id"`
	SubTotal        string                 `gorm:"column:sub_total" json:"sub_total"`
	TotalDiscount   string                 `gorm:"column:total_discount" json:"total_discount"`
	TaxEnabled      bool                   `gorm:"column:tax_enabled" json:"tax_enabled"`
	TaxPercentage   string                 `gorm:"column:tax_percentage" json:"tax_percentage"`
	TaxAmount       string                 `gorm:"column:tax_amount" json:"tax_amount"`
	TotalAmount     string                 `gorm:"column:total_amount" json:"total_amount"`
	CreditedDate    string                 `gorm:"column:credited_date" json:"credited_date"`
	CreatedAt       string                 `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy       string                 `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt       string                 `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy       string                 `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete        bool                   `gorm:"column:isDelete" json:"isDelete"`
	InvoiceNumber   string                 `gorm:"column:invoiceNumber" json:"invoiceNumber"`
	InvoiceStatus   bool                   `gorm:"column:invoiceStatus" json:"invoiceStatus"`
	Products        []PurchaseOrderProduct `gorm:"-" json:"products"` // transient, not in DB
}

type PurchaseOrderProduct struct {
	POProductID     int              `gorm:"primaryKey;column:po_product_id" json:"poProductId"`
	PurchaseOrderID int              `gorm:"column:purchase_order_id" json:"purchaseOrderId"`
	CategoryID      int              `gorm:"column:category_id" json:"categoryId"`
	Description     string           `gorm:"column:description" json:"description"`
	UnitPrice       string           `gorm:"column:unit_price" json:"unitPrice"`
	Discount        string           `gorm:"column:discount" json:"discount"`
	Quantity        string           `gorm:"column:quantity" json:"quantity"`
	Total           string           `gorm:"column:total" json:"total"`
	CreatedAt       string           `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy       string           `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt       string           `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy       string           `gorm:"column:updatedBy" json:"updatedBy"`
	CategoryDetails *InitialCategory `gorm:"-" json:"categoryDetails,omitempty"` // important
}

type SupplierDetails struct {
	SupplierId           int    `json:"supplierId"`
	SupplierName         string `json:"supplierName"`
	SupplierCompanyName  string `json:"supplierCompanyName"`
	SupplierCode         string `json:"supplierCode"`
	SupplierEmail        string `json:"supplierEmail"`
	SupplierMobile       string `json:"supplierMobile"`
	SupplierGSTNumber    string `json:"supplierGSTNumber"`
	SupplierPaymentTerms string `json:"supplierPaymentTerms"`
	// Add more fields as needed
}

type BranchDetails struct {
	RefBranchId   int    `json:"refBranchId"`
	RefBranchName string `json:"refBranchName"`
	RefBranchCode string `json:"refBranchCode"`
	RefLocation   string `json:"refLocation"`
	RefMobile     string `json:"refMobile"`
	RefEmail      string `json:"refEmail"`
	IsMainBranch  bool   `json:"isMainBranch"`
	IsActive      bool   `json:"isActive"`
	// Add more fields as needed
}

type InitialCategory struct {
	InitialCategoryId   int    `json:"initialCategoryId" gorm:"column:initialCategoryId;primaryKey;autoIncrement"`
	InitialCategoryName string `json:"initialCategoryName" gorm:"column:initialCategoryName"`
	InitialCategoryCode string `json:"initialCategoryCode" gorm:"column:initialCategoryCode"`
	IsDelete            bool   `json:"isDelete" gorm:"column:isDelete"`
	CreatedAt           string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy           string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt           string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy           string `json:"updatedBy" gorm:"column:updatedBy"`
}

type PurchaseOrderPayload struct {
	PurchaseOrderID int             `json:"purchaseOrderId"`
	InvoiceNumber   string          `json:"invoiceNumber"`
	Supplier        SupplierDetails `json:"supplier"`
	Branch          BranchDetails   `json:"branch"`
	Summary         struct {
		SubTotal      string `json:"subTotal"`
		TotalDiscount string `json:"totalDiscount"`
		TaxEnabled    bool   `json:"taxEnabled"`
		TaxPercentage string `json:"taxPercentage"`
		TaxAmount     string `json:"taxAmount"`
		TotalAmount   string `json:"totalAmount"`
	} `json:"summary"`
	CreditedDate string                 `json:"creditedDate"`
	Products     []PurchaseOrderProduct `json:"products"`
}

type PurchaseOrderResponse struct {
	PurchaseOrderID int    `gorm:"column:purchase_order_id"`
	SupplierID      int    `gorm:"column:supplierId"`      // po.supplier_id AS supplierId
	SupplierName    string `gorm:"column:supplierName"`    // s.supplierName AS supplierName
	SupplierCompany string `gorm:"column:supplierCompany"` // s.supplierCompanyName AS supplierCompany
	SupplierCode    string `gorm:"column:supplierCode"`    // s.supplierCode AS supplierCode
	SupplierEmail   string `gorm:"column:supplierEmail"`   // s.supplierEmail AS supplierEmail
	SupplierMobile  string `gorm:"column:supplierMobile"`  // s.supplierContactNumber AS supplierMobile
	SupplierGST     string `gorm:"column:supplierGST"`     // s.supplierGSTNumber AS supplierGST
	SupplierTerms   string `gorm:"column:supplierTerms"`   // s.supplierPaymentTerms AS supplierTerms

	BranchID       int    `gorm:"column:branchId"`       // po.branch_id AS branchId
	BranchName     string `gorm:"column:branchName"`     // b.refBranchName AS branchName
	BranchCode     string `gorm:"column:branchCode"`     // b.refBranchCode AS branchCode
	BranchLocation string `gorm:"column:branchLocation"` // b.refLocation AS branchLocation
	BranchMobile   string `gorm:"column:branchMobile"`   // b.refMobile AS branchMobile
	BranchEmail    string `gorm:"column:branchEmail"`    // b.refEmail AS branchEmail
	IsMainBranch   bool   `gorm:"column:isMainBranch"`   // b.isMainBranch AS isMainBranch
	IsActive       bool   `gorm:"column:isActive"`       // b.isActive AS isActive

	SubTotal      string `gorm:"column:sub_total"`
	TotalDiscount string `gorm:"column:total_discount"`
	TaxEnabled    bool   `gorm:"column:tax_enabled"`
	TaxPercentage string `gorm:"column:tax_percentage"`
	TaxAmount     string `gorm:"column:tax_amount"`
	TotalAmount   string `gorm:"column:total_amount"`
	CreditedDate  string `gorm:"column:credited_date"`
	InvoiceNumber string `gorm:"column:invoiceNumber"`
	CreatedAt     string `gorm:"column:createdAt"`
	CreatedBy     string `gorm:"column:createdBy"`
}
