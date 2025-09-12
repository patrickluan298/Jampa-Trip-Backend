package contract

// LoginClienteResponse - resposta de login para clientes
type LoginClienteResponse struct {
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Token      string          `json:"token"`
	Dados      ClienteResponse `json:"dados"`
}

// ClienteResponse - dados do cliente logado
type ClienteResponse struct {
	ID              int    `json:"id"`
	Nome            string `json:"nome"`
	Email           string `json:"email"`
	CPF             string `json:"cpf"`
	Telefone        string `json:"telefone"`
	DataNascimento  string `json:"data_nascimento"`
	MomentoCadastro string `json:"momento_cadastro"`
}
