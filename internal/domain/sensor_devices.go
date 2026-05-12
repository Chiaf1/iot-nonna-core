package domain

type SensorDevices struct {
	DeviceId string     `json:"device_id"`
	Sensor   SensorType `json:"sensor"`
}

type AssociateSensorRequest struct {
	SensorId string `json:"sensor_id" validate:"required,uuid"`
}
