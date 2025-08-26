package model

type Supplier struct {
	SupplierID             int    `json:"supplierId" gorm:"column:supplierId;primaryKey;autoIncrement"`
	SupplierName           string `json:"supplierName" gorm:"column:supplierName"`
	SupplierCompanyName    string `json:"supplierCompanyName" gorm:"column:supplierCompanyName"`
	SupplierCode           string `json:"supplierCode" gorm:"column:supplierCode"`
	SupplierEmail          string `json:"supplierEmail" gorm:"column:supplierEmail"`
	SupplierGSTNumber      string `json:"supplierGSTNumber" gorm:"column:supplierGSTNumber"`
	SupplierPaymentTerms   string `json:"supplierPaymentTerms" gorm:"column:supplierPaymentTerms"`
	SupplierBankACNumber   string `json:"supplierBankACNumber" gorm:"column:supplierBankACNumber"`
	SupplierIFSC           string `json:"supplierIFSC" gorm:"column:supplierIFSC"`
	SupplierBankName       string `json:"supplierBankName" gorm:"column:supplierBankName"`
	SupplierUPI            string `json:"supplierUPI" gorm:"column:supplierUPI"`
	SupplierIsActive       string `json:"supplierIsActive" gorm:"column:supplierIsActive"`
	SupplierContactNumber  string `json:"supplierContactNumber" gorm:"column:supplierContactNumber"`
	EmergencyContactName   string `json:"emergencyContactName" gorm:"column:emergencyContactName"`
	EmergencyContactNumber string `json:"emergencyContactNumber" gorm:"column:emergencyContactNumber"`
	SupplierDoorNumber     string `json:"supplierDoorNumber" gorm:"column:supplierDoorNumber"`
	SupplierStreet         string `json:"supplierStreet" gorm:"column:supplierStreet"`
	SupplierCity           string `json:"supplierCity" gorm:"column:supplierCity"`
	SupplierState          string `json:"supplierState" gorm:"column:supplierState"`
	SupplierCountry        string `json:"supplierCountry" gorm:"column:supplierCountry"`
	CreatedAt              string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy              string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt              string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy              string `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete               bool   `json:"isDelete" gorm:"column:isDelete"`
	CreditedDays           int    `json:"creditedDays" gorm:"column:creditedDays"`
}

type BulkDeleteRequest struct {
	IDs      []int `json:"ids"`      // array of supplier IDs
	IsDelete bool  `json:"isDelete"` // true = delete, false = restore
}
