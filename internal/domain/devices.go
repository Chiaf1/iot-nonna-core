package domain

import "time"

type Device struct {
	Id         string     `json:"id"`
	Code       string     `json:"code"`
	DeviceType DeviceType `json:"device_type"`
	Room       *Room      `json:"room,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type DeviceRequest struct {
	Code         string  `json:"code" validate:"required,max=50"`
	DeviceTypeId string  `json:"device_type_id" validate:"required,uuid"`
	RoomId       *string `json:"room_id,omitempty" validate:"omitempty,uuid"`
}
