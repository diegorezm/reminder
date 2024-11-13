package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/diegorezm/reminder/internal/store"
)

type CreateReminderInput struct {
	Title   string
	DueDate time.Time
	Repeat  string
}

func createNotifications(ctx context.Context, db *sql.DB, s *store.Queries, reminderID int64, date []time.Time) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	qtx := s.WithTx(tx)

	for _, d := range date {
		err = qtx.CreateNotification(ctx, store.CreateNotificationParams{
			ReminderID: reminderID,
			DueDate:    d,
		})
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func CreateReminder(ctx context.Context, db *sql.DB, s *store.Queries, input CreateReminderInput) error {
	if input.Repeat == "" {
		id, err := s.CreateReminder(ctx, input.Title)
		if err != nil {
			return err
		}
		err = s.CreateNotification(ctx, store.CreateNotificationParams{
			ReminderID: id,
			DueDate:    input.DueDate,
		})
		if err != nil {
			return err
		}
		return nil
	}

	repeat := input.Repeat[1:]
	numberPart := repeat[:len(repeat)-1]
	unit := repeat[len(repeat)-1:]

	numberOfRepeats, err := strconv.Atoi(numberPart)

	if err != nil {
		return fmt.Errorf("Invalid repeat interval format: %w", err)
	}

	var interval time.Duration

	switch unit {
	case "d":
		interval = 24 * time.Hour
	case "w":
		interval = 7 * 24 * time.Hour
	case "m":
		interval = 30 * 24 * time.Hour
	case "y":
		interval = 365 * 24 * time.Hour
	default:
		return fmt.Errorf("Invalid repeat unit. Expected one of: d, w, m, y.")
	}

	// Multiply interval by the number of repeats
	interval *= time.Duration(numberOfRepeats)

	id, err := s.CreateReminder(ctx, input.Title)

	if err != nil {
		return err
	}

	dates := []time.Time{}

	d := input.DueDate

	// Generate reminders
	for i := 0; i < numberOfRepeats; i++ {
		dates = append(dates, d)
		d = d.Add(interval)
	}

	err = createNotifications(ctx, db, s, id, dates)

	if err != nil {
		return err
	}

	return nil
}
