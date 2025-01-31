package data

import "time"

type Message struct {
	Time        time.Time `json:"time"`
	DeviceId    string    `json:"device_id"`
	DeviceType  string    `json:"device_type"`
	Temperature float64   `json:"temp"`
}
