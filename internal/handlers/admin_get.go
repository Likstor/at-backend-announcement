package handlers

import (
	"at-backend-announcement/internal/pkg/apperror"
	"at-backend-announcement/internal/pkg/logs"
	"at-backend-announcement/internal/pkg/responses"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func (h announcementHandlerForAdmins) getAnnouncementByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	announcement, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		logs.Error(r.Context(), err)

		switch {
		case errors.Is(err, apperror.ErrAnnouncementNotExists):
			responses.NotFound(r.Context(), w)
		default:
			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	announcementResp := map[string]any{
		"id":          announcement.ID,
		"title":       announcement.Title,
		"description": announcement.Description,
		"created_at":  announcement.CreatedAt,
		"created_by":  announcement.CreatedBy,
		"updated_at":  announcement.UpdatedAt,
		"updated_by":  announcement.UpdatedBy,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		announcementResp,
	)
}
