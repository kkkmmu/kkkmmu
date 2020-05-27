package cproxy

import (
	"encoding/json"
	"errors"
)

type CControl struct {
	Type  string
	Value string
}

func EncodeCControl(typ, value string) ([]byte, error) {
	return json.Marshal(&CControl{
		Type:  typ,
		Value: value,
	})
}

func DecodeCControl(data []byte) (*CControl, error) {
	var cm CControl
	err := json.Unmarshal(data, &cm)
	if err != nil {
		return nil, errors.New("Cannot parse control message")
	}

	return &cm, nil
}
