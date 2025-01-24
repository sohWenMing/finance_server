package auth

type CreateUserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
