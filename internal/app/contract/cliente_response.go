package contract

// LoginClienteResponse - resposta de login para clientes
type LoginClienteResponse struct {
	Mensagem string       `json:"mensagem"`
	Token    string       `json:"token"`
	Dados    ClienteLogin `json:"dados"`
}

// CadastrarClienteResponse - resposta de cadastro para clientes
type CadastrarClienteResponse struct {
	Mensagem string  `json:"mensagem"`
	Dados    Cliente `json:"dados"`
}

// AtualizarClienteResponse - resposta de atualização para clientes
type AtualizarClienteResponse struct {
	Mensagem string  `json:"mensagem"`
	Dados    Cliente `json:"dados"`
}

// ListarClienteResponse - resposta de listagem para clientes
type ListarClienteResponse struct {
	Mensagem string    `json:"mensagem"`
	Dados    []Cliente `json:"dados"`
}

// ClienteLogin - dados do cliente logado
type ClienteLogin struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Cliente - dados do cliente
type Cliente struct {
	ID                 int    `json:"id"`
	Nome               string `json:"nome"`
	Email              string `json:"email"`
	CPF                string `json:"cpf"`
	Telefone           string `json:"telefone"`
	DataNascimento     string `json:"data_nascimento"`
	MomentoCadastro    string `json:"momento_cadastro"`
	MomentoAtualizacao string `json:"momento_atualizacao"`
}
