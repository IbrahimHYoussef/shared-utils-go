package jsonutil

import (
	"encoding/json"
	"errors"
)

func EncodeMessage(message any) ([]byte, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("failed to encode message: " + err.Error())
	}
	return data, nil
}

func DecodeMessage[T any](raw []byte) (*T, error) {
	// Configure the Conversation and the first message
	var message T
	if err := json.Unmarshal(raw, &message); err != nil {
		return nil, errors.New("failed to decode message: " + err.Error())
	}
	return &message, nil
}
