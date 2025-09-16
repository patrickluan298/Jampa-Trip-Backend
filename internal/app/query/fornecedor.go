package query

var (
	ObterPorEmail = `
		SELECT * FROM fornecedores WHERE email = ?;
	`

	Cadastrar = `
		INSERT INTO fornecedores (nome, email, senha, cnpj, telefone, endereco, momento_cadastro)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`
)
