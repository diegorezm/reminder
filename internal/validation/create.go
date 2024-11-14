package validation

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EXPECTS: ./binary reminder create "title" "01-01-2024" +1d
const create_reminder_error = "Please consider using the following format: ./binary create \"title\" \"01-01-2024\" +1d"

func CreateReminderArgsValidator(args []string) (string, *time.Time, string, error) {
	if len(args) < 4 {
		return "", nil, "", fmt.Errorf("Invalid number of arguments. Expected 3, got %d. %s", len(args), create_reminder_error)
	}

	createArgs := args[2:]

	if len(createArgs) < 2 {
		return "", nil, "", fmt.Errorf("Invalid number of arguments. Expected at least 2, got %d. %s", len(createArgs), create_reminder_error)
	}

	title := createArgs[0]

	if title == "" {
		return "", nil, "", fmt.Errorf("Invalid title. Expected a non-empty string. %s", create_reminder_error)
	}

	// Parse due date in "d-m-y" format (e.g., "01-01-2024 00:00:00")
	dueDate, err := parseDueDate(createArgs[1])
	if err != nil {
		return "", nil, "", err
	}

	if len(createArgs) < 3 {
		return title, &dueDate, "", nil
	}

	repeatFor := createArgs[2]
	err = RepeatForValidator(repeatFor)
	if err != nil {
		return "", nil, "", err
	}

	return title, &dueDate, repeatFor, nil
}

const create_notification_error = "Please consider using the following format: ./binary create notification <id> <date>"

func CreateNotificationArgsValidator(args []string) (int64, *time.Time, error) {
	if len(args) < 4 {
		return 0, nil, fmt.Errorf("Invalid number of arguments. Expected 3, got %d. %s", len(args), create_notification_error)
	}
	// bin [0] create [1] notification [2] id [3] date [4]
	createArgs := args[3:]

	if len(createArgs) < 2 {
		return 0, nil, fmt.Errorf("Invalid number of arguments. Expected at least 2, got %d. %s", len(createArgs), create_notification_error)
	}

	id, err := strconv.Atoi(createArgs[0])
	if err != nil {
		return 0, nil, fmt.Errorf("Invalid ID: %w", err)
	}

	// Parse due date in "d-m-y" format (e.g., "01-01-2024 00:00:00")
	dueDate, err := parseDueDate(createArgs[1])
	if err != nil {
		return 0, nil, err
	}

	return int64(id), &dueDate, nil
}

func parseDueDate(input string) (time.Time, error) {
	if !strings.Contains(input, " ") {
		input += " 00:00"
	}
	parsedDate := strings.Replace(input, "/", "-", -1)

	dueDate, err := time.Parse("02-01-2006 15:04", parsedDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected 'dd-mm-yyyy HH:mm' or 'dd-mm-yyyy': %w", err)
	}

	return dueDate, nil
}
