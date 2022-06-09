package errors

var (
	TarikDanaRecordNotFound = New("TarikDana record not found")
	TarikDanaIsDisable      = New("TarikDana is disabled")
	TarikDanaAlreadyExists  = New("TarikDana already exists")
)
