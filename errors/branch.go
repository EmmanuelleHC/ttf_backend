package errors

var (
	BranchRecordNotFound = New("Branch record not found")
	BranchIsDisable      = New("Branch is disabled")
	BranchAlreadyExists  = New("Branch already exists")
)
