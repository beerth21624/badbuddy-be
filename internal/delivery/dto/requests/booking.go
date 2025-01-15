package requests

// CreateBookingRequest represents the request to create a new court booking
type CreateBookingRequest struct {
	CourtID   string  `json:"court_id" validate:"required,uuid"`
	Date      string  `json:"date" validate:"required,datetime"`
	StartTime string  `json:"start_time" validate:"required,datetime"`
	EndTime   string  `json:"end_time" validate:"required,datetime"`
	Notes     *string `json:"notes" validate:"omitempty,min=1,max=500"`
}

// UpdateBookingRequest represents the request to update an existing booking
type UpdateBookingRequest struct {
	Status string  `json:"status" validate:"omitempty,oneof=confirmed cancelled"`
	Notes  *string `json:"notes" validate:"omitempty,min=1,max=500"`
}

// CreatePaymentRequest represents the request to create a payment for a booking
type CreatePaymentRequest struct {
	PaymentMethod string  `json:"payment_method" validate:"required,oneof=cash transfer card qr"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	TransactionID *string `json:"transaction_id" validate:"omitempty,min=1"`
}

//UpdatePaymentRequest represents the request to update a payment for a booking
type UpdatePaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"omitempty,oneof=cash transfer card qr"`
	Status        string `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

// ListBookingsRequest represents the request to list bookings with filters
type ListBookingsRequest struct {
	CourtID  string `json:"court_id" validate:"omitempty,uuid"`
	VenueID  string `json:"venue_id" validate:"omitempty,uuid"`
	DateFrom string `json:"date_from" validate:"omitempty,datetime"`
	DateTo   string `json:"date_to" validate:"omitempty,datetime"`
	Status   string `json:"status" validate:"omitempty,oneof=pending confirmed cancelled completed"`
	Limit    int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int    `json:"offset" validate:"omitempty,min=0"`
}

// CheckAvailabilityRequest represents the request to check court availability
type CheckAvailabilityRequest struct {
	CourtID   string `json:"court_id" validate:"required,uuid"`
	Date      string `json:"date" validate:"required,datetime"`
	StartTime string `json:"start_time" validate:"required,datetime"`
	EndTime   string `json:"end_time" validate:"required,datetime"`
}
