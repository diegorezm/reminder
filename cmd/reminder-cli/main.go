package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const DEFAULT_SQLITE_PATH = ".local/share/reminder/reminder.db"

func createDB() error {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return fmt.Errorf("HOME environment variable not set")
	}

	dbPath := filepath.Join(homeDir, DEFAULT_SQLITE_PATH)

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("No reminder.db file found. Creating...")
		err := os.MkdirAll(filepath.Dir(dbPath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		_, err = os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("failed to create reminder.db: %w", err)
		}
	}
	return nil
}

func main() {
	err := createDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
