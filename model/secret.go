package model

type (
	Secret struct {
		App  string            `json:"app"`
		Vars map[string]string `json:"vars"`
	}
)
