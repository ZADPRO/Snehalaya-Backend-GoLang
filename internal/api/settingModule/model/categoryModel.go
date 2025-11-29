package model

type TransactionHistory struct {
	RefTransHisId   int    `json:"refTransHisId" gorm:"column:refTransHisId;primaryKey;autoIncrement"`
	RefTransTypeId  int    `json:"refTransTypeId" gorm:"column:refTransTypeId"`
	RefTransHisData string `json:"refTransHisData" gorm:"column:refTransHisData"`
	CreatedAt       string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy       string `json:"createdBy" gorm:"column:createdBy"`
	RefUserId       int    `json:"refUserId" gorm:"column:refUserId"`
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

type Category struct {
	RefCategoryId int    `json:"refCategoryId" gorm:"column:refCategoryid;primaryKey;autoIncrement"`
	CategoryName  string `json:"categoryName" gorm:"column:categoryName"`
	CategoryCode  string `json:"categoryCode" gorm:"column:categoryCode"`
	IsActive      bool   `json:"isActive" gorm:"column:isActive"`
	IsDelete      bool   `json:"isDelete" gorm:"column:isDelete"`
	CreatedAt     string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy     string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt     string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy     string `json:"updatedBy" gorm:"column:updatedBy"`
	ProfitMargin  string `json:"profitMargin" gorm:"column:profitMargin"`
}

type SubCategory struct {
	RefSubCategoryId int    `json:"refSubCategoryId" gorm:"column:refSubCategoryId;primaryKey;autoIncrement"`
	SubCategoryName  string `json:"subCategoryName" gorm:"column:subCategoryName"`
	RefCategoryId    int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	SubCategoryCode  string `json:"subCategoryCode" gorm:"column:subCategoryCode"`
	IsActive         bool   `json:"isActive" gorm:"column:isActive"`
	CreatedAt        string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy        string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt        string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy        string `json:"updatedBy" gorm:"column:updatedBY"` // 🔧 Fix is here
	IsDelete         bool   `json:"isDelete" gorm:"column:isDelete"`
}

type SubCategoryResponse struct {
	RefSubCategoryId int    `json:"refSubCategoryId" gorm:"column:refSubCategoryId"`
	SubCategoryName  string `json:"subCategoryName" gorm:"column:subCategoryName"`
	RefCategoryId    int    `json:"refCategoryId" gorm:"column:refCategoryId"`
	CategoryName     string `json:"categoryName" gorm:"column:categoryName"`
	SubCategoryCode  string `json:"subCategoryCode" gorm:"column:subCategoryCode"`
	IsActive         bool   `json:"isActive" gorm:"column:isActive"`
	CreatedAt        string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy        string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt        string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy        string `json:"updatedBy" gorm:"column:updatedBY"`
}

type Branch struct {
	RefBranchId   int    `gorm:"column:refBranchId;primaryKey;autoIncrement" json:"refBranchId"`
	RefBranchName string `gorm:"column:refBranchName" json:"refBranchName"`
	RefBranchCode string `gorm:"column:refBranchCode" json:"refBranchCode"`
	RefLocation   string `gorm:"column:refLocation" json:"refLocation"`
	RefMobile     string `gorm:"column:refMobile" json:"refMobile"`
	RefEmail      string `gorm:"column:refEmail" json:"refEmail"`
	IsMainBranch  bool   `gorm:"column:isMainBranch" json:"isMainBranch"`
	IsActive      bool   `gorm:"column:isActive" json:"isActive"`
	RefBTId       int    `gorm:"column:refBTId" json:"refBTId"`
	CreatedAt     string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy     string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt     string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy     string `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete      bool   `gorm:"column:isDelete" json:"isDelete"`
}

type BranchWithFloor struct {
	RefBranchId      int    `gorm:"column:refBranchId;primaryKey;autoIncrement"`
	RefBranchName    string `gorm:"column:refBranchName" json:"refBranchName"`
	RefBranchCode    string `gorm:"column:refBranchCode" json:"refBranchCode"`
	RefLocation      string `gorm:"column:refLocation" json:"refLocation"`
	RefMobile        string `gorm:"column:refMobile" json:"refMobile"`
	RefEmail         string `gorm:"column:refEmail" json:"refEmail"`
	IsMainBranch     bool   `gorm:"column:isMainBranch" json:"isMainBranch"`
	IsActive         bool   `gorm:"column:isActive" json:"isActive"`
	RefBTId          int    `gorm:"column:refBTId" json:"refBTId"`
	CreatedAt        string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy        string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt        string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete         bool   `gorm:"column:isDelete" json:"isDelete"`
	IsOnline         bool   `gorm:"column:isOnline" json:"isOnline"`
	IsOffline        bool   `gorm:"column:isOffline" json:"isOffline"`
	RefBranchDoorNo  string `gorm:"column:refBranchDoorNo" json:"refBranchDoorNo"`
	RefBranchStreet  string `gorm:"column:refBranchStreet" json:"refBranchStreet"`
	RefBranchCity    string `gorm:"column:refBranchCity" json:"refBranchCity"`
	RefBranchState   string `gorm:"column:refBranchState" json:"refBranchState"`
	RefBranchPincode string `gorm:"column:refBranchPincode" json:"refBranchPincode"`
}

type Floors struct {
	RefFloorId   int    `gorm:"column:refFloorId;primaryKey;autoIncrement" json:"refFloorId"`
	RefBranchId  int    `gorm:"column:refBranchId" json:"refBranchId"`
	RefFloorName string `gorm:"column:refFloorName" json:"refFloorName"`
	RefFloorCode string `gorm:"column:refFloorCode" json:"refFloorCode"`
	IsActive     string `gorm:"column:isActive" json:"isActive"`
	CreatedAt    string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy    string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt    string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy    string `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete     bool   `gorm:"column:isDelete" json:"isDelete"`
}

type Sections struct {
	RefSectionId     int    `gorm:"column:refSectionId;primaryKey;autoIncrement" json:"refSectionId"`
	RefFloorId       int    `gorm:"column:refFloorId" json:"refFloorId"`
	RefSectionName   string `gorm:"column:refSectionName" json:"refSectionName"`
	RefSectionCode   string `gorm:"column:refSectionCode" json:"refSectionCode"`
	RefCategoryId    int    `gorm:"column:refCategoryId" json:"refCategoryId"`
	RefSubCategoryId int    `gorm:"column:refSubCategoryId" json:"refSubCategoryId"`
	IsActive         string `gorm:"column:isActive" json:"isActive"`
	CreatedAt        string `gorm:"column:createdAt" json:"createdAt"`
	CreatedBy        string `gorm:"column:createdBy" json:"createdBy"`
	UpdatedAt        string `gorm:"column:updatedAt" json:"updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy" json:"updatedBy"`
	IsDelete         bool   `gorm:"column:isDelete" json:"isDelete"`
}

type BranchResponse struct {
	RefBranchId      int             `json:"refBranchId"`
	RefBranchName    string          `json:"refBranchName"`
	RefBranchCode    string          `json:"refBranchCode"`
	RefLocation      string          `json:"refLocation"`
	RefMobile        string          `json:"refMobile"`
	RefEmail         string          `json:"refEmail"`
	IsMainBranch     bool            `json:"isMainBranch"`
	IsActive         bool            `json:"isActive"`
	IsOnline         bool            `json:"isOnline"`
	IsOffline        bool            `json:"isOffline"`
	RefBranchDoorNo  string          `gorm:"column:refBranchDoorNo" json:"refBranchDoorNo"`
	RefBranchStreet  string          `gorm:"column:refBranchStreet" json:"refBranchStreet"`
	RefBranchCity    string          `gorm:"column:refBranchCity" json:"refBranchCity"`
	RefBranchState   string          `gorm:"column:refBranchState" json:"refBranchState"`
	RefBranchPincode string          `gorm:"column:refBranchPincode" json:"refBranchPincode"`
	Floors           []FloorResponse `json:"floors"`
}

type FloorResponse struct {
	RefFloorId int               `json:"refFloorId"`
	FloorName  string            `json:"floorName"`
	FloorCode  string            `json:"floorCode"`
	Sections   []SectionResponse `json:"sections"`
}

type SectionResponse struct {
	RefSectionId     int    `json:"refSectionId"`
	SectionName      string `json:"sectionName"`
	SectionCode      string `json:"sectionCode"`
	CategoryId       int    `json:"categoryId"`
	RefSubCategoryId int    `json:"refSubCategoryId"`
}

type RoleType struct {
	RefRTId   int    `json:"refRTId" gorm:"column:refRTId"`
	RefRTName string `json:"refRTName" gorm:"column:refRTName"`
}

func (RoleType) TableName() string {
	return "RoleType"
}

type User struct {
	RefUserId          int    `gorm:"primaryKey;autoIncrement;column:refUserId"`
	RefUserCustId      string `gorm:"column:refUserCustId"`
	RefRTId            int    `gorm:"column:refRTId"`
	RefUserFName       string `gorm:"column:refUserFName"`
	RefUserLName       string `gorm:"column:refUserLName"`
	RefUserDesignation string `gorm:"column:refUserDesignation"`
	RefUserStatus      string `gorm:"column:refUserStatus"`
	RefUserBranchId    int    `gorm:"column:refUserBranchId"`
	CreatedAt          string `gorm:"column:createdAt"`
	CreatedBy          string `gorm:"column:createdBy"`
	UpdatedAt          string `gorm:"column:updatedAt"`
	UpdatedBy          string `gorm:"column:updatedBy"`
	IsDelete           bool   `json:"isDelete" gorm:"column:isDelete"`
}

type UserAuth struct {
	RefUACId             int    `gorm:"primaryKey;autoIncrement;column:refUACId"`
	RefUserId            int    `gorm:"column:refUserId"`
	RefUACPassword       string `gorm:"column:refUACPassword"`
	RefUACHashedPassword string `gorm:"column:refUACHashedPassword"`
	RefUACUsername       string `gorm:"column:refUACUsername"`
	CreatedAt            string `gorm:"column:createdAt"`
	CreatedBy            string `gorm:"column:createdBy"`
	UpdatedAt            string `gorm:"column:updatedAt"`
	UpdatedBy            string `gorm:"column:updatedBy"`
}

type UserCommunication struct {
	RefUserComDetId int    `gorm:"primaryKey;autoIncrement;column:refUserComDetId"`
	RefUserId       int    `gorm:"column:refUserId"`
	RefUCDMobile    string `gorm:"column:refUCDMobile"`
	RefUCDEmail     string `gorm:"column:refUCDEmail"`
	RefUCDDoorNo    string `gorm:"column:refUCDDoorNo"`
	RefUCDStreet    string `gorm:"column:refUCDStreet"`
	RefUCDCity      string `gorm:"column:refUCDCity"`
	RefUCDState     string `gorm:"column:refUCDState"`
	CreatedAt       string `gorm:"column:createdAt"`
	CreatedBy       string `gorm:"column:createdBy"`
	UpdatedAt       string `gorm:"column:updatedAt"`
	UpdatedBy       string `gorm:"column:updatedBy"`
}

type EmployeePayload struct {
	RefUserId       int    `json:"refUserId"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Designation     string `json:"designation"`
	RoleTypeId      int    `json:"roleTypeId"`
	RefUserStatus   bool   `json:"refUserStatus"`
	RefUserBranchId int    `json:"refUserBranchId" gorm:"column:refUserBranchId"`
	Username        string `json:"username"`
	Mobile          string `json:"mobile"`
	Email           string `json:"email"`
	DoorNumber      string `json:"doorNumber"`
	StreetName      string `json:"streetName"`
	City            string `json:"city"`
	State           string `json:"state"`
}

type ProfilePayload struct {
	RefUserId       int    `json:"refUserId"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Designation     string `json:"designation"`
	RoleTypeId      int    `json:"roleTypeId"`
	RefUserBranchId int    `gorm:"column:refUserBranchId"`
	Username        string `json:"username"`
	Mobile          string `json:"mobile"`
	Email           string `json:"email"`
	DoorNumber      string `json:"doorNumber"`
	StreetName      string `json:"streetName"`
	City            string `json:"city"`
	State           string `json:"state"`
}

type EmployeeResponse struct {
	User
	Username   string `json:"username"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Doorno     string `json:"doorNumber"`
	Street     string `json:"streetName"`
	City       string `json:"city"`
	State      string `json:"state"`
	Role       string `json:"role" gorm:"column:role"`
	BranchName string `json:"branch" gorm:"column:branch"`
}

type AttributeGroupTable struct {
	AttributeGroupId   int    `gorm:"primaryKey;autoIncrement;column:attributeGroupId"`
	AttributeGroupName string `gorm:"column:attributeGroupName"`
}

type AttributesTable struct {
	AttributeId      int    `gorm:"primaryKey;autoIncrement;column:attributeId"`
	AttributeGroupId int    `gorm:"column:attributeGroupId"`
	AttributeKey     string `gorm:"column:attributeKey"`
	AttributeValue   string `gorm:"column:attributeValue"`
	CreatedAt        string `gorm:"column:createdAt"`
	CreatedBy        string `gorm:"column:createdBy"`
	UpdatedAt        string `gorm:"column:updatedAt"`
	UpdatedBy        string `gorm:"column:updatedBy"`
	IsDelete         bool   `gorm:"column:isDelete" json:"isDelete"`
}

type AttributeWithGroup struct {
	AttributesTable
	AttributeGroupName string `gorm:"column:attributeGroupName" json:"attributeGroupName"`
}

type OverviewCards struct {
	Categories int `json:"Categories" gorm:"column:Categories"`
	Branches   int `json:"Branches" gorm:"column:Branches"`
	Supplier   int `json:"Supplier" gorm:"column:Supplier"`
	Attributes int `json:"Attributes" gorm:"column:Attributes"`
}

type LatestSupplier struct {
	RefSupplierId   int    `json:"refSupplierId" gorm:"column:supplierId"`
	RefSupplierName string `json:"refSupplierName" gorm:"column:supplierName"`
	CreatedAt       string `json:"createdAt" gorm:"column:createdAt"`
}

type LatestCategory struct {
	RefCategoryId   int    `json:"refCategoryId" gorm:"column:refCategoryid"`
	RefCategoryName string `json:"refCategoryName" gorm:"column:categoryName"`
	CreatedAt       string `json:"createdAt" gorm:"column:createdAt"`
}

type CategorySubCategoryCount struct {
	Month         string `json:"month" gorm:"column:month" `
	Categories    int    `json:"Categories" gorm:"column:Categories"`
	SubCategories int    `json:"SubCategories" gorm:"column:SubCategories"`
}

type SettingsOverview struct {
	Cards            OverviewCards              `json:"cards" `
	LatestSuppliers  []LatestSupplier           `json:"latestSuppliers"`
	LatestCategories []LatestCategory           `json:"latestCategories"`
	MonthlyCounts    []CategorySubCategoryCount `json:"monthlyCounts"`
}

type ProductFieldDefinition struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	ColumnName  *string `json:"column_name" gorm:"column:column_name"`
	ColumnLabel string  `json:"column_label" gorm:"column:column_label"`
	DataType    string  `json:"data_type" gorm:"column:data_type"`
	Type        string  `json:"type" gorm:"column:type"`
	IsRequired  bool    `json:"is_required" gorm:"column:is_required"`
	CreatedAt   string  `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy   string  `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt   string  `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy   string  `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete    bool    `json:"isDelete" gorm:"column:isDelete"`
}

func (ProductFieldDefinition) TableName() string {
	return `"purchaseOrder".product_field_definitions`
}

type SettingsProduct struct {
	Id            int    `json:"id" gorm:"column:id;primaryKey"`
	CategoryId    int    `json:"categoryId" gorm:"column:categoryId"`
	SubCategoryId int    `json:"subCategoryId" gorm:"column:subCategoryId"`
	ProductName   string `json:"productName" gorm:"column:productName"`
	HsnCode       string `json:"hsn" gorm:"column:hsnCode"`
	TaxPercentage string `json:"tax" gorm:"column:taxPercentage"`
	ProductCode   string `json:"productCode" gorm:"column:productCode"`
	CreatedAt     string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy     string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt     string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy     string `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete      bool   `json:"isDelete" gorm:"column:isDelete"`
}

// Bind struct to table name
func (SettingsProduct) TableName() string {
	return `"SettingsProducts"`
}

type SettingsProductResponse struct {
	Id            int    `json:"id" gorm:"column:id"`
	CategoryId    int    `json:"categoryId" gorm:"column:categoryId"`
	SubCategoryId int    `json:"subCategoryId" gorm:"column:subCategoryId"`
	ProductName   string `json:"productName" gorm:"column:productName"`
	HsnCode       string `json:"hsnCode" gorm:"column:hsnCode"`
	TaxPercentage string `json:"taxPercentage" gorm:"column:taxPercentage"`
	ProductCode   string `json:"productCode" gorm:"column:productCode"`
	CreatedAt     string `json:"createdAt" gorm:"column:createdAt"`
	CreatedBy     string `json:"createdBy" gorm:"column:createdBy"`
	UpdatedAt     string `json:"updatedAt" gorm:"column:updatedAt"`
	UpdatedBy     string `json:"updatedBy" gorm:"column:updatedBy"`
	IsDelete      bool   `json:"isDelete" gorm:"column:isDelete"`

	// JOIN FIELDS
	CategoryName    string `json:"categoryName" gorm:"column:categoryName"`
	SubCategoryName string `json:"subCategoryName" gorm:"column:subCategoryName"`
}

type MasterPayload struct {
	Name string `json:"name" binding:"required"`
}

type MasterUpdatePayload struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
