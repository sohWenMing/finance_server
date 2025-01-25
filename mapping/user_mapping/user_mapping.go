package usermapping

type CreateUserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreatedUserResponse struct {
	IsSuccess bool
	UserId    string
}
