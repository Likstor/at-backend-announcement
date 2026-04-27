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
)

type createAnnouncementRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (car createAnnouncementRequest) Validate() error {
	var err error

	if car.Title == nil {
		err = errors.Join(err, errors.New("title field is missing or null"))
	} else {
		if *car.Title == "" {
			err = errors.Join(err, errors.New("title field is empty"))
		}
	}

	if car.Description == nil {
		err = errors.Join(err, errors.New("description field is missing or null"))
	} else {
		if *car.Description == "" {
			err = errors.Join(err, errors.New("description field is empty"))
		}
	}

	return err
}

func (h announcementHandlerForAdmins) createAnnouncement(w http.ResponseWriter, r *http.Request) {
	var dto createAnnouncementRequest

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
		Title:       *dto.Title,
		Description: *dto.Description,
		CreatedBy:   *userID,
	}

	id, err := h.usecase.Create(r.Context(), announcement)
	if err != nil {
		slog.ErrorContext(apperror.GetErrorCtx(r.Context(), err), err.Error())

		responses.InternalServerError(r.Context(), w)
		return
	}

	resp := map[string]any{
		"announcement_id": id,
	}

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}
