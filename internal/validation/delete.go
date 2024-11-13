package validation

import "fmt"

// EXPECTS: ./binary reminder delete <reminder_id>

const delete_reminder_error = "Please consider using the following format: ./binary reminder delete <reminder_id>"

func DeleteReminderArgsValidator(args []string) (string, error) {
	if len(args) < 4 {
		return "", fmt.Errorf("Invalid number of arguments. Expected 2, got %d. %s", len(args), create_reminder_error)
	}
	deleteArgs := args[3:]
	id := deleteArgs[0]
	if id == "" {
		return "", fmt.Errorf("Invalid reminder id. Expected a non-empty string. %s", delete_reminder_error)
	}

	return id, nil
}
