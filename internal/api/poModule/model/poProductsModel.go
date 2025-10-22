package poModuleModel

type PurchaseOrderProductPayload struct {
	PoInvoiceNumber string                   `json:"poInvoiceNumber"`
	PoId            int                      `json:"poId"`
	SupplierId      int                      `json:"supplierId"`
	BranchId        int                      `json:"branchId"`
	TotalAmount     string                   `json:"totalAmount"`
	Products        []ProductPayloadProducts `json:"products"`
}

type ProductPayloadProducts struct {
	CategoryId  int    `json:"categoryId"`
	ProductName string `json:"productName"`
	OrderedQty  int    `json:"orderedQty"`
	ReceivedQty int    `json:"receivedQty"`
	RejectedQty int    `json:"rejectedQty"`
	UnitPrice   string `json:"unitPrice"`
	TotalPrice  int    `json:"totalPrice"`
	Status      string `json:"status"`
}

type PurchaseOrdersProducts struct {
	PurchaseOrderId     int    `gorm:"column:purchase_order_id;primaryKey"`
	SupplierID          int    `gorm:"column:supplier_id"`
	BranchID            int    `gorm:"column:branch_id"`
	TotalAmount         string `gorm:"column:total_amount"`
	CreditedDate        string `gorm:"column:credited_date"`
	PurchaseOrderNumber string `gorm:"column:purchaseOrderNumber"` // 👈 match exact DB column
	InvoiceStatus       bool   `gorm:"column:invoiceStatus"`       // 👈 match exact DB column
	CreatedAt           string `gorm:"column:createdAt"`
	CreatedBy           string `gorm:"column:createdBy"`
	UpdatedAt           string `gorm:"column:updatedAt"`
	UpdatedBy           string `gorm:"column:updatedBy"`
	InvoiceFinalNumber  string `gorm:"column:invoiceFinalNumber" json:"invoiceFinalNumber"`
}

type PurchaseOrderProducts struct {
	PoProductId      int    `gorm:"primaryKey;column:po_product_id"`
	PurchaseOrderID  int    `gorm:"column:purchase_order_id"`
	CategoryID       int    `gorm:"column:category_id"`
	Description      string `gorm:"column:description"`
	UnitPrice        string `gorm:"column:unit_price"`
	Quantity         string `gorm:"column:quantity"`
	AcceptedQuantity string `gorm:"column:accepted_quantity"`
	RejectedQuantity string `gorm:"column:rejected_quantity"`
	Status           string `gorm:"column:status"`
	AcceptedTotal    string `gorm:"column:accepted_total"`
	CreatedAt        string `gorm:"column:createdAt"` // 👈 matches DB
	CreatedBy        string `gorm:"column:createdBy"` // 👈 matches DB
	UpdatedAt        string `gorm:"column:updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy"`
}

type PurchaseOrderProductInstances struct {
	ProductInstanceId  int    `gorm:"primaryKey;column:product_instance_id"`
	PoProductID        string `gorm:"column:po_product_id"`
	SerialNo           string `gorm:"column:serial_no"`
	CategoryID         int    `gorm:"column:category_id"`
	ProductDescription string `gorm:"column:product_description"`
	UnitPrice          string `gorm:"column:unit_price"`
	Status             string `gorm:"column:status"`
	CreatedAt          string `gorm:"column:createdAt"` // 👈 matches DB
	CreatedBy          string `gorm:"column:createdBy"` // 👈 matches DB
}

type RejectedProducts struct {
	RejectedProductId  int    `gorm:"primaryKey;column:rejected_product_id"`
	PoProductID        int    `gorm:"column:po_product_id"`
	CategoryID         int    `gorm:"column:category_id"`
	ProductDescription string `gorm:"column:product_description"`
	UnitPrice          string `gorm:"column:unit_price"`
	RejectedQty        string `gorm:"column:rejected_qty"`
	Reason             string `gorm:"column:reason"`
	CreatedAt          string `gorm:"column:created_at"`
	CreatedBy          string `gorm:"column:created_by"`
}

type AcceptedProduct struct {
	PoProductID     int    `json:"po_product_id"`
	CategoryID      int    `json:"category_id"`
	ProductDesc     string `json:"product_description"`
	UnitPrice       string `json:"unit_price"`
	AcceptedQty     string `json:"accepted_quantity"`
	AcceptedTotal   string `json:"accepted_total"`
	OrderedQuantity string `json:"ordered_quantity"`
	OrderedTotal    string `json:"ordered_total"`
	Status          string `json:"status"`
	UpdatedAt       string `json:"updated_at"`
	UpdatedBy       string `json:"updated_by"`
}

type AcceptedPOResponse struct {
	PurchaseOrderID     int               `json:"purchase_order_id"`
	PurchaseOrderNumber string            `json:"purchaseOrderNumber"`
	BranchID            int               `json:"branch_id"`
	SupplierID          int               `json:"supplier_id"`
	TotalAmount         string            `json:"total_amount"`
	CreatedAt           string            `json:"created_at"`
	AcceptedProducts    []AcceptedProduct `json:"accepted_products"`
}
