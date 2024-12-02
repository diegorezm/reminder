-- name: CreateReminder :one
INSERT INTO reminders (title) VALUES (?) RETURNING id;

-- name: GetReminder :one
SELECT * FROM reminders WHERE id = ? LIMIT 1;

-- name: GetReminders :many
SELECT * FROM reminders;

-- name: UpdateReminder :exec
UPDATE reminders SET title = ? WHERE id = ?;

-- name: DeleteReminder :execrows
DELETE FROM reminders WHERE id = ?;

-- name: CreateNotification :exec
INSERT INTO notifications (reminder_id, due_date) VALUES (?, ?);

-- name: GetNotifications :many
SELECT * FROM notifications 
WHERE dismissed_at IS NULL
AND reminder_id = ?
ORDER BY due_date ASC;

-- name: GetAllNotifications :many
SELECT * FROM notifications 
WHERE reminder_id = ?
ORDER BY due_date ASC;

-- name: GetAboutToExpireNotifications :many
SELECT * FROM notifications 
WHERE dismissed_at IS NULL
AND reminder_id = ?
AND due_date BETWEEN datetime('now', 'start of day', 'utc') AND datetime('now', '+1 day', 'start of day', 'utc')
ORDER BY due_date ASC;

-- name: GetDismissedNotifications :many
SELECT * FROM notifications 
WHERE dismissed_at IS NOT NULL
ORDER BY reminder_id DESC;

-- name: GetExpiredNotifications :many
SELECT * FROM notifications 
WHERE due_date < CURRENT_TIMESTAMP
ORDER BY reminder_id DESC;

-- name: DismissNotification :exec
UPDATE notifications SET dismissed_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: GetNotificationByID :one
SELECT * FROM notifications WHERE id = ?;

-- name: GetAllAboutToExpireNotifications :many
SELECT * FROM notifications 
WHERE dismissed_at IS NULL
AND due_date BETWEEN datetime('now') AND datetime('now', '+1 day')
ORDER BY due_date ASC;

-- name: CleanExpired :exec
DELETE FROM notifications WHERE due_date < CURRENT_TIMESTAMP;

-- name: Clean :exec
DELETE FROM reminders;
