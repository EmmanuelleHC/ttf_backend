package errors

var (
	BKKHeaderRecordNotFound = New("BKKHeader record not found")
	BKKHeaderIsDisable      = New("BKKHeader is disabled")
	BKKHeaderAlreadyExists  = New("BKKHeader already exists")
)
