package query

var (
	GetClientByID = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(name, '') AS name,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(cpf, '') AS cpf,
			COALESCE(phone, '') AS phone,
			COALESCE(birth_date, '') AS birth_date,
			created_at,
			updated_at 
		FROM clients
		WHERE id = ?;
	`

	GetClientByEmail = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(name, '') AS name,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(cpf, '') AS cpf,
			COALESCE(phone, '') AS phone,
			COALESCE(birth_date, '') AS birth_date,
			created_at,
			updated_at 
		FROM clients
		WHERE email = ?;
	`

	CreateClient = `
		INSERT INTO clients (name, email, password, cpf, phone, birth_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	UpdateClient = `
		UPDATE clients 
		SET name = ?, email = ?, password = ?, cpf = ?, phone = ?, birth_date = ?, updated_at = ?
		WHERE id = ?;
	`

	ListAllClients = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(name, '') AS name,
			COALESCE(email, '') AS email,
			COALESCE(cpf, '') AS cpf,
			COALESCE(phone, '') AS phone,
			COALESCE(birth_date, '') AS birth_date,
			created_at,
			updated_at 
		FROM clients
		WHERE CASE WHEN (? <> '') THEN name = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN email = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN cpf = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN phone = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN birth_date = ? ELSE TRUE END
		ORDER BY created_at DESC;
	`
)
