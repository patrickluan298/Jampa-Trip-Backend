package query

var (
	ObterEmpresaPorID = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(nome, '') AS nome,
			COALESCE(email, '') AS email,
			COALESCE(senha, '') AS senha,
			COALESCE(cnpj, '') AS cnpj,
			COALESCE(telefone, '') AS telefone,
			COALESCE(endereco, '') AS endereco,
			momento_cadastro,
			momento_atualizacao
		FROM empresas
		WHERE id = ?;
	`

	ObterEmpresaPorEmail = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(nome, '') AS nome,
			COALESCE(email, '') AS email,
			COALESCE(senha, '') AS senha,
			COALESCE(cnpj, '') AS cnpj,
			COALESCE(telefone, '') AS telefone,
			COALESCE(endereco, '') AS endereco,
			momento_cadastro,
			momento_atualizacao
		FROM empresas
		WHERE email = ?;
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
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(nome, '') AS nome,
			COALESCE(email, '') AS email,
			COALESCE(cnpj, '') AS cnpj,
			COALESCE(telefone, '') AS telefone,
			COALESCE(endereco, '') AS endereco,
			momento_cadastro,
			momento_atualizacao
		FROM empresas
		WHERE CASE WHEN (? <> '') THEN nome = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN email = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN cnpj = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN telefone = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN endereco = ? ELSE TRUE END
		ORDER BY momento_cadastro DESC;
	`
)
