package errors

var (
	CompanyRecordNotFound = New("company record not found")
	CompanyIsDisable      = New("company is disabled")
	CompanyAlreadyExists  = New("company already exists")
)
