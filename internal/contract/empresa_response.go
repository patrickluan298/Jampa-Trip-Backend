package contract

// LoginEmpresaResponse - resposta de login para empresas
type LoginEmpresaResponse struct {
	Mensagem string       `json:"mensagem"`
	Dados    EmpresaLogin `json:"dados"`
}

// CadastrarEmpresaResponse - resposta de cadastro para empresas
type CadastrarEmpresaResponse struct {
	Mensagem string  `json:"mensagem"`
	Dados    Empresa `json:"dados"`
}

// AtualizarEmpresaResponse - resposta de atualização para empresas
type AtualizarEmpresaResponse struct {
	Mensagem string  `json:"mensagem"`
	Dados    Empresa `json:"dados"`
}

// ListarEmpresasResponse - resposta de listagem para empresas
type ListarEmpresasResponse struct {
	Mensagem string    `json:"mensagem"`
	Dados    []Empresa `json:"dados"`
}

// ObterEmpresaResponse - resposta de obter empresa por ID
type ObterEmpresaResponse struct {
	Mensagem string  `json:"mensagem"`
	Dados    Empresa `json:"dados"`
}

// EmpresaLogin - dados do login da empresa
type EmpresaLogin struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Empresa - dados do empresa logada
type Empresa struct {
	ID                 int    `json:"id"`
	Nome               string `json:"nome"`
	Email              string `json:"email"`
	CNPJ               string `json:"cnpj"`
	Telefone           string `json:"telefone"`
	Endereco           string `json:"endereco"`
	MomentoCadastro    string `json:"momento_cadastro"`
	MomentoAtualizacao string `json:"momento_atualizacao"`
}
