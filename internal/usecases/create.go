package usecase

import (
	"at-backend-announcement/internal/domain"
	"context"
)

func (as announcementUsecase) Create(ctx context.Context, announcement domain.Announcement) (uint64, error) {
	id, err := as.repo.Create(ctx, announcement)
	if err != nil {
		return 0, err
	}

	return id, nil
}