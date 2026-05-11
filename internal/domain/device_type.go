package domain

type Device_type struct {
	Id          string  `json:"id"`
	Code        string  `json:"code"`
	Topic       string  `json:"topic"`
	Description *string `json:"description,omitempty"`
}

type Device_typeRequest struct {
	Code        string  `json:"code" validate:"required"`
	Topic       string  `json:"topic" validate:"required,max=50"`
	Description *string `json:"description,omitempty"`
}
