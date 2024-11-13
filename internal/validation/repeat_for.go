package validation

import (
	"fmt"
	"strconv"
	"strings"
)

const invalid_repeat_for_error = "Invalid repeat argument. Consider the format: +1d, +1w, +1m, +1y"

func RepeatForValidator(repeatFor string) error {
	if repeatFor == "" {
		return nil
	}

	parsedFor := strings.Split(repeatFor, "")
	if len(parsedFor) < 3 {
		return fmt.Errorf(invalid_repeat_for_error)
	}

	if parsedFor[0] != "+" {
		return fmt.Errorf("Expected +, got %s. \n%s", parsedFor[0], invalid_repeat_for_error)
	}

	_, err := strconv.Atoi(parsedFor[1])

	if err != nil {
		return fmt.Errorf("Expected a number, got %s. \n%s", parsedFor[1], invalid_repeat_for_error)
	}

	if parsedFor[2] != "d" && parsedFor[2] != "w" && parsedFor[2] != "m" && parsedFor[2] != "y" {
		return fmt.Errorf(invalid_repeat_for_error)
	}

	return nil
}
