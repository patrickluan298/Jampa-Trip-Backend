package query

var (
	CreateImage = `
        INSERT INTO images (
            user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id, uploaded_at, updated_at
    `

	GetImageByID = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE id = $1
    `

	GetImageByIDAndUser = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE id = $1 AND user_id = $2
    `

	UpdateImage = `
        UPDATE images 
        SET 
            tour_id = $2,
            description = $3,
            alt_text = $4,
            is_primary = $5,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING updated_at
    `

	DeleteImage = `
        DELETE FROM images WHERE id = $1
    `

	ListImages = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE user_id = $1
            AND ($2::int IS NULL OR tour_id = $2)
            AND ($3::text IS NULL OR format = $3)
        ORDER BY 
            CASE WHEN $4 = 'uploaded_at' THEN uploaded_at END DESC,
            CASE WHEN $4 = 'size' THEN size END DESC,
            CASE WHEN $4 = 'filename' THEN filename END ASC,
            sort_order ASC, id ASC
        LIMIT $5 OFFSET $6
    `

	CountImages = `
        SELECT COUNT(*)
        FROM images 
        WHERE user_id = $1
            AND ($2::int IS NULL OR tour_id = $2)
            AND ($3::text IS NULL OR format = $3)
    `

	ListImagesByTour = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE tour_id = $1 AND user_id = $2
        ORDER BY sort_order ASC, id ASC
    `

	IsImageOwnedByUser = `
        SELECT COUNT(*) > 0
        FROM images 
        WHERE id = $1 AND user_id = $2
    `

	IsImageUsedInActiveTour = `
        SELECT COUNT(*) > 0
        FROM images i
        INNER JOIN tours t ON i.tour_id = t.id
        WHERE i.id = $1 AND t.id IS NOT NULL
    `

	GetImageUsage = `
        SELECT 
            t.name as tour_name,
            CASE WHEN t.id IS NOT NULL THEN true ELSE false END as is_used
        FROM images i
        LEFT JOIN tours t ON i.tour_id = t.id
        WHERE i.id = $1
    `

	UpdateImageSortOrder = `
        UPDATE images 
        SET sort_order = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND user_id = $3
    `

	BatchUpdateSortOrder = `
        UPDATE images 
        SET sort_order = CASE id
            $1
        END, updated_at = CURRENT_TIMESTAMP
        WHERE id = ANY($2::int[]) AND user_id = $3
    `

	RemovePrimaryFromTour = `
        UPDATE images 
        SET is_primary = false, updated_at = CURRENT_TIMESTAMP
        WHERE tour_id = $1 AND user_id = $2 AND id != $3
    `

	SetImageAsPrimary = `
        UPDATE images 
        SET is_primary = true, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1 AND user_id = $2
    `

	GetImageStats = `
        SELECT 
            COUNT(*) as total_images,
            COALESCE(SUM(size), 0) as total_size,
            COALESCE(AVG(size), 0) as average_size,
            COUNT(CASE WHEN is_primary = true THEN 1 END) as primary_images,
            COUNT(CASE WHEN tour_id IS NULL THEN 1 END) as unused_images,
            COUNT(CASE WHEN uploaded_at >= NOW() - INTERVAL '24 hours' THEN 1 END) as recent_uploads
        FROM images 
        WHERE user_id = $1
    `

	GetImageFormatCounts = `
        SELECT format, COUNT(*) as count
        FROM images 
        WHERE user_id = $1
        GROUP BY format
        ORDER BY count DESC
    `

	GetImagesByIDs = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE id = ANY($1::int[]) AND user_id = $2
        ORDER BY id
    `

	BatchDeleteImages = `
        DELETE FROM images 
        WHERE id = ANY($1::int[]) AND user_id = $2
    `

	GetImagesByTourID = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE tour_id = $1
        ORDER BY sort_order ASC, id ASC
    `

	CheckImageExists = `
        SELECT COUNT(*) > 0
        FROM images 
        WHERE id = $1
    `

	GetImageWithTourInfo = `
        SELECT 
            i.id, i.user_id, i.tour_id, i.filename, i.original_name, i.url, i.thumbnail_url,
            i.size, i.width, i.height, i.format, i.description, i.alt_text, i.is_primary, i.sort_order,
            i.uploaded_at, i.updated_at,
            t.name as tour_name
        FROM images i
        LEFT JOIN tours t ON i.tour_id = t.id
        WHERE i.id = $1
    `

	GetRecentImages = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE user_id = $1
        ORDER BY uploaded_at DESC
        LIMIT $2
    `

	SearchImages = `
        SELECT 
            id, user_id, tour_id, filename, original_name, url, thumbnail_url,
            size, width, height, format, description, alt_text, is_primary, sort_order,
            uploaded_at, updated_at
        FROM images 
        WHERE user_id = $1
            AND (
                filename ILIKE $2 
                OR original_name ILIKE $2 
                OR description ILIKE $2
            )
        ORDER BY uploaded_at DESC
        LIMIT $3 OFFSET $4
    `

	CountSearchImages = `
        SELECT COUNT(*)
        FROM images 
        WHERE user_id = $1
            AND (
                filename ILIKE $2 
                OR original_name ILIKE $2 
                OR description ILIKE $2
            )
    `
)
