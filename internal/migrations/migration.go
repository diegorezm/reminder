package migrations

import "database/sql"

type Migration struct {
	db *sql.DB
}

func New(db *sql.DB) *Migration {
	return &Migration{db: db}
}

func (m *Migration) Up() {
	if err := upRemindersTable(m.db); err != nil {
		panic(err)
	}
	if err := upNotificationsTable(m.db); err != nil {
		panic(err)
	}
}

func (m *Migration) Down() {
	if err := downRemindersTable(m.db); err != nil {
		panic(err)
	}

	if err := downNotificationsTable(m.db); err != nil {
		panic(err)
	}
}

func (m *Migration) Reset() {
	m.Down()
	m.Up()
}

func upRemindersTable(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS reminders (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT NOT NULL,
      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
  `)

	if err != nil {
		return err
	}
	return nil
}

func downRemindersTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE reminders;`)
	if err != nil {
		return err
	}
	return nil
}

func upNotificationsTable(db *sql.DB) error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS notifications (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      reminder_id INTEGER NOT NULL,
      due_date DATETIME NOT NULL,
      dismissed_at DATETIME,
      FOREIGN KEY (reminder_id) REFERENCES reminders(id) ON DELETE CASCADE
    );
  `)

	if err != nil {
		return err
	}
	return nil
}

func downNotificationsTable(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE notifications;`)
	if err != nil {
		return err
	}
	return nil
}
