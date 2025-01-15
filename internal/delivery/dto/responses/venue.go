package responses

import "time"

type CourtResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	PricePerHour float64 `json:"price_per_hour"`
	Status       string  `json:"status"`
}

type VenueResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Address      string              `json:"address"`
	Location     string              `json:"location"`
	Phone        string              `json:"phone"`
	Email        string              `json:"email"`
	OpenRange    []OpenRangeResponse `json:"open_range" validate:"required"`
	ImageURLs    string              `json:"image_urls"`
	Status       string              `json:"status"`
	Rating       float64             `json:"rating"`
	TotalReviews int                 `json:"total_reviews"`
	Courts       []CourtResponse     `json:"courts"`
	Facilities   []FacilityResponse  `json:"facilities"`
	Rules        []RuleResponse      `json:"rules"`
	Latitude     float64             `json:"latitude"`
	Longitude    float64             `json:"longitude"`
}

type OpenRangeResponse struct {
	Day       string    `json:"day"`
	IsOpen    bool      `json:"is_open"`
	OpenTime  time.Time `json:"open_time"`
	CloseTime time.Time `json:"close_time"`
}

type RuleResponse struct {
	Rule string `json:"rule"`
}

type VenueResponseDTO struct {
	Venues []VenueResponse `json:"venues"`
	Total  int             `json:"total"`
}

type ListVenueResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ReviewResponse struct {
	ID        string           `json:"id"`
	Rating    int              `json:"rating"`
	Comment   string           `json:"comment"`
	CreatedAt string           `json:"created_at"`
	Reviewer  ReviewerResponse `json:"reviewer"`
}

type ReviewerResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}
