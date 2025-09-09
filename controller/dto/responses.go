package dto

// AuthResponse representa a resposta de autenticação
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse representa a resposta com dados do usuário
type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsBlocked bool   `json:"is_blocked"`
}
