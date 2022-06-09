package errors

var (
	EmployeeRecordNotFound = New("Employee record not found")
	EmployeeIsDisable      = New("Employee is disabled")
	EmployeeAlreadyExists  = New("Employee already exists")
)
