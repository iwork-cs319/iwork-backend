package mail

import (
	"time"
)

type EmailParams struct {
	Name          string
	Email         string
	WorkspaceName string
	FloorName     string
	Start         time.Time
	End           time.Time
}

type EmailClient interface {
	SendConfirmation(typeS string, params *EmailParams) error
	SendCancellation(typeS string, params *EmailParams) error
}
