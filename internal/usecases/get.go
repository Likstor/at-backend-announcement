package usecase

import (
	"at-backend-announcement/internal/domain"
	"context"
)

func (as announcementUsecase) GetPage(ctx context.Context, cursor, pageSize uint64) ([]domain.Announcement, error) {
	if pageSize > as.maxPageSize {
		pageSize = as.maxPageSize
	}

	announcements, err := as.repo.GetPage(ctx, cursor, pageSize)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func (as announcementUsecase) GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Announcement, error) {
	if pageSize > as.maxPageSize {
		pageSize = as.maxPageSize
	}
	
	announcements, err := as.repo.GetFirstPage(ctx, pageSize)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func (as announcementUsecase) GetByID(ctx context.Context, id uint64) (domain.Announcement, error) {
	announcement, err := as.repo.GetByID(ctx, id)
	if err != nil {
		return dummyAnnouncement, err
	}

	return announcement, nil
}