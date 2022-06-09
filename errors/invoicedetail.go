package errors

var (
	InvoiceDetailRecordNotFound = New("InvoiceDetail record not found")
	InvoiceDetailIsDisable      = New("InvoiceDetail is disabled")
	InvoiceDetailAlreadyExists  = New("InvoiceDetail already exists")
)
