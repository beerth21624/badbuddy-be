package responses

type FacilityListResponse struct {
	Facilities []FacilityResponse `json:"facilities"`
}

type FacilityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}