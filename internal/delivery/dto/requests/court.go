package requests

type CreateCourtRequest struct {
	VenueID      string  `json:"venue_id" validate:"required,uuid"`
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Description  string  `json:"description" validate:"omitempty,max=500"`
	PricePerHour float64 `json:"price_per_hour" validate:"required,gt=0"`
}

type UpdateCourtRequest struct {
	CourtID      string  `json:"court_id"`
	Name         string  `json:"name" validate:"omitempty,min=2,max=100"`
	Description  string  `json:"description" validate:"omitempty,max=500"`
	PricePerHour float64 `json:"price_per_hour" validate:"omitempty,gt=0"`
	Status       string  `json:"status" validate:"omitempty,oneof=available occupied maintenance"`
}

type UpdateCourtStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=available occupied maintenance"`
}

type ListCourtsRequest struct {
	VenueID  string  `json:"venue_id" validate:"omitempty,uuid"`
	Status   string  `json:"status" validate:"omitempty,oneof=available occupied maintenance"`
	Location string  `json:"location" validate:"omitempty,max=100"`
	PriceMin float64 `json:"price_min" validate:"omitempty,min=0"`
	PriceMax float64 `json:"price_max" validate:"omitempty,gtefield=PriceMin"`
	Limit    int     `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int     `json:"offset" validate:"omitempty,min=0"`
}

type CheckCourtAvailabilityRequest struct {
	CourtID   string `json:"court_id" validate:"required,uuid"`
	Date      string `json:"date" validate:"required,datetime=2006-01-02"`
	StartTime string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime   string `json:"end_time" validate:"required,datetime=15:04,gtfield=StartTime"`
}
