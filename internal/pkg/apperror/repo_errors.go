package apperror

import "errors"

var (
	ErrRepository = errors.New("repository internal error")
	
	ErrAnnouncementNotExists = errors.New("announcement is not found")
)
