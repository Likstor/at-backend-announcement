package handlers

import (
	"at-backend-announcement/internal/pkg/responses"
	"at-backend-announcement/internal/pkg/apperror"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func (h announcementHandlerForAdmins) deleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)

		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, apperror.ErrAnnouncementNotExists):
			responses.NotFound(r.Context(), w)
		default:
			slog.ErrorContext(apperror.GetErrorCtx(r.Context(), err), err.Error())

			responses.InternalServerError(r.Context(), w)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}