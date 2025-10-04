package query

var (
	GetCompanyByID = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(name, '') AS name,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(cnpj, '') AS cnpj,
			COALESCE(phone, '') AS phone,
			COALESCE(address, '') AS address,
			created_at,
			updated_at
		FROM companies
		WHERE id = ?;
	`

	GetCompanyByEmail = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(name, '') AS name,
			COALESCE(email, '') AS email,
			COALESCE(password, '') AS password,
			COALESCE(cnpj, '') AS cnpj,
			COALESCE(phone, '') AS phone,
			COALESCE(address, '') AS address,
			created_at,
			updated_at
		FROM companies
		WHERE email = ?;
	`

	CreateCompany = `
		INSERT INTO companies (name, email, password, cnpj, phone, address, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	UpdateCompany = `
		UPDATE companies 
		SET name = ?, email = ?, password = ?, cnpj = ?, phone = ?, address = ?, updated_at = ?
		WHERE id = ?;
	`

	ListAllCompanies = `
		SELECT
			COALESCE(id, 0) AS id,
			COALESCE(name, '') AS name,
			COALESCE(email, '') AS email,
			COALESCE(cnpj, '') AS cnpj,
			COALESCE(phone, '') AS phone,
			COALESCE(address, '') AS address,
			created_at,
			updated_at
		FROM companies
		WHERE CASE WHEN (? <> '') THEN name = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN email = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN cnpj = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN phone = ? ELSE TRUE END AND
			CASE WHEN (? <> '') THEN address = ? ELSE TRUE END
		ORDER BY created_at DESC;
	`
)
