package model

import "encoding/json"

type NullableInt struct {
	IsSet bool
	Value *int
}

func (n *NullableInt) UnmarshalJSON(data []byte) error {
	n.IsSet = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}

	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.Value = &v
	return nil
}
