package errors

var (
	InvoiceHeaderRecordNotFound = New("InvoiceHeader record not found")
	InvoiceHeaderIsDisable      = New("InvoiceHeader is disabled")
	InvoiceHeaderAlreadyExists  = New("InvoiceHeader already exists")
)
