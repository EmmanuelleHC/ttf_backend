package errors

var (
	BKKDetailRecordNotFound = New("BKKDetail record not found")
	BKKDetailIsDisable      = New("BKKDetail is disabled")
	BKKDetailAlreadyExists  = New("BKKDetail already exists")
)
