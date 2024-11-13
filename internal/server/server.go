package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/diegorezm/reminder/internal/store"
	"github.com/diegorezm/reminder/internal/templates/pages"
	"github.com/diegorezm/reminder/internal/usecases"
)

type Server struct {
	st *store.Queries
	db *sql.DB
}

func New(st *store.Queries, db *sql.DB) *Server {
	return &Server{
		st: st,
		db: db,
	}
}

func (s *Server) StartServer() {
	port := ":4242"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reminders, err := s.st.GetReminders(ctx)

		if err != nil {
			fmt.Println(err)
		}

		notifications, err := s.st.GetAllAboutToExpireNotifications(ctx)

		index := pages.Index(reminders, notifications)

		if err := index.Render(ctx, w); err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			log.Println("Error rendering page:", err)
			return
		}
	})

	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pathId := r.URL.Path[1:]

		id, err := strconv.Atoi(pathId)

		reminder, err := s.st.GetReminder(ctx, int64(id))

		if err != nil {
			fmt.Println(err)
		}

		notifications, err := usecases.ListNotification(usecases.ListNotificationInput{
			ReminderID: reminder.ID,
			Type:       usecases.NOTIFICATION_TYPE_ALL_IGNORE_DATE,
		}, s.st, ctx)

		if err != nil {
			fmt.Println(err)
		}

		reminderPage := pages.Reminder(reminder, notifications)

		if err := reminderPage.Render(ctx, w); err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			log.Println("Error rendering page:", err)
			return
		}
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		create := pages.Create()

		if err := create.Render(ctx, w); err != nil {
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
			log.Println("Error rendering page:", err)
			return
		}
	})

	http.HandleFunc("/api/create", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		title := r.FormValue("title")

		if title == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Title is required"))
			return
		}

		dateStr := r.FormValue("date")

		if dateStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Date is required"))
			return
		}

		d := strings.Split(dateStr, "T")
		dt := strings.Join(d, " ")

		date, err := time.Parse("2006-01-02 15:04", dt)

		if err != nil {
			log.Printf("Error parsing date: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid date format"))
			return
		}

		repeat := r.FormValue("repeat")

		err = usecases.CreateReminder(ctx, s.db, s.st, usecases.CreateReminderInput{
			Title:   title,
			DueDate: date,
			Repeat:  repeat,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error creating reminder:", err)
			w.Write([]byte("Failed to create reminder"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Reminder created successfully"))
	})

	http.HandleFunc("/api/delete", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.FormValue("id")

		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			w.Write([]byte("ID is required"))
			return
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			w.Write([]byte("Invalid ID"))
			return
		}

		err = usecases.DeleteReminder(int64(idInt), s.st, ctx)

		if err != nil {
			http.Error(w, "Failed to delete reminder", http.StatusInternalServerError)
			log.Println("Error deleting reminder:", err)
			w.Write([]byte("Failed to delete reminder"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reminder deleted successfully"))
	})

	http.HandleFunc("/api/dismiss", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.FormValue("id")

		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			w.Write([]byte("ID is required"))
			return
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			w.Write([]byte("Invalid ID"))
			return
		}

		err = usecases.DismissNotification(int64(idInt), s.st, ctx)

		if err != nil {
			http.Error(w, "Failed to dismiss notification", http.StatusInternalServerError)
			log.Println("Error dismissing notification:", err)
			w.Write([]byte("Failed to dismiss notification"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification dismissed successfully"))
	})

	fmt.Printf("Server started on http://localhost%s/\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
