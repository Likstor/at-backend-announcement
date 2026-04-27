package handlers

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/pkg/responses"
	"at-backend-announcement/internal/pkg/apperror"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func (h announcementHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	announcement, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrAnnouncementNotExists):
			responses.NotFound(r.Context(), w)
		default:
			slog.ErrorContext(apperror.GetErrorCtx(r.Context(), err), err.Error())

			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	announcementResp := map[string]any{
		"id":          announcement.ID,
		"title":       announcement.Title,
		"description": announcement.Description,
		"created_at":  announcement.CreatedAt,
		"updated_at":  announcement.UpdatedAt,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		announcementResp,
	)
}

func (h announcementHandler) pageToRespPage(page []domain.Announcement) []map[string]any {
	pageResp := make([]map[string]any, 0, len(page))
	for _, announcement := range page {
		announcementResp := map[string]any{
			"id":          announcement.ID,
			"title":       announcement.Title,
			"description": announcement.Description,
			"created_at":  announcement.CreatedAt,
			"updated_at":  announcement.UpdatedAt,
		}

		pageResp = append(pageResp, announcementResp)
	}

	return pageResp
}

func (h announcementHandler) getFirstPage(w http.ResponseWriter, r *http.Request, pageSize uint64) {
	page, err := h.usecase.GetFirstPage(r.Context(), pageSize)
	if err != nil {
		slog.ErrorContext(apperror.GetErrorCtx(r.Context(), err), err.Error())

		responses.InternalServerError(r.Context(), w)
		return
	}

	announcementsResp := h.pageToRespPage(page)

	resp := map[string]any{
		"announcements": announcementsResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (h announcementHandler) getPage(w http.ResponseWriter, r *http.Request) {
	pageSize := getPageSize(r.URL.Query())

	cursorString := r.URL.Query().Get("cursor")
	if cursorString == "" {
		h.getFirstPage(w, r, pageSize)
		return
	}

	cursor, err := strconv.ParseUint(cursorString, 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	page, err := h.usecase.GetPage(r.Context(), cursor, pageSize)
	if err != nil {
		slog.ErrorContext(apperror.GetErrorCtx(r.Context(), err), err.Error())
		
		responses.InternalServerError(r.Context(), w)
		return
	}

	announcementsResp := h.pageToRespPage(page)

	resp := map[string]any{
		"announcements": announcementsResp,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
