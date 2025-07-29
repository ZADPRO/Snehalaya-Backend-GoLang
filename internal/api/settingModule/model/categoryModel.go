package model

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
	DoorNo     string `json:"doorNumber"`
	Street     string `json:"streetName"`
	City       string `json:"city"`
	State      string `json:"state"`
	Role       string `json:"role" gorm:"column:role"`
	BranchName string `json:"branch" gorm:"column:branch"`
}
