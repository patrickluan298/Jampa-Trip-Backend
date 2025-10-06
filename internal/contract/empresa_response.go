package contract

// CreateCompanyResponse - resposta de cadastro para empresas
type CreateCompanyResponse struct {
	Message string  `json:"message"`
	Data    Company `json:"data"`
}

// UpdateCompanyResponse - resposta de atualização para empresas
type UpdateCompanyResponse struct {
	Message string  `json:"message"`
	Data    Company `json:"data"`
}

// ListCompanyResponse - resposta de listagem para empresas
type ListCompanyResponse struct {
	Message string    `json:"message"`
	Data    []Company `json:"data"`
}

// GetCompanyResponse - resposta de obter empresa por ID
type GetCompanyResponse struct {
	Message string  `json:"message"`
	Data    Company `json:"data"`
}

// Company - dados do empresa logada
type Company struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CNPJ      string `json:"cnpj"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}