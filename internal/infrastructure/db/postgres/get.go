package postgres

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/pkg/apperror"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

const queryGetByID = `
	SELECT
		id, title, description, created_at, created_by, updated_at, updated_by
	FROM announcements
	WHERE 
		id = $1
`

func (ar announcementRepository) GetByID(ctx context.Context, id uint64) (domain.Announcement, error) {
	rows, err := ar.Conn(ctx).Query(ctx, queryGetByID, id)
	if err != nil {
		return dummyAnnouncement, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	announcement, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[domain.Announcement])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dummyAnnouncement, apperror.NewErrorCtxWithoutMsg(ctx, apperror.ErrAnnouncementNotExists)
		}

		return dummyAnnouncement, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	announcement.ID = id

	return announcement, nil
}

func (ar announcementRepository) getAnnouncementFromRows(ctx context.Context, rows pgx.Rows) ([]domain.Announcement, error) {
	claims, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.Announcement])
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return claims, nil
}

const queryGetFirstPage = `
	SELECT 
		id, title, description, created_at, updated_at
	FROM announcements
	ORDER BY id DESC
	LIMIT $1
`

func (ar announcementRepository) GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Announcement, error) {
	rows, err := ar.Conn(ctx).Query(ctx, queryGetFirstPage, pageSize)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return ar.getAnnouncementFromRows(ctx, rows)
}

const queryGetPage = `
	SELECT 
		id, title, description, created_at, updated_at
	FROM announcements
	WHERE id < $1
	ORDER BY id DESC
	LIMIT $2
`

func (ar announcementRepository) GetPage(ctx context.Context, cursor, pageSize uint64) ([]domain.Announcement, error) {
	rows, err := ar.Conn(ctx).Query(ctx, queryGetPage, cursor, pageSize)
	if err != nil {
		return nil, apperror.NewErrorCtx(ctx, err.Error(), apperror.ErrRepository)
	}

	return ar.getAnnouncementFromRows(ctx, rows)
}
