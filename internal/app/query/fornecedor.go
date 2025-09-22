package query

var (
	ObterFornecedorPorID = `
		SELECT id, nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM fornecedores WHERE id = ?;
	`

	ObterFornecedorPorEmail = `
		SELECT id, nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM fornecedores WHERE email = ?;
	`

	CadastrarFornecedor = `
		INSERT INTO fornecedores (nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	AtualizarFornecedor = `
		UPDATE fornecedores 
		SET nome = ?, email = ?, senha = ?, cnpj = ?, telefone = ?, endereco = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	ListarTodosFornecedores = `
		SELECT id, nome, email, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM fornecedores 
		ORDER BY momento_cadastro DESC;
	`
)
