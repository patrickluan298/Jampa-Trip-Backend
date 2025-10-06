package contract

// LoginResponse - resposta de login
type LoginResponse struct {
	Message string        `json:"message"`
	Type    string        `json:"type"`
	Data    UserLoginData `json:"data"`
}

// UserLoginData - dados do usuário logado
type UserLoginData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
