package contract

// Login ...
type Login struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// LoginResponse ...
type LoginResponse struct {
	ResponseJSON
	Token string `json:"token"`
}
