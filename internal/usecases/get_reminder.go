package usecases

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/diegorezm/reminder/internal/store"
)

func GetReminder(id int64, s *store.Queries, ctx context.Context) (store.Reminder, error) {
	reminder, err := s.GetReminder(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return store.Reminder{}, fmt.Errorf("reminder with ID %d not found", id)
		}
		return store.Reminder{}, err
	}
	return reminder, nil
}
