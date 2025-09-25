package query

var (
	ObterEmpresaPorID = `
		SELECT id, nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM empresas WHERE id = ?;
	`

	ObterEmpresaPorEmail = `
		SELECT id, nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM empresas WHERE email = ?;
	`

	CadastrarEmpresa = `
		INSERT INTO empresas (nome, email, senha, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	AtualizarEmpresa = `
		UPDATE empresas 
		SET nome = ?, email = ?, senha = ?, cnpj = ?, telefone = ?, endereco = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	ListarTodasEmpresas = `
		SELECT id, nome, email, cnpj, telefone, endereco, momento_cadastro, momento_atualizacao 
		FROM empresas 
		ORDER BY momento_cadastro DESC;
	`
)
