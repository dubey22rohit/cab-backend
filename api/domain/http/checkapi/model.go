package checkapi

import "encoding/json"

type ready struct {
	Status string `json:"status"`
}

func (r ready) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}
