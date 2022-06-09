package errors

var (
	DepartmentRecordNotFound = New("Department record not found")
	DepartmentIsDisable      = New("Department is disabled")
	DepartmentAlreadyExists  = New("Department already exists")
)
