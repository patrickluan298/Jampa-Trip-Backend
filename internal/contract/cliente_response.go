package contract

// CreateClientResponse - resposta de cadastro para clientes
type CreateClientResponse struct {
	Message string  `json:"message"`
	Data    Client  `json:"data"`
}

// UpdateClientResponse - resposta de atualização para clientes
type UpdateClientResponse struct {
	Message string  `json:"message"`
	Data   Client `json:"data"`
}

// ListClientResponse - resposta de listagem para clientes
type ListClientResponse struct {
	Message string   `json:"message"`
	Data    []Client `json:"data"`
}

// GetClientResponse - resposta de obter cliente por ID
type GetClientResponse struct {
	Message string  `json:"message"`
	Data    Client  `json:"data"`
}

// Client - dados do cliente
type Client struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CPF       string `json:"cpf"`
	Phone     string `json:"phone"`
	BirthDate string `json:"birth_date"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}