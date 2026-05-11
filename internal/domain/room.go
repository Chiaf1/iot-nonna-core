package domain

// Complete room struct
type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Struct to read the body of the reuqest
type RoomRequest struct {
	Name string `json:"name"`
}
