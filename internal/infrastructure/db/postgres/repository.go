package postgres

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/pkg/transactor"

	"github.com/jackc/pgx/v5/pgxpool"
)

type announcementRepository struct {
	transactor.RepositoryWithTransactor
}

func NewAnnouncementRepository(pool *pgxpool.Pool) *announcementRepository {
	return &announcementRepository{
		RepositoryWithTransactor: *transactor.NewRepositoryWithTransactor(pool),
	}
}

var dummyAnnouncement = domain.Announcement{}