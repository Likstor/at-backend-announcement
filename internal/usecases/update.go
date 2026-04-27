package usecase

import (
	"at-backend-announcement/internal/domain"
	"context"
)

func (as announcementUsecase) Update(ctx context.Context, announcement domain.Announcement) error {
	if err := as.repo.Update(ctx, announcement); err != nil {
		return err
	}

	return nil
}