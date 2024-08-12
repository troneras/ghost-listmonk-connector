package utils

import (
	"fmt"
	"time"
)

func ParseDuration(durationStr string) (time.Duration, error) {
	// First, try to parse with Go's built-in time.ParseDuration
	duration, err := time.ParseDuration(durationStr)
	if err == nil {
		return duration, nil
	}

	// If that fails, check for days or weeks
	var numericValue int
	var unit string

	_, err = fmt.Sscanf(durationStr, "%d%s", &numericValue, &unit)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	switch unit {
	case "d":
		return time.Duration(numericValue) * 24 * time.Hour, nil
	case "w":
		return time.Duration(numericValue) * 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported duration unit: %s", unit)
	}
}

func FormatDuration(d time.Duration) string {
	if d.Hours() >= 24*7 {
		weeks := int(d.Hours() / (24 * 7))
		return fmt.Sprintf("%dw", weeks)
	} else if d.Hours() >= 24 {
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
	return d.String()
}
