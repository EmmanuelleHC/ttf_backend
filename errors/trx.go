package errors

var (
	TrxRecordNotFound = New("Trx record not found")
	TrxIsDisable      = New("Trx is disabled")
	TrxAlreadyExists  = New("Trx already exists")
)
