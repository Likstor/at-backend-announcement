package postgres

import (
	"at-backend-announcement/internal/pkg/apperror"
	"context"
)

const queryDelete = `
	DELETE FROM announcements
	WHERE id = $1
`

func (ar announcementRepository) Delete(ctx context.Context, id uint64) error {
	resp, err := ar.Conn(ctx).Exec(ctx, queryDelete, id)
	if err != nil {
		return apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	if resp.RowsAffected() < 1 {
		return apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrAnnouncementNotExists)
	}

	return nil
}
