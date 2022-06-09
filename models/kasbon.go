package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Kasbon struct {
	database.Model
	database.ModelTrans
	ID          string            `gorm:"column:id;size:36;not null;index;" json:"id"`
	Num         string            `gorm:"column:num;size:10;not null;index:idx_num,unique;" json:"num"`
	Type        string            `gorm:"column:type;size:15;index;not null;" json:"type"`
	Amount      int64             `gorm:"column:amount;default:0;" json:"amount"`
	Description string            `gorm:"column:description;not null;" json:"description"`
	Date        database.Datetime `gorm:"column:date;" json:"date"`
	File        string            `gorm:"column:file;not null;" json:"file"`
	Status      string            `gorm:"column:status;size:15;index;not null;" json:"status"`
	ReleaseDate database.Datetime `gorm:"column:release_date;" json:"release_date"`
	PaideDate   database.Datetime `gorm:"column:paid_date;" json:"paid_date"`
	LunasDate   database.Datetime `gorm:"column:lunas_date;" json:"lunas_date"`
	EmployeeID  string            `gorm:"column:employee_id;size:36;index;not null;" json:"employee_id"`
	BKKHeaderID string            `gorm:"column:bkk_header_id;size:36;index;not null;" json:"bkk_header_id"`
	CompanyID   string            `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID    string            `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	DeptID      string            `gorm:"column:department_id;size:36;index;not null;" json:"department_id"`

	Company    Company    `gorm:"references:ID" json:"company" yaml:"company"`
	Branch     Branch     `gorm:"references:ID" json:"branch" yaml:"branch"`
	BKKHeader  BKKHeader  `gorm:"foreignKey:BKKHeaderID;references:ID" json:"bkkHeader" yaml:"bkkHeader"`
	Employee   Employee   `gorm:"references:ID" json:"employee" yaml:"employee"`
	Department Department `gorm:"foreignKey:DeptID;references:ID" json:"department" yaml:"department"`
}

type Kasbons []*Kasbon

type KasbonQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs         []string `query:"ids"`
	Num         string   `query:"num"`
	Type        string   `query:"type"`
	Amount      int64    `query:"amount"`
	Description string   `query:"description"`
	Date        string   `query:"date"`
	File        string   `query:"file"`
	Status      string   `query:"status"`
	ReleaseDate string   `query:"release_date"`
	PaideDate   string   `query:"paid_date"`
	LunasDate   string   `query:"lunas_date"`
	EmployeeID  string   `query:"employee_id"`
	BKKHeaderID string   `query:"bkk_header_id"`
	CompanyID   string   `query:"company_id"`
	BranchID    string   `query:"branch_id"`
	DeptID      string   `query:"department_id"`
	QueryValue  string   `query:"query_value"`
}

type KasbonQueryResult struct {
	List       Kasbons         `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Kasbons) ToMap() map[string]*Kasbon {
	m := make(map[string]*Kasbon)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
