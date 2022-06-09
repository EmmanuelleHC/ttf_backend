package database

// Model base model
type ModelTrans struct {
	CreatedBy string `gorm:"column:created_by;not null;" json:"created_by"`
	UpdateBy  string `gorm:"column:update_by;not null;" json:"update_by"`
}
