package poModuleModel

type PurchaseOrder struct {
	PurchaseOrderID     int                    `gorm:"primaryKey;column:purchase_order_id" json:"purchase_order_id"`
	SupplierID          int                    `gorm:"column:supplier_id" json:"supplier_id"`
	BranchID            int                    `gorm:"column:branch_id" json:"branch_id"`
	SubTotal            string                 `gorm:"column:sub_total" json:"sub_total"`
	TotalDiscount       string                 `gorm:"column:total_discount" json:"total_discount"`
	TaxEnabled          bool                   `gorm:"column:tax_enabled" json:"tax_enabled"`
	TaxPercentage       string                 `gorm:"column:tax_percentage" json:"tax_percentage"`
	TaxAmount           string                 `gorm:"column:tax_amount" json:"tax_amount"`
	TotalAmount         string                 `gorm:"column:total_amount" json:"total_amount"`
	CreditedDate        string                 `gorm:"column:credited_date" json:"credited_date"`
	CreatedAt           string                 `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy           string                 `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt           string                 `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy           string                 `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete            bool                   `gorm:"column:isDelete" json:"isDelete"`
	PurchaseOrderNumber string                 `gorm:"column:purchaseOrderNumber" json:"purchaseOrderNumber"`
	InvoiceStatus       bool                   `gorm:"column:invoiceStatus" json:"invoiceStatus"`
	Products            []PurchaseOrderProduct `gorm:"-" json:"products"`
	InvoiceFinalNumber  string                 `gorm:"column:invoiceFinalNumber" json:"invoiceFinalNumber"`
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
	PurchaseOrderID     int             `json:"purchaseOrderId"`
	PurchaseOrderNumber string          `json:"purchaseOrderNumber"`
	Supplier            SupplierDetails `json:"supplier"`
	Branch              BranchDetails   `json:"branch"`
	Summary             struct {
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

	SubTotal            string `gorm:"column:sub_total"`
	TotalDiscount       string `gorm:"column:total_discount"`
	TaxEnabled          bool   `gorm:"column:tax_enabled"`
	TaxPercentage       string `gorm:"column:tax_percentage"`
	TaxAmount           string `gorm:"column:tax_amount"`
	TotalAmount         string `gorm:"column:total_amount"`
	CreditedDate        string `gorm:"column:credited_date"`
	PurchaseOrderNumber string `gorm:"column:purchaseOrderNumber"`
	CreatedAt           string `gorm:"column:createdAt"`
	CreatedBy           string `gorm:"column:createdBy"`
}

// type PurchaseOrderListResponse struct {
// 	PurchaseOrderID     int64  `json:"purchase_order_id" db:"purchase_order_id"`
// 	PurchaseOrderNumber string `json:"purchase_order_number" db:"purchase_order_number"`
// 	Status              string `json:"status" db:"status"`

// 	// Quantities
// 	TotalOrderedQuantity  int64 `json:"total_ordered_quantity" db:"total_ordered_quantity"`
// 	TotalAcceptedQuantity int64 `json:"total_accepted_quantity" db:"total_accepted_quantity"`
// 	TotalRejectedQuantity int64 `json:"total_rejected_quantity" db:"total_rejected_quantity"`

// 	// Totals
// 	TotalAmount   float64 `json:"total_amount" db:"total_amount"`
// 	TaxableAmount float64 `json:"taxable_amount" db:"taxable_amount"`

// 	CreatedAt string `json:"created_at" db:"created_at"`

// 	// Supplier & Branch
// 	SupplierID   int64  `json:"supplier_id" db:"supplier_id"`
// 	BranchID     int64  `json:"branch_id" db:"branch_id"`
// 	SupplierName string `json:"supplier_name" db:"supplier_name"`
// 	BranchName   string `json:"branch_name" db:"branch_name"`

// 	// Category & Product Info
// 	Description         string `json:"description" db:"description"`
// 	InitialCategoryID   int64  `json:"initial_category_id" db:"initial_category_id"`
// 	InitialCategoryName string `json:"initial_category_name" db:"initial_category_name"`

// 	// Product-level breakdown
// 	ProductOrderedQuantity  int64 `json:"product_ordered_quantity" db:"product_ordered_quantity"`
// 	ProductAcceptedQuantity int64 `json:"product_accepted_quantity" db:"product_accepted_quantity"`
// 	ProductRejectedQuantity int64 `json:"product_rejected_quantity" db:"product_rejected_quantity"`
// }

type PurchaseOrderListResponseFlat struct {
	PurchaseOrderID       int
	PurchaseOrderNumber   string
	Status                string
	TotalOrderedQuantity  int
	TotalAcceptedQuantity int
	TotalRejectedQuantity int
	TotalAmount           float64
	TaxableAmount         float64
	CreatedAt             string
	SupplierID            int
	SupplierName          string
	BranchID              int
	BranchName            string
	POProductID           int
	Description           string
	CategoryID            int
	UnitPrice             string
	Discount              string
	Quantity              string
	Total                 string
	ProductCreatedAt      string
	ProductCreatedBy      string
	ProductUpdatedAt      string
	ProductUpdatedBy      string
	InitialCategoryID     int
	InitialCategoryName   string
	InitialCategoryCode   string
	IsDelete              bool
	CatCreatedAt          string
	CatCreatedBy          string
	CatUpdatedAt          string
	CatUpdatedBy          string
}

type CategoryDetails struct {
	InitialCategoryId   int     `json:"initialCategoryId"`
	InitialCategoryName string  `json:"initialCategoryName"`
	InitialCategoryCode string  `json:"initialCategoryCode"`
	IsDelete            *bool   `json:"isDelete"`
	CreatedAt           *string `json:"createdAt"`
	CreatedBy           *string `json:"createdBy"`
	UpdatedAt           *string `json:"updatedAt"`
	UpdatedBy           *string `json:"updatedBy"`
}

type PurchaseOrderProductLatest struct {
	PoProductId      int    `json:"poProductId"`
	PurchaseOrderId  int    `json:"purchaseOrderId"`
	CategoryId       int    `json:"categoryId"`
	Description      string `json:"description"`
	UnitPrice        string `json:"unitPrice"`
	Discount         string `json:"discount"`
	Quantity         string `json:"quantity"`
	Total            string `json:"total"`
	CreatedAt        string `json:"createdAt"`
	CreatedBy        string `json:"createdBy"`
	UpdatedAt        string `json:"updatedAt"`
	UpdatedBy        string `json:"updatedBy"`
	AcceptedQuantity string `json:"accepted_quantity"`
	RejectedQuantity string `json:"rejected_quantity"`

	InitialCategoryId   int     `json:"-"` // hidden in JSON
	InitialCategoryName string  `json:"-"`
	InitialCategoryCode string  `json:"-"`
	IsDelete            *bool   `json:"-"`
	CategoryCreatedAt   *string `json:"-"`
	CategoryCreatedBy   *string `json:"-"`
	CategoryUpdatedAt   *string `json:"-"`
	CategoryUpdatedBy   *string `json:"-"`

	CategoryDetails CategoryDetails `json:"categoryDetails" gorm:"-"`
}

type PurchaseOrderListResponse struct {
	PurchaseOrderId       int                          `json:"purchaseOrderId"`
	PurchaseOrderNumber   string                       `json:"purchaseOrderNumber"`
	Status                string                       `json:"status"`
	TotalOrderedQuantity  int64                        `json:"totalOrderedQuantity"`
	TotalAcceptedQuantity int64                        `json:"totalAcceptedQuantity"`
	TotalRejectedQuantity int64                        `json:"totalRejectedQuantity"`
	TotalAmount           string                       `json:"totalAmount"`
	CreatedAt             string                       `json:"createdAt"`
	TaxableAmount         string                       `json:"taxableAmount"`
	SupplierId            int                          `json:"supplierId"`
	SupplierName          string                       `json:"supplierName"`
	BranchId              int                          `json:"branchId"`
	BranchName            string                       `json:"branchName"`
	Products              []PurchaseOrderProductLatest `json:"products" gorm:"-"` // ✅ correct
}

type PurchaseOrderAcceptedProductResponse struct {
	ProductInstanceID  int    `json:"productInstanceId" gorm:"column:product_instance_id"`
	PoProductID        int    `json:"poProductId" gorm:"column:po_product_id"`
	LineNumber         string `json:"lineNumber" gorm:"column:line_number"`
	ReferenceNumber    string `json:"referenceNumber" gorm:"column:reference_number"`
	ProductDescription string `json:"productDescription" gorm:"column:product_description"`
	Discount           string `json:"discount" gorm:"column:discount"`
	UnitPrice          string `json:"unitPrice" gorm:"column:unit_price"`
	DiscountPrice      string `json:"discountPrice" gorm:"column:discount_price"`
	Margin             string `json:"margin" gorm:"column:margin"`
	TotalAmount        string `json:"totalAmount" gorm:"column:total_amount"`
	CategoryID         int    `json:"categoryId" gorm:"column:category_id"`
	SubCategoryID      int    `json:"subCategoryId" gorm:"column:sub_category_id"`
	Status             string `json:"status" gorm:"column:status"`
	CreatedAt          string `json:"createdAt" gorm:"column:created_at"`
	CreatedBy          string `json:"createdBy" gorm:"column:created_by"`
	UpdatedAt          string `json:"updatedAt" gorm:"column:updated_at"`
	UpdatedBy          string `json:"updatedBy" gorm:"column:updated_by"`
	IsDelete           bool   `json:"isDelete" gorm:"column:is_delete"`
	ProductName        string `json:"productName" gorm:"column:product_name"`
	PurchaseOrderId    int    `json:"purchaseOrderId" gorm:"column:purchase_order_id"`
	SKU                string `json:"sku" gorm:"column:sku"`
	ProductBranchId    int    `json:"productBranchId" gorm:"column:product_branch_id"`
	Quantity           string `json:"quantity" gorm:"column:quantity"`
	InvoiceFinalNumber string `json:"invoiceFinalNumber" gorm:"column:invoice_final_number"`
	CategoryName       string `json:"categoryName" gorm:"column:category_name"`
	SubCategoryName    string `json:"subCategoryName" gorm:"column:sub_category_name"`
	BranchName         string `json:"branchName" gorm:"column:branch_name"`
}

type PurchaseOrderAcceptedProduct struct {
	ProductInstanceId  int     `json:"product_instance_id" gorm:"column:product_instance_id"`
	PoProductId        int     `json:"po_product_id" gorm:"column:po_product_id"`
	PurchaseOrderId    int     `json:"purchaseOrderId" gorm:"column:purchaseOrderId"`
	LineNumber         string  `json:"line_number" gorm:"column:line_number"`
	ReferenceNumber    string  `json:"reference_number" gorm:"column:reference_number"`
	ProductDescription string  `json:"product_description" gorm:"column:product_description"`
	Discount           string  `json:"discount" gorm:"column:discount"`
	UnitPrice          string  `json:"unit_price" gorm:"column:unit_price"`
	DiscountPrice      string  `json:"discount_price" gorm:"column:discount_price"`
	Margin             string  `json:"margin" gorm:"column:margin"`
	TotalAmount        string  `json:"total_amount" gorm:"column:total_amount"`
	CategoryId         int     `json:"category_id" gorm:"column:category_id"`
	SubCategoryId      int     `json:"sub_category_id" gorm:"column:sub_category_id"`
	ProductName        string  `json:"product_name" gorm:"column:product_name"`
	SKU                string  `json:"SKU" gorm:"column:SKU"`
	Status             string  `json:"status" gorm:"column:status"`
	CreatedAt          string  `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy          string  `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt          string  `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy          string  `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete           bool    `json:"isDelete" gorm:"column:isDelete"`
	ProductBranchId    *int    `json:"productBranchId" gorm:"column:productBranchId"`
	Quantity           *string `json:"quantity" gorm:"column:quantity"`
}

type AcceptedRow struct {
	PoProductId        int     `gorm:"column:po_product_id"`
	LineNumber         int     `gorm:"column:line_number"`
	ReferenceNumber    string  `gorm:"column:reference_number"`
	ProductDescription string  `gorm:"column:product_description"`
	Discount           float64 `gorm:"column:discount"`
	UnitPrice          float64 `gorm:"column:unit_price"`
	DiscountPrice      float64 `gorm:"column:discount_price"`
	Margin             float64 `gorm:"column:margin"`
	TotalAmount        string  `gorm:"column:total_amount"`
	CategoryId         int     `gorm:"column:category_id"`
	CategoryName       string  `gorm:"column:category_name"`
	SubCategoryId      int     `gorm:"column:sub_category_id"`
	SubCategoryName    string  `gorm:"column:sub_category_name"`
	ProductName        string  `gorm:"column:product_name"`
	SKU                string  `gorm:"column:SKU"`
}

type ProductRow struct {
	PoProductId     int     `gorm:"column:po_product_id"`
	LineNumber      int     `gorm:"column:line_number"`
	Description     string  `gorm:"column:description"`
	CategoryId      int     `gorm:"column:category_id"`
	CategoryName    string  `gorm:"column:category_name"`
	SubCategoryId   int     `gorm:"column:sub_category_id"`
	SubCategoryName string  `gorm:"column:sub_category_name"`
	Quantity        int     `gorm:"column:quantity"`
	UnitPrice       float64 `gorm:"column:unit_price"`
}

type DialogRow struct {
	SNo                int     `json:"sNo"`
	LineNumber         int     `json:"lineNumber"`
	ReferenceNumber    string  `json:"referenceNumber"`
	ProductDescription string  `json:"productDescription"`
	Discount           float64 `json:"discount"`
	Price              float64 `json:"price"`
	DiscountPrice      float64 `json:"discountPrice"`
	Margin             float64 `json:"margin"`
	TotalAmount        string  `json:"totalAmount"`
}

type Product struct {
	SNo             int         `json:"sNo"`
	LineNumber      int         `json:"lineNumber"`
	ProductName     string      `json:"productName"`
	Brand           string      `json:"brand"`
	CategoryId      int         `json:"categoryId"`
	SubCategoryId   int         `json:"subCategoryId"`
	TaxClass        string      `json:"taxClass"`
	Quantity        int         `json:"quantity"`
	Cost            float64     `json:"cost"`
	ProfitMargin    float64     `json:"profitMargin"`
	SellingPrice    float64     `json:"sellingPrice"`
	Mrp             float64     `json:"mrp"`
	DiscountPercent float64     `json:"discountPercent"`
	DiscountPrice   float64     `json:"discountPrice"`
	TotalAmount     string      `json:"totalAmount"`
	DialogRows      []DialogRow `json:"dialogRows"`
}

type PurchaseOrderDetailsResponse struct {
	PurchaseOrderId    int       `json:"purchaseOrderId"`
	InvoiceFinalNumber string    `json:"invoiceFinalNumber"`
	Products           []Product `json:"products"`
}
