package phase

type Status string

// These are the valid statuses of APIMock.
const (
	Pending Status = "Pending"
	Running Status = "Running"
	Error   Status = "Error"
	Failed  Status = "Failed"
	Unknown Status = "Unknown"
)
