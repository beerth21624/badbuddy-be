package responses

import "time"

type UserResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	PlayLevel    string    `json:"play_level"`
	Location     string    `json:"location"`
	Bio          string    `json:"bio"`
	Gender       string    `json:"gender"`
	PlayHand     string    `json:"play_hand"`
	AvatarURL    string    `json:"avatar_url"`
	LastActiveAt time.Time `json:"last_active_at"`
	Role         string    `json:"role"`
	Venues       []Venue   `json:"venues"`
}

type UserProfileResponse struct {
	UserResponse
	HostedSessions  int     `json:"hosted_sessions"`
	JoinedSessions  int     `json:"joined_sessions"`
	AverageRating   float64 `json:"average_rating"`
	TotalReviews    int     `json:"total_reviews"`
	RegularPartners int     `json:"regular_partners"`
	Venues          []Venue `json:"venues"`
}

type Venue struct {
	ID string `json:"id"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}
