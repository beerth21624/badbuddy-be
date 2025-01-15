package requests

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	PlayLevel string `json:"play_level" validate:"required"`
	Gender    string `json:"gender" validate:"required"`
	PlayHand  string `json:"play_hand" validate:"required"`
	Location  string `json:"location" validate:"required"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	PlayLevel string `json:"play_level"`
	Location  string `json:"location"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type SearchFilters struct {
	PlayLevel string `query:"play_level"`
	Location  string `query:"location"`
	Limit     int    `query:"limit" validate:"required,min=1,max=100"`
	Offset    int    `query:"offset" validate:"min=0"`
}

type UpdateRolesRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Role  string  `json:"role" validate:"required"`
}
