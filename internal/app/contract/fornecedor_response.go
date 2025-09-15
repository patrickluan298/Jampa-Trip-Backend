package contract

// LoginFornecedorResponse - resposta de login para fornecedores
type LoginFornecedorResponse struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Token      string     `json:"token"`
	Dados      Fornecedor `json:"dados"`
}

// CadastrarFornecedorResponse - resposta de cadastro para fornecedores
type CadastrarFornecedorResponse struct {
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Dados      Fornecedor `json:"dados"`
}

// Fornecedor - dados do fornecedor logado
type Fornecedor struct {
	ID              int    `json:"id"`
	Nome            string `json:"nome"`
	Email           string `json:"email"`
	Senha           string `json:"senha"`
	CNPJ            string `json:"cnpj"`
	Telefone        string `json:"telefone"`
	Endereco        string `json:"endereco"`
	MomentoCadastro string `json:"momento_cadastro"`
}
