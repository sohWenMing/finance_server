package usermapping

type UserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreatedUserResponse struct {
	IsSuccess bool
	UserId    string
}

type LoginResponse struct {
	IsSuccess    bool   `json:"is_success"`
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
