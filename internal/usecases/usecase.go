package usecase

import (
	"at-backend-announcement/internal/domain"
	"context"
)

var dummyAnnouncement = domain.Announcement{}

type announcementRepo interface {
	GetByID(ctx context.Context, id uint64) (domain.Announcement, error)
	GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Announcement, error)
	GetPage(ctx context.Context, cursor, pageSize uint64) ([]domain.Announcement, error)

	Create(ctx context.Context, announcement domain.Announcement) (uint64, error)
	Update(ctx context.Context, announcement domain.Announcement) error
	Delete(ctx context.Context, id uint64) error
}

type announcementUsecase struct {
	repo        announcementRepo
	maxPageSize uint64
}

func NewAnnouncementUsecase(repo announcementRepo, maxPageSize uint64) *announcementUsecase {
	return &announcementUsecase{
		repo:        repo,
		maxPageSize: maxPageSize,
	}
}
