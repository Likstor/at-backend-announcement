package handlers

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/handlers/middleware"
	"at-backend-announcement/internal/pkg/roles"
	"context"
	"fmt"
	"net/http"
)

type announcementUsecase interface {
	GetPage(ctx context.Context, cursor, pageSize uint64) ([]domain.Announcement, error)
	GetFirstPage(ctx context.Context, pageSize uint64) ([]domain.Announcement, error)
	GetByID(ctx context.Context, id uint64) (domain.Announcement, error)
}

type announcementHandler struct {
	usecase announcementUsecase
}

func NewAnnouncementHandler(usecase announcementUsecase) *announcementHandler {
	return &announcementHandler{
		usecase: usecase,
	}
}

func (a announcementHandler) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("GET /page", a.getPage)
	muxWithAuth.HandleFunc("GET /{id}", a.getByID)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.User, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}

type announcementUsecaseForAdmins interface {
	announcementUsecase

	Create(ctx context.Context, announcement domain.Announcement) (uint64, error)
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, announcement domain.Announcement) error
}

type announcementHandlerForAdmins struct {
	usecase announcementUsecaseForAdmins
}

func NewAnnouncementHandlerForAdmins(usecase announcementUsecaseForAdmins) *announcementHandlerForAdmins {
	return &announcementHandlerForAdmins{
		usecase: usecase,
	}
}

func (a announcementHandlerForAdmins) Setup(prefix string, verifier middleware.Verifier, mux *http.ServeMux) {
	muxWithAuth := http.NewServeMux()

	muxWithAuth.HandleFunc("POST /", a.createAnnouncement)
	muxWithAuth.HandleFunc("GET /{id}", a.getAnnouncementByID)
	muxWithAuth.HandleFunc("DELETE /{id}", a.deleteAnnouncement)
	muxWithAuth.HandleFunc("PUT /{id}", a.updateAnnouncement)

	muxWithAuthWrapped := middleware.Authorization(muxWithAuth, roles.Operator, verifier)

	mux.Handle(fmt.Sprintf("%s/", prefix), http.StripPrefix(prefix, muxWithAuthWrapped))
}