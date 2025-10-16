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
	Products        []PurchaseOrderProduct `gorm:"-" json:"products"` // transient, not in DB
}

type PurchaseOrderProduct struct {
	POProductID     int    `gorm:"primaryKey;column:po_product_id" json:"poProductId"`
	PurchaseOrderID int    `gorm:"column:purchase_order_id" json:"purchaseOrderId"`
	CategoryID      int    `gorm:"column:category_id" json:"categoryId"`
	Description     string `gorm:"column:description" json:"description"`
	UnitPrice       string `gorm:"column:unit_price" json:"unitPrice"`
	Discount        string `gorm:"column:discount" json:"discount"`
	Quantity        string `gorm:"column:quantity" json:"quantity"`
	Total           string `gorm:"column:total" json:"total"`
	CreatedAt       string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy       string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt       string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy       string `gorm:"column:updatedBy" json:"updatedBy"`
}

type PurchaseOrderPayload struct {
	PurchaseOrderID int                       `json:"purchaseOrderId"`
	Supplier        struct{ SupplierId int }  `json:"supplier"`
	Branch          struct{ RefBranchId int } `json:"branch"`
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
	PurchaseOrderID int    `json:"purchaseOrderId"`
	SupplierID      int    `json:"supplierId"`
	BranchID        int    `json:"branchId"`
	SubTotal        string `json:"subTotal"`
	TotalDiscount   string `json:"totalDiscount"`
	TaxEnabled      bool   `json:"taxEnabled"`
	TaxPercentage   string `json:"taxPercentage"`
	TaxAmount       string `json:"taxAmount"`
	TotalAmount     string `json:"totalAmount"`
	CreditedDate    string `json:"creditedDate"`
	CreatedAt       string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy       string `gorm:"column:createdBy" json:"createdBy"`
}
