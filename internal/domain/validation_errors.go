package domain

type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}
