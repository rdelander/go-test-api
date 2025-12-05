package address

// CreateAddressRequest represents the request to create an address
type CreateAddressRequest struct {
	EntityType  string `json:"entity_type" validate:"required,oneof=user"`
	EntityID    int32  `json:"entity_id" validate:"required,min=1"`
	AddressType string `json:"address_type" validate:"required,oneof=shipping billing"`
	StreetLine1 string `json:"street_line1" validate:"required,max=255"`
	StreetLine2 string `json:"street_line2" validate:"omitempty,max=255"`
	City        string `json:"city" validate:"required,max=100"`
	State       string `json:"state" validate:"required,max=100"`
	PostalCode  string `json:"postal_code" validate:"required,max=20"`
	Country     string `json:"country" validate:"required,max=100"`
}

// UpdateAddressRequest represents the request to update an address
type UpdateAddressRequest struct {
	StreetLine1 string `json:"street_line1" validate:"required,max=255"`
	StreetLine2 string `json:"street_line2" validate:"omitempty,max=255"`
	City        string `json:"city" validate:"required,max=100"`
	State       string `json:"state" validate:"required,max=100"`
	PostalCode  string `json:"postal_code" validate:"required,max=20"`
	Country     string `json:"country" validate:"required,max=100"`
}

// AddressResponse represents an address in API responses
type AddressResponse struct {
	ID          string `json:"id"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	AddressType string `json:"address_type"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}
