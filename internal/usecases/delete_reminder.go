package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diegorezm/reminder/internal/store"
)

func DeleteReminder(id int64, s *store.Queries, ctx context.Context) error {
	_, err := s.DeleteReminder(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reminder with ID %d not found", id)
		}
		return err
	}
	return nil
}
