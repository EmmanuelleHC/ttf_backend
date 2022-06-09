package errors

var (
	CostCentreRecordNotFound = New("CostCentre record not found")
	CostCentreIsDisable      = New("CostCentre is disabled")
	CostCentreAlreadyExists  = New("CostCentre already exists")
)
