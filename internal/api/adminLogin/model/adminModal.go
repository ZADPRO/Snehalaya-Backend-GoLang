package model

// import "time"

type AdminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	User    *AdminLoginModelReq `json:"user,omitempty"`
	Token   string              `json:"token,omitempty"`
	Email   string              `json:"email"`
}

type AdminLoginModelReq struct {
	UserId            int    `json:"refUserId" gorm:"column:refUserId"`
	CustId            string `json:"refUserCustId" gorm:"column:refUserCustId"`
	RoleTypeId        string `json:"refRTId" gorm:"column:refRTId"`
	UserFName         string `json:"refUserFName" gorm:"column:refUserFName"`
	UserLName         string `json:"refUserLName" gorm:"column:refUserLName"`
	UserStatus        bool   `json:"refUserStatus" gorm:"column:refUserStatus"`
	UserBranchId      int    `json:"refUserBranchId" gorm:"refUserBranchId"`
	UACUsername       string `json:"refUACUsername" gorm:"column:refUACUsername"`
	UCDMobile         string `json:"refUCDMobile" gorm:"column:refUCDMobile"`
	UCDEmail          string `json:"refUCDEmail" gorm:"column:refUCDEmail"`
	UCDPassword       string `json:"refUACPassword" gorm:"column:refUACPassword"`
	UCDHashedPassword string `json:"refUACHashedPassword" gorm:"column:refUACHashedPassword"`
}

type OTPVerification struct {
	ID         int    `gorm:"primaryKey" json:"-"`
	Email      string `gorm:"column:email"`
	OTP        string `gorm:"column:otp"`
	ExpiresAt  string `gorm:"column:expires_at"` // stored as text in DB
	IsVerified bool   `gorm:"column:is_verified"`
	CreatedAt  string `gorm:"column:createdAt"`
	CreatedBy  string `gorm:"column:createdBy"`
	UpdatedAt  string `gorm:"column:updatedAt"`
	UpdatedBy  string `gorm:"column:updatedBy"`
}
