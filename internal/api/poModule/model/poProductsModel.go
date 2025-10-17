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
	PurchaseOrderId int    `gorm:"column:purchase_order_id;primaryKey"`
	SupplierID      int    `gorm:"column:supplier_id"`
	BranchID        int    `gorm:"column:branch_id"`
	TotalAmount     string `gorm:"column:total_amount"`
	CreditedDate    string `gorm:"column:credited_date"`
	InvoiceNumber   string `gorm:"column:invoiceNumber"` // 👈 match exact DB column
	InvoiceStatus   bool   `gorm:"column:invoiceStatus"` // 👈 match exact DB column
	CreatedAt       string `gorm:"column:createdAt"`
	CreatedBy       string `gorm:"column:createdBy"`
	UpdatedAt       string `gorm:"column:updatedAt"`
	UpdatedBy       string `gorm:"column:updatedBy"`
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
