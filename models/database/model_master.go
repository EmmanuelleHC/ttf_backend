package database

// Model base model
type ModelMaster struct {
	ActiveFlag   bool     `gorm:"column:active_flag;default:true;" json:"active_flag"`
	InactiveDate Datetime `gorm:"column:inactive_date;" json:"inactive_date"`
	CreatedBy    string   `gorm:"column:created_by;not null;" json:"created_by"`
	UpdateBy     string   `gorm:"column:update_by;not null;" json:"update_by"`
}
