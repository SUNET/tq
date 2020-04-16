package monitor

type Status struct {
	Mode  string   `json:"mode"`
	Peers []string `json:"peers"`
}

func NewStatus() *Status {
	status := &Status{}
	return status
}
