package posManagementModel

type AddCustomer struct {
	RefCustomerId       int    `gorm:"column:refCustomerId;primaryKey;autoIncrement" json:"refCustomerId"`
	RefCustomerName     string `gorm:"column:refCustomerName;type:varchar(255);not null" json:"refCustomerName" binding:"required"`
	RefMobileNo         string `gorm:"column:refMobileNo;type:varchar(20);not null;uniqueIndex" json:"refMobileNo" binding:"required"`
	RefAddress          string `gorm:"column:refAddress;type:varchar(255);not null" json:"refAddress" binding:"required"`
	RefCity             string `gorm:"column:refCity;type:varchar(100);not null" json:"refCity" binding:"required"`
	RefPincode          string `gorm:"column:refPincode;type:varchar(20)" json:"refPincode"`
	RefState            string `gorm:"column:refState;type:varchar(100)" json:"refState"`
	RefCountry          string `gorm:"column:refCountry;type:varchar(100);not null" json:"refCountry" binding:"required"`
	RefMembershipNumber string `gorm:"column:refMembershipNumber;type:varchar(100)" json:"refMembershipNumber"`
	RefTaxNumber        string `gorm:"column:refTaxNumber;type:varchar(50)" json:"refTaxNumber"`
	// RefCustId           int    `gorm:"column:refCustId" json:"refCustId"` // Not sure if used?
	CreatedAt string `gorm:"column:createdAt;type:timestamp" json:"-"`
	CreatedBy string `gorm:"column:createdBy;type:varchar(100)" json:"-"`
	UpdatedAt string `gorm:"column:updatedAt;type:timestamp" json:"-"`
	UpdatedBy string `gorm:"column:updatedBy;type:varchar(100)" json:"-"`
	IsDelete  bool   `gorm:"column:isDelete;default:false" json:"-"`
}
