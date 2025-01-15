package responses

// CourtListResponse represents the response for listing courts
type CourtListResponse struct {
	Courts []CourtResponse `json:"courts"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
