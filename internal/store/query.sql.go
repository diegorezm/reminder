// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package store

import (
	"context"
	"time"
)

const clean = `-- name: Clean :exec
DELETE FROM reminders
`

func (q *Queries) Clean(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clean)
	return err
}

const cleanExpired = `-- name: CleanExpired :exec
DELETE FROM notifications WHERE due_date < CURRENT_TIMESTAMP
`

func (q *Queries) CleanExpired(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, cleanExpired)
	return err
}

const createNotification = `-- name: CreateNotification :exec
INSERT INTO notifications (reminder_id, due_date) VALUES (?, ?)
`

type CreateNotificationParams struct {
	ReminderID int64
	DueDate    time.Time
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) error {
	_, err := q.db.ExecContext(ctx, createNotification, arg.ReminderID, arg.DueDate)
	return err
}

const createReminder = `-- name: CreateReminder :one
INSERT INTO reminders (title) VALUES (?) RETURNING id
`

func (q *Queries) CreateReminder(ctx context.Context, title string) (int64, error) {
	row := q.db.QueryRowContext(ctx, createReminder, title)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteReminder = `-- name: DeleteReminder :execrows
DELETE FROM reminders WHERE id = ?
`

func (q *Queries) DeleteReminder(ctx context.Context, id int64) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteReminder, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const dismissNotification = `-- name: DismissNotification :exec
UPDATE notifications SET dismissed_at = CURRENT_TIMESTAMP WHERE id = ?
`

func (q *Queries) DismissNotification(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, dismissNotification, id)
	return err
}

const getAboutToExpireNotifications = `-- name: GetAboutToExpireNotifications :many
SELECT id, reminder_id, due_date, dismissed_at FROM notifications 
WHERE dismissed_at IS NULL
AND reminder_id = ?
AND due_date BETWEEN datetime('now', 'start of day', 'utc') AND datetime('now', '+1 day', 'start of day', 'utc')
ORDER BY due_date ASC
`

func (q *Queries) GetAboutToExpireNotifications(ctx context.Context, reminderID int64) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, getAboutToExpireNotifications, reminderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.ReminderID,
			&i.DueDate,
			&i.DismissedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllAboutToExpireNotifications = `-- name: GetAllAboutToExpireNotifications :many
SELECT id, reminder_id, due_date, dismissed_at FROM notifications 
WHERE dismissed_at IS NULL
AND due_date BETWEEN datetime('now') AND datetime('now', '+1 day')
ORDER BY due_date ASC
`

func (q *Queries) GetAllAboutToExpireNotifications(ctx context.Context) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, getAllAboutToExpireNotifications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.ReminderID,
			&i.DueDate,
			&i.DismissedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllNotifications = `-- name: GetAllNotifications :many
SELECT id, reminder_id, due_date, dismissed_at FROM notifications 
WHERE reminder_id = ?
ORDER BY due_date ASC
`

func (q *Queries) GetAllNotifications(ctx context.Context, reminderID int64) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, getAllNotifications, reminderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.ReminderID,
			&i.DueDate,
			&i.DismissedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDismissedNotifications = `-- name: GetDismissedNotifications :many
SELECT id, reminder_id, due_date, dismissed_at FROM notifications 
WHERE dismissed_at IS NOT NULL
ORDER BY reminder_id DESC
`

func (q *Queries) GetDismissedNotifications(ctx context.Context) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, getDismissedNotifications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.ReminderID,
			&i.DueDate,
			&i.DismissedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getExpiredNotifications = `-- name: GetExpiredNotifications :many
SELECT id, reminder_id, due_date, dismissed_at FROM notifications 
WHERE due_date < CURRENT_TIMESTAMP
ORDER BY reminder_id DESC
`

func (q *Queries) GetExpiredNotifications(ctx context.Context) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, getExpiredNotifications)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.ReminderID,
			&i.DueDate,
			&i.DismissedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotificationByID = `-- name: GetNotificationByID :one
SELECT id, reminder_id, due_date, dismissed_at FROM notifications WHERE id = ?
`

func (q *Queries) GetNotificationByID(ctx context.Context, id int64) (Notification, error) {
	row := q.db.QueryRowContext(ctx, getNotificationByID, id)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.ReminderID,
		&i.DueDate,
		&i.DismissedAt,
	)
	return i, err
}

const getNotifications = `-- name: GetNotifications :many
SELECT id, reminder_id, due_date, dismissed_at FROM notifications 
WHERE dismissed_at IS NULL
AND reminder_id = ?
ORDER BY due_date ASC
`

func (q *Queries) GetNotifications(ctx context.Context, reminderID int64) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, getNotifications, reminderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.ReminderID,
			&i.DueDate,
			&i.DismissedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReminder = `-- name: GetReminder :one
SELECT id, title, created_at FROM reminders WHERE id = ? LIMIT 1
`

func (q *Queries) GetReminder(ctx context.Context, id int64) (Reminder, error) {
	row := q.db.QueryRowContext(ctx, getReminder, id)
	var i Reminder
	err := row.Scan(&i.ID, &i.Title, &i.CreatedAt)
	return i, err
}

const getReminders = `-- name: GetReminders :many
SELECT id, title, created_at FROM reminders
`

func (q *Queries) GetReminders(ctx context.Context) ([]Reminder, error) {
	rows, err := q.db.QueryContext(ctx, getReminders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reminder
	for rows.Next() {
		var i Reminder
		if err := rows.Scan(&i.ID, &i.Title, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateReminder = `-- name: UpdateReminder :exec
UPDATE reminders SET title = ? WHERE id = ?
`

type UpdateReminderParams struct {
	Title string
	ID    int64
}

func (q *Queries) UpdateReminder(ctx context.Context, arg UpdateReminderParams) error {
	_, err := q.db.ExecContext(ctx, updateReminder, arg.Title, arg.ID)
	return err
}
