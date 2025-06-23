package model

type ReqVal struct {
	EncryptedData []string `json:"encryptedData"`
}

type CreateUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Designation string `json:"designation"`
	Status      bool   `json:"status"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RoleTypeId  int    `json:"roleTypeId"`
	BranchId    int    `json:"branchId"`
	CreatedBy   string `json:"createdBy"`
}

type User struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	RoleTypeId  int
	FirstName   string
	LastName    string
	Designation string
	Status      bool
	BranchId    int
	CreatedAt   string
	CreatedBy   string
	UpdatedAt   string
	UpdatedBy   string
}

type UserComm struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserId    uint
	Mobile    string
	Email     string
	Address   string
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy string
}

type UserAuth struct {
	ID             uint `gorm:"primaryKey;autoIncrement"`
	UserId         uint
	Username       string
	Password       string // original (should be optional)
	HashedPassword string
	CreatedAt      string
	CreatedBy      string
	UpdatedAt      string
	UpdatedBy      string
}
