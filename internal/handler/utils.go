package handler

import "time"

// parseDate parses an ISO date string (YYYY-MM-DD or full ISO) into *time.Time
func parseDate(s string) *time.Time {
	// Try YYYY-MM-DD first
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		// Try full ISO
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return nil
		}
	}
	return &t
}
