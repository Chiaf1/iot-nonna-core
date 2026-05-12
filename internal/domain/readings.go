package domain

import "time"

type DhtReadings struct {
	Id          string    `json:"id"`
	DeviceId    string    `json:"device_id"`
	SensorId    string    `json:"sensor_id"`
	Timestamp   time.Time `json:"timestamp"`
	Temperature *float32  `json:"temperature,omitempty"`
	Humidity    *float32  `json:"humidity,omitempty"`
}

type StatusReadings struct {
	Id        string    `json:"id"`
	DeviceId  string    `json:"device_id"`
	SensorId  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    bool      `json:"status"`
}
