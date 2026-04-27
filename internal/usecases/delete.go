package usecase

import "context"

func (as announcementUsecase) Delete(ctx context.Context, id uint64) error {
	if err := as.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}