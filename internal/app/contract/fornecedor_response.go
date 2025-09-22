package contract

// LoginFornecedorResponse - resposta de login para fornecedores
type LoginFornecedorResponse struct {
	Mensagem string          `json:"mensagem"`
	Token    string          `json:"token"`
	Dados    FornecedorLogin `json:"dados"`
}

// CadastrarFornecedorResponse - resposta de cadastro para fornecedores
type CadastrarFornecedorResponse struct {
	Mensagem string     `json:"mensagem"`
	Dados    Fornecedor `json:"dados"`
}

// AtualizarFornecedorResponse - resposta de atualização para fornecedores
type AtualizarFornecedorResponse struct {
	Mensagem string     `json:"mensagem"`
	Dados    Fornecedor `json:"dados"`
}

// ListarFornecedorResponse - resposta de listagem para fornecedores
type ListarFornecedorResponse struct {
	Mensagem string       `json:"mensagem"`
	Dados    []Fornecedor `json:"dados"`
}

// FornecedorLogin - dados do fornecedor logado
type FornecedorLogin struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Fornecedor - dados do fornecedor logado
type Fornecedor struct {
	ID                 int    `json:"id"`
	Nome               string `json:"nome"`
	Email              string `json:"email"`
	CNPJ               string `json:"cnpj"`
	Telefone           string `json:"telefone"`
	Endereco           string `json:"endereco"`
	MomentoCadastro    string `json:"momento_cadastro"`
	MomentoAtualizacao string `json:"momento_atualizacao"`
}
