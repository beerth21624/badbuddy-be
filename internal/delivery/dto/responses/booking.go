package responses

// BookingResponse represents the response for a court booking
type BookingListResponse struct {
	Bookings []BookingResponse `json:"bookings"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type BookingResponse struct {
	ID            string           `json:"id"`
	CourtName     string           `json:"court_name"`
	VenueName     string           `json:"venue_name"`
	VenueLocation string           `json:"venue_location"`
	UserName      string           `json:"user_name"`
	Date          string           `json:"date"`
	StartTime     string           `json:"start_time"`
	EndTime       string           `json:"end_time"`
	Duration      string           `json:"duration"`
	TotalAmount   float64          `json:"total_amount"`
	Status        string           `json:"status"`
	Notes         string           `json:"notes,omitempty"`
	CreatedAt     string           `json:"created_at"`
	UpdatedAt     string           `json:"updated_at"`
	CancelledAt   string           `json:"cancelled_at,omitempty"`
	Payment       *PaymentResponse `json:"payment,omitempty"`
}

// PaymentResponse represents the response for a booking payment
type PaymentResponse struct {
	ID            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TransactionID string  `json:"transaction_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// CourtAvailabilityResponse represents the response for court availability check
type CourtAvailabilityResponse struct {
	CourtID   string        `json:"court_id"`
	CourtName string        `json:"court_name"`
	Date      string        `json:"date"`
	Available bool          `json:"available"`
	TimeSlots []TimeSlot    `json:"time_slots"`
	Conflicts []BookingSlot `json:"conflicts,omitempty"`
}

// TimeSlot represents an available time slot
type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// BookingSlot represents a conflicting booking slot
type BookingSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}
