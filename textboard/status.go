package textboard

type Status string

const (
	Authenticating Status = "Authenticating"
	Authenticated  Status = "Authenticated"
	Disconnected   Status = "Disconnected"
	Reconnecting   Status = "Reconnecting"
)
