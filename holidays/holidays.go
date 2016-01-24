package holidays

import "time"

// midnight returns date with zero time.
func midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// ISO8601DateFormat represents date in YYYY-MM-DD format.
const ISO8601DateFormat = "2006-01-02"

// parseDate returns parsed or zero date.
func parseDate(s string) (time.Time, error) {
	d, err := time.Parse(ISO8601DateFormat, s)
	return midnight(d), err
}