package query

var (
	CreateFeedback = `
		INSERT INTO feedbacks (cliente_id, empresa_id, reserva_id, nota, comentario, status, momento_criacao, momento_atualizacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	GetFeedbackByID = `
		SELECT
			f.id,
			f.cliente_id,
			f.empresa_id,
			f.reserva_id,
			f.nota,
			COALESCE(f.comentario, '') AS comentario,
			f.status,
			f.momento_criacao,
			f.momento_atualizacao
		FROM feedbacks f
		WHERE f.id = ?;
	`

	GetFeedbacksByClienteID = `
		SELECT
			f.id,
			f.cliente_id,
			f.empresa_id,
			f.reserva_id,
			f.nota,
			COALESCE(f.comentario, '') AS comentario,
			f.status,
			f.momento_criacao,
			f.momento_atualizacao
		FROM feedbacks f
		WHERE f.cliente_id = ?
		ORDER BY f.momento_criacao DESC
		LIMIT ? OFFSET ?;
	`

	CountFeedbacksByClienteID = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE cliente_id = ?;
	`

	GetFeedbacksByEmpresaID = `
		SELECT
			f.id,
			f.cliente_id,
			f.empresa_id,
			f.reserva_id,
			f.nota,
			COALESCE(f.comentario, '') AS comentario,
			f.status,
			f.momento_criacao,
			f.momento_atualizacao
		FROM feedbacks f
		WHERE f.empresa_id = ?
		ORDER BY f.momento_criacao DESC
		LIMIT ? OFFSET ?;
	`

	CountFeedbacksByEmpresaID = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE empresa_id = ?;
	`

	GetFeedbacksByStatus = `
		SELECT
			f.id,
			f.cliente_id,
			f.empresa_id,
			f.reserva_id,
			f.nota,
			COALESCE(f.comentario, '') AS comentario,
			f.status,
			f.momento_criacao,
			f.momento_atualizacao
		FROM feedbacks f
		WHERE f.status = ?
		ORDER BY f.momento_criacao DESC
		LIMIT ? OFFSET ?;
	`

	CountFeedbacksByStatus = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE status = ?;
	`

	GetFeedbacksByRating = `
		SELECT
			f.id,
			f.cliente_id,
			f.empresa_id,
			f.reserva_id,
			f.nota,
			COALESCE(f.comentario, '') AS comentario,
			f.status,
			f.momento_criacao,
			f.momento_atualizacao
		FROM feedbacks f
		WHERE f.nota = ?
		ORDER BY f.momento_criacao DESC
		LIMIT ? OFFSET ?;
	`

	CountFeedbacksByRating = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE nota = ?;
	`

	UpdateFeedback = `
		UPDATE feedbacks 
		SET nota = ?, comentario = ?, status = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	UpdateFeedbackStatus = `
		UPDATE feedbacks 
		SET status = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	DeleteFeedback = `
		DELETE FROM feedbacks WHERE id = ?;
	`

	GetAverageRating = `
		SELECT 
			COALESCE(AVG(nota), 0) AS average, 
			COUNT(*) AS count
		FROM feedbacks
		WHERE empresa_id = ? AND status = ?;
	`

	GetRatingDistribution = `
		SELECT 
			nota, 
			COUNT(*) AS count
		FROM feedbacks
		WHERE empresa_id = ? AND status = ?
		GROUP BY nota
		ORDER BY nota;
	`

	GetRecentFeedbacks = `
		SELECT
			f.id,
			f.cliente_id,
			f.empresa_id,
			f.reserva_id,
			f.nota,
			COALESCE(f.comentario, '') AS comentario,
			f.status,
			f.momento_criacao,
			f.momento_atualizacao
		FROM feedbacks f
		WHERE f.empresa_id = ? AND f.momento_criacao >= ?
		ORDER BY f.momento_criacao DESC
		LIMIT ? OFFSET ?;
	`

	CountRecentFeedbacks = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE empresa_id = ? AND momento_criacao >= ?;
	`
)
