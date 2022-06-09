package errors

var (
	CounterRecordNotFound = New("Counter record not found")
	CounterIsDisable      = New("Counter is disabled")
	CounterAlreadyExists  = New("Counter already exists")
)
