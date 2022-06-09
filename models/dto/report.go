package dto

type ReportMonitoring struct {
	DateParams []string `json:"dateparam"`
	CompanyID  string   `json:"company_id"`
	BranchID   string   `json:"branch_id"`
	UserId     string
}
