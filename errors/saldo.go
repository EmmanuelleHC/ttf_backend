package errors

var (
	SaldoRecordNotFound = New("Saldo record not found")
	SaldoIsDisable      = New("Saldo is disabled")
	SaldoAlreadyExists  = New("Saldo already exists")
	SaldoNeedsNew       = New("Saldo needs news")
)
