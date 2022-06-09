package errors

var (
	AccountRecordNotFound = New("Account record not found")
	AccountIsDisable      = New("Account is disabled")
	AccountAlreadyExists  = New("Account already exists")
)
