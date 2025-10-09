package query

var (
	// CreateFeedback - cria um novo feedback
	CreateFeedback = `
		INSERT INTO feedbacks (cliente_id, empresa_id, reserva_id, nota, comentario, status, momento_criacao, momento_atualizacao)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`

	// GetFeedbackByID - busca um feedback pelo ID com relacionamentos
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

	// GetFeedbacksByClienteID - busca feedbacks por cliente com paginação
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

	// CountFeedbacksByClienteID - conta total de feedbacks por cliente
	CountFeedbacksByClienteID = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE cliente_id = ?;
	`

	// GetFeedbacksByEmpresaID - busca feedbacks por empresa com paginação
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

	// CountFeedbacksByEmpresaID - conta total de feedbacks por empresa
	CountFeedbacksByEmpresaID = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE empresa_id = ?;
	`

	// GetFeedbacksByStatus - busca feedbacks por status com paginação
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

	// CountFeedbacksByStatus - conta total de feedbacks por status
	CountFeedbacksByStatus = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE status = ?;
	`

	// GetFeedbacksByRating - busca feedbacks por nota com paginação
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

	// CountFeedbacksByRating - conta total de feedbacks por nota
	CountFeedbacksByRating = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE nota = ?;
	`

	// UpdateFeedback - atualiza um feedback
	UpdateFeedback = `
		UPDATE feedbacks 
		SET nota = ?, comentario = ?, status = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	// UpdateFeedbackStatus - atualiza apenas o status de um feedback
	UpdateFeedbackStatus = `
		UPDATE feedbacks 
		SET status = ?, momento_atualizacao = ?
		WHERE id = ?;
	`

	// DeleteFeedback - remove um feedback
	DeleteFeedback = `
		DELETE FROM feedbacks WHERE id = ?;
	`

	// GetAverageRating - calcula a média de avaliações de uma empresa
	GetAverageRating = `
		SELECT 
			COALESCE(AVG(nota), 0) AS average, 
			COUNT(*) AS count
		FROM feedbacks
		WHERE empresa_id = ? AND status = ?;
	`

	// GetRatingDistribution - obtém a distribuição de notas de uma empresa
	GetRatingDistribution = `
		SELECT 
			nota, 
			COUNT(*) AS count
		FROM feedbacks
		WHERE empresa_id = ? AND status = ?
		GROUP BY nota
		ORDER BY nota;
	`

	// GetRecentFeedbacks - busca feedbacks recentes com paginação
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

	// CountRecentFeedbacks - conta total de feedbacks recentes
	CountRecentFeedbacks = `
		SELECT COUNT(*)
		FROM feedbacks
		WHERE empresa_id = ? AND momento_criacao >= ?;
	`
)
