package query

var (
	ObterClientePorID = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(nome, '') AS nome,
			COALESCE(email, '') AS email,
			COALESCE(senha, '') AS senha,
			COALESCE(cpf, '') AS cpf,
			COALESCE(telefone, '') AS telefone,
			COALESCE(data_nascimento, '') AS data_nascimento,
			momento_cadastro,
			momento_atualizacao 
		FROM clientes
		WHERE id = ?;
	`

	ObterClientePorEmail = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(nome, '') AS nome,
			COALESCE(email, '') AS email,
			COALESCE(senha, '') AS senha,
			COALESCE(cpf, '') AS cpf,
			COALESCE(telefone, '') AS telefone,
			COALESCE(data_nascimento, '') AS data_nascimento,
			momento_cadastro,
			momento_atualizacao 
		FROM clientes
		WHERE email = ?;
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
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(nome, '') AS nome,
			COALESCE(email, '') AS email,
			COALESCE(cpf, '') AS cpf,
			COALESCE(telefone, '') AS telefone,
			COALESCE(data_nascimento, '') AS data_nascimento,
			momento_cadastro,
			momento_atualizacao 
		FROM clientes
		WHERE CASE WHEN (? <> '') THEN nome = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN email = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN cpf = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN telefone = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN data_nascimento = ? ELSE TRUE END
		ORDER BY momento_cadastro DESC;
	`
)
