package postgres

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/pkg/apperror"
	"context"
)

const queryCreate = `
	INSERT INTO announcements
		(title, description, created_by, updated_by)
	VALUES
		($1, $2, $3, $3)
	RETURNING id
`

func (ar announcementRepository) Create(ctx context.Context, announcement domain.Announcement) (uint64, error) {
	var id uint64

	if err := ar.Conn(ctx).QueryRow(ctx, queryCreate, announcement.Title, announcement.Description, announcement.CreatedBy).Scan(&id); err != nil {
		return id, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return id, nil
}
