package model

type (
	// Secret represents the variables attached to an specific app
	Secret struct {
		App  string            `json:"app"`
		Vars map[string]string `json:"vars"`
	}
)
