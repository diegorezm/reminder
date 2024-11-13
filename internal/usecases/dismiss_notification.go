package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diegorezm/reminder/internal/store"
)

func DismissNotification(id int64, s *store.Queries, ctx context.Context) error {
	not, err := s.GetNotificationByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("notification with ID %d not found", id)
		}
		return fmt.Errorf("failed to get notification with ID %d: %w", id, err)
	}

	if not.DismissedAt.Valid {
		return fmt.Errorf("notification with ID %d is already dismissed", id)
	}

	return s.DismissNotification(ctx, id)
}
