package handlers

import (
	"at-backend-announcement/internal/domain"
	"at-backend-announcement/internal/pkg/responses"
	"at-backend-announcement/internal/pkg/apperror"
	"at-backend-announcement/internal/pkg/reqctx"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type announcementUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (aur announcementUpdateRequest) Validate() error {
	var err error

	if aur.Title == nil {
		err = errors.Join(err, errors.New("title field is missing or null"))
	} else {
		if *aur.Title == "" {
			err = errors.Join(err, errors.New("title field is empty"))
		}
	}

	if aur.Description == nil {
		err = errors.Join(err, errors.New("description field is missing or null"))
	} else {
		if *aur.Description == "" {
			err = errors.Join(err, errors.New("description field is empty"))
		}
	}

	return err
}

func (h announcementHandlerForAdmins) updateAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcementID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.NotFound(r.Context(), w)
		return
	}

	var dto announcementUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	if err := dto.Validate(); err != nil {
		slog.ErrorContext(r.Context(), err.Error())

		responses.Error(
			r.Context(),
			w,
			http.StatusUnprocessableEntity,
			err.Error(),
		)

		return
	}

	userID := reqctx.GetUserID(r.Context())

	announcement := domain.Announcement{
		ID:          announcementID,
		Title:       *dto.Title,
		Description: *dto.Description,
		UpdatedBy:   *userID,
	}

	if err := h.usecase.Update(r.Context(), announcement); err != nil {
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
