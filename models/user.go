package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable 0: Disable
type User struct {
	database.Model
	ID        string    `gorm:"column:id;size:36;index;not null;" json:"id"`
	Username  string    `gorm:"column:username;size:64;not null;index;" json:"username" validate:"required"`
	Realname  string    `gorm:"column:realname;size:64;not null;" json:"realname" validate:"required"`
	Password  string    `gorm:"column:password;not null;" json:"password"`
	Email     string    `gorm:"column:email;default:'';" json:"email"`
	Phone     string    `gorm:"column:phone;default:'';" json:"phone"`
	Status    int       `gorm:"column:status;not null;default:0;" json:"status" validate:"required,max=1,min=-1"`
	CreatedBy string    `gorm:"column:created_by;not null;" json:"created_by"`
	CompanyID string    `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID  string    `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	Company   Company   `gorm:"references:ID" json:"company" yaml:"company"`
	Branch    Branch    `gorm:"references:ID" json:"branch" yaml:"branch"`
	UserRoles UserRoles `gorm:"-" json:"user_roles"`
}

type Users []*User

type UserInfo struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Realname string `json:"realname"`
	Roles    Roles  `json:"roles"`
}

type UserQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	QueryPassword bool
	ID            string   `query:"user_id"`
	Username      string   `query:"username"`
	Realname      string   `query:"realname"`
	CompanyID     string   `query:"company_id"`
	BranchID      string   `query:"branch_id"`
	Email         string   `query:"email"`
	QueryValue    string   `query:"query_value"`
	Status        int      `query:"status" validate:"max=1,min=-1"`
	RoleIDs       []string `query:"-"`
}

type UserQueryResult struct {
	List       Users           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

func (a Users) ToIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.ID
	}
	return ids
}
