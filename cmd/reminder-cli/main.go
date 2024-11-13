package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
	"github.com/diegorezm/reminder/internal/migrations"
	"github.com/diegorezm/reminder/internal/server"
	"github.com/diegorezm/reminder/internal/store"
	"github.com/diegorezm/reminder/internal/usecases"
	"github.com/diegorezm/reminder/internal/validation"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DEFAULT_SQLITE_PATH = ".local/share/reminder/reminder.db"
	DEFAULT_DATE_FORMAT = "02/01/2006"
)

func printTableNotifications(notifications []store.Notification, r store.Reminder) {
	t := table.New(os.Stdout)
	t.SetHeaders("#", "Title", "Date", "Dismissed", "isExpired")
	for _, n := range notifications {
		id := strconv.FormatInt(n.ID, 10)

		var expiredText string

		isExpired := n.DueDate.Before(time.Now()) && n.DismissedAt.Valid == false

		if isExpired {
			expiredText = "Yes"
		} else {
			expiredText = "No"
		}
		var dismissedText string

		if n.DismissedAt.Valid {
			dismissedText = n.DismissedAt.Time.Format(DEFAULT_DATE_FORMAT)
		} else {
			dismissedText = "Not dismissed yet"
		}

		t.AddRow(id, r.Title, n.DueDate.Format(DEFAULT_DATE_FORMAT), dismissedText, expiredText)
	}
	t.Render()
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

func help() {
	fmt.Println("\nReminder CLI - Command Line Tool")
	fmt.Println("Usage:")
	fmt.Println("  reminder-cli <command> [options]")
	fmt.Println("\nCommands:")

	fmt.Println("  help")
	fmt.Println("      Displays this help message.")

	fmt.Println("\n  serve")
	fmt.Println("      Starts the server, so you can manage reminders from your browser.")

	fmt.Println("\n  create <title> <date> [repeat]")
	fmt.Println("      Creates a new reminder.")
	fmt.Println("      - <title>: Title of the reminder.")
	fmt.Println("      - <date>:  Due date (e.g., 31/12/2024 15:00:00).")
	fmt.Println("      - [repeat]: Optional repetition interval (e.g., '3d' for 3 days).")

	fmt.Println("\n  delete <id>")
	fmt.Println("      Deletes a reminder by its ID (also deletes associated notifications).")

	fmt.Println("\n  dismiss <id>")
	fmt.Println("      Dismisses a notification by its ID.")

	fmt.Println("\n  list")
	fmt.Println("      List reminders, notifications, or filtered notifications.")

	fmt.Println("\n    reminders")
	fmt.Println("      Lists all reminders.")

	fmt.Println("\n    notifications <subcommand>")
	fmt.Println("      Lists notifications based on subcommand:")
	fmt.Println("        all                Lists all notifications.")
	fmt.Println("        about-to-expire    Lists notifications that are about to expire (within 1 day). This runs by default.")
	fmt.Println("        dismissed          Lists all dismissed notifications.")
	fmt.Println("        expired            Lists all expired notifications.")
	fmt.Println("        reminder <id>      List notifications for a specific reminder.")

	fmt.Println("\nExamples:")
	fmt.Println("  reminder-cli create \"Meeting\" \"31/12/2024 15:00:00\" \"+3d\"")
	fmt.Println("  reminder-cli list notifications all")
	fmt.Println("  reminder-cli dismiss 5")
	fmt.Println("  reminder-cli delete 2")
	fmt.Println()
}

func createDB() (*sql.DB, error) {
	var home string
	if home = os.Getenv("HOME"); home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}

	isFresh := false

	pathDir := filepath.Join(home, ".local", "share", "reminder")
	// check if the directory exists
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		isFresh = true
		fmt.Println("Creating directory:", pathDir)
		if err := os.MkdirAll(pathDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		if _, err := os.Create(filepath.Join(pathDir, "reminder.db")); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	dbPath := filepath.Join(pathDir, "reminder.db")
	_, err := os.Stat(dbPath)

	// check if the file exists
	if os.IsNotExist(err) {
		isFresh = true
		fmt.Println("Creating database...")
		if _, err := os.Create(filepath.Join(pathDir, "reminder.db")); err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if isFresh {
		m := migrations.New(db)
		m.Up()
	}

	return db, nil
}

func main() {
	db, err := createDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()
	s := store.New(db)

	if len(os.Args) < 2 {
		help()
		return
	}

	command := os.Args[1]
	switch command {
	case "serve":
		fmt.Println("Starting server...")
		sv := server.New(s, db)
		sv.StartServer()
	case "create":
		title, date, repeat, err := validation.CreateReminderArgsValidator(os.Args)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = usecases.CreateReminder(ctx, db, s, usecases.CreateReminderInput{
			Title:   title,
			DueDate: *date,
			Repeat:  repeat,
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Notifications created for reminder with title: %s\n", title)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Insufficient arguments for deleting a reminder.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		err = usecases.DeleteReminder(int64(id), s, ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reminder with ID: %d deleted\n", id)

	case "dismiss":
		if len(os.Args) < 3 {
			fmt.Println("Insufficient arguments for dismissing a notification.")
			return
		}
		// bin dismiss <id>
		id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		err = usecases.DismissNotification(int64(id), s, ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Notification with ID: %d dismissed\n", id)
	case "list":
		if len(os.Args) < 3 {
			fmt.Println("Specify what to list (reminders, notifications, dismissed , expired).")
			return
		}
		switch os.Args[2] {
		case "reminders":
			reminders, err := s.GetReminders(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			t := table.New(os.Stdout)
			t.SetHeaders("#", "Title")
			if len(reminders) == 0 {
				fmt.Println("No reminders found.")
				return
			}
			for _, r := range reminders {
				id := strconv.FormatInt(r.ID, 10)
				t.AddRow(id, r.Title)
			}
			t.Render()
		case "notifications":
			reminders, err := s.GetReminders(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(os.Args) < 4 {
				for _, r := range reminders {
					notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
						ReminderID: r.ID,
						Type:       usecases.NOTIFICATION_TYPE_ABOUT_TO_EXPIRE,
					}, s, ctx)

					if err != nil {
						fmt.Println(err)
						return
					}

					if len(notifications) == 0 {
						continue
					}

					printTableNotifications(notifications, r)
				}
				return
			}

			switch os.Args[3] {
			case "all":
				for _, r := range reminders {
					notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
						ReminderID: r.ID,
						Type:       usecases.NOTIFICATION_TYPE_ALL,
					}, s, ctx)
					if err != nil {
						fmt.Println(err)
						return
					}
					if len(notifications) == 0 {
						continue
					}
					printTableNotifications(notifications, r)
				}
			case "reminder":
				if len(os.Args) < 4 {
					fmt.Println("Insufficient arguments for listing notifications for a specific reminder.")
					return
				}
				id, err := strconv.Atoi(os.Args[4])
				if err != nil {
					fmt.Println("Invalid ID:", err)
					return
				}
				reminder, err := usecases.GetReminder(int64(id), s, ctx)
				if err != nil {
					fmt.Println(err)
					return
				}
				notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
					ReminderID: reminder.ID,
					Type:       usecases.NOTIFICATION_TYPE_ALL_IGNORE_DATE,
				}, s, ctx)
				if err != nil {
					fmt.Println(err)
					return
				}
				printTableNotifications(notifications, reminder)
			case "dismissed":
				reminders, err := s.GetReminders(ctx)
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, r := range reminders {
					notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
						ReminderID: r.ID,
						Type:       usecases.NOTIFICATION_TYPE_DISMISSED,
					}, s, ctx)
					if err != nil {
						fmt.Println(err)
						return
					}
					if len(notifications) == 0 {
						continue
					}
					printTableNotifications(notifications, r)
				}
			case "expired":
				reminders, err := s.GetReminders(ctx)
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, r := range reminders {
					notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
						ReminderID: r.ID,
						Type:       usecases.NOTIFICATION_TYPE_EXPIRED,
					}, s, ctx)
					if err != nil {
						fmt.Println(err)
						return
					}
					if len(notifications) == 0 {
						continue
					}
					printTableNotifications(notifications, r)
				}
			default:
				for _, r := range reminders {
					notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
						ReminderID: r.ID,
						Type:       usecases.NOTIFICATION_TYPE_ABOUT_TO_EXPIRE,
					}, s, ctx)
					if err != nil {
						fmt.Println(err)
						return
					}
					if len(notifications) == 0 {
						continue
					}
					printTableNotifications(notifications, r)
				}
			}
		default:
			fmt.Println("Unknown list option. Use one of: reminders, notifications, dismissed, expired.")
		}
	case "help":
		help()

	default:
		fmt.Println("Unknown command.")
		help()
	}
	os.Exit(0)
}
