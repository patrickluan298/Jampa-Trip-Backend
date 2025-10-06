package query

var (
	CreateTour = `
		INSERT INTO tours (company_id, name, dates, departure_time, arrival_time, max_people, description, images, price, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING id;
	`

	UpdateTour = `
		UPDATE tours 
		SET name = ?, dates = ?, departure_time = ?, arrival_time = ?, max_people = ?, description = ?, images = ?, price = ?, updated_at = NOW()
		WHERE id = ?;
	`

	GetTourByID = `
		SELECT
			COALESCE(t.id, 0) AS id,
			COALESCE(t.company_id, 0) AS company_id,
			COALESCE(t.name, '') AS name,
			t.dates,
			COALESCE(t.departure_time, '') AS departure_time,
			COALESCE(t.arrival_time, '') AS arrival_time,
			COALESCE(t.max_people, 0) AS max_people,
			COALESCE(t.description, '') AS description,
			t.images,
			COALESCE(t.price, 0) AS price,
			t.created_at,
			t.updated_at,
			COALESCE(c.name, '') AS company_name
		FROM tours t
		LEFT JOIN companies c ON t.company_id = c.id
		WHERE t.id = ?;
	`

	ListTours = `
		SELECT
			COALESCE(t.id, 0) AS id,
			COALESCE(t.company_id, 0) AS company_id,
			COALESCE(t.name, '') AS name,
			t.dates,
			COALESCE(t.departure_time, '') AS departure_time,
			COALESCE(t.arrival_time, '') AS arrival_time,
			COALESCE(t.max_people, 0) AS max_people,
			COALESCE(t.description, '') AS description,
			t.images,
			COALESCE(t.price, 0) AS price,
			t.created_at,
			t.updated_at,
			COALESCE(c.name, '') AS company_name
		FROM tours t
		LEFT JOIN companies c ON t.company_id = c.id
		WHERE CASE WHEN (? <> '') THEN t.name ILIKE ? ELSE TRUE END
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?;
	`

	ListMyTours = `
		SELECT
			COALESCE(t.id, 0) AS id,
			COALESCE(t.company_id, 0) AS company_id,
			COALESCE(t.name, '') AS name,
			t.dates,
			COALESCE(t.departure_time, '') AS departure_time,
			COALESCE(t.arrival_time, '') AS arrival_time,
			COALESCE(t.max_people, 0) AS max_people,
			COALESCE(t.description, '') AS description,
			t.images,
			COALESCE(t.price, 0) AS price,
			t.created_at,
			COALESCE(COUNT(r.id), 0) AS reservations_count
		FROM tours t
		LEFT JOIN reservas r ON t.id = r.tour_id
		WHERE t.company_id = ?
		GROUP BY t.id, t.company_id, t.name, t.dates, t.departure_time, t.arrival_time, t.max_people, t.description, t.images, t.price, t.created_at
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?;
	`

	DeleteTour = `
		DELETE FROM tours WHERE id = ?;
	`

	CountTours = `
		SELECT COUNT(*)
		FROM tours t
		WHERE CASE WHEN (? <> '') THEN t.name ILIKE ? ELSE TRUE END;
	`

	CountMyTours = `
		SELECT COUNT(*)
		FROM tours t
		WHERE t.company_id = ?;
	`

	CheckTourOwnership = `
		SELECT COUNT(*)
		FROM tours
		WHERE id = ? AND company_id = ?;
	`

	CountReservationsByTourID = `
		SELECT COUNT(*)
		FROM reservas
		WHERE tour_id = ? AND status IN ('pendente', 'confirmada');
	`
)
