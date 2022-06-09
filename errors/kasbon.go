package errors

var (
	KasbonRecordNotFound = New("Kasbon record not found")
	KasbonIsDisable      = New("Kasbon is disabled")
	KasbonAlreadyExists  = New("Kasbon already exists")
)
