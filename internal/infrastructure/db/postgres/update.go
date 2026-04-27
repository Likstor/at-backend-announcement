package postgres

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/pkg/apperror"
	"context"
)

const queryUpdate = `
	UPDATE announcements
	SET 
		title = $1,
		description = $2,
		updated_at = NOW() AT TIME ZONE 'utc';
		updated_by = $3
	WHERE 
		id = $4;
`

func (ar announcementRepository) Update(ctx context.Context, announcement domain.Announcement) error {
	resp, err := ar.Conn(ctx).Exec(
		ctx,
		queryUpdate,
		announcement.Title,
		announcement.Description,
		announcement.UpdatedBy,
		announcement.ID,
	)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrAnnouncementNotExists)
	}

	return nil
}
