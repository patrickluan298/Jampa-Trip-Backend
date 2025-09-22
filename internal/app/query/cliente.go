package query

var (
	ObterClientePorID = `
		SELECT id, nome, email, senha, cpf, telefone, data_nascimento, momento_cadastro, momento_atualizacao 
		FROM clientes WHERE id = ?;
	`

	ObterClientePorEmail = `
		SELECT id, nome, email, senha, cpf, telefone, data_nascimento, momento_cadastro, momento_atualizacao 
		FROM clientes WHERE email = ?;
	`

	CadastrarCliente = `
		INSERT INTO clientes (nome, email, senha, cpf, telefone, data_nascimento, momento_cadastro, momento_atualizacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	AtualizarCliente = `
		UPDATE clientes 
		SET nome = ?, email = ?, senha = ?, cpf = ?, telefone = ?, data_nascimento = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	ListarTodosClientes = `
		SELECT id, nome, email, cpf, telefone, data_nascimento, momento_cadastro, momento_atualizacao 
		FROM clientes 
		ORDER BY momento_cadastro DESC;
	`
)
