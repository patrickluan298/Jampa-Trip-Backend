package query

var (
	ObterPorEmail = `
		SELECT id, nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM fornecedores WHERE email = ?;
	`

	Cadastrar = `
		INSERT INTO fornecedores (nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	Atualizar = `
		UPDATE fornecedores 
		SET nome = ?, email = ?, senha = ?, cnpj = ?, telefone = ?, endereco = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	ObterPorID = `
		SELECT id, nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM fornecedores WHERE id = ?;
	`
)
