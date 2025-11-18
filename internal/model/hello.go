package model

// HelloRequest represents the request body for the hello endpoint
type HelloRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// HelloResponse represents the response body for the hello endpoint
type HelloResponse struct {
	Message string `json:"message"`
}
