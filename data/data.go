package data

import (
	"fmt"
	"time"
)

type Message struct {
	Time        time.Time `json:"time"`
	DeviceId    string    `json:"device_id"`
	DeviceType  string    `json:"device_type"`
	Temperature float64   `json:"temp"`
}

func ValidateMessageBody(m Message) (bool, error) {
	if m.DeviceId == "" {
		return false, fmt.Errorf("invalid device id")
	} else if m.DeviceType == "" {
		return false, fmt.Errorf("invalid device id")
	}
	return true, nil
}

func ValidateDeviceId(devId string) (bool, error) {
	if devId == "" {
		return false, fmt.Errorf("invalid device id")
	}
	return true, nil
}
