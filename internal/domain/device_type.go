package domain

type DeviceType struct {
	Id          string  `json:"id"`
	Code        string  `json:"code"`
	Topic       string  `json:"topic"`
	Description *string `json:"description,omitempty"`
}

type DeviceTypeRequest struct {
	Code        string  `json:"code" validate:"required"`
	Topic       string  `json:"topic" validate:"required,max=50"`
	Description *string `json:"description,omitempty"`
}
