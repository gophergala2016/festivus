package holidays

import (
	"fmt"
	"math"
	"time"
)

// New returns all holidays for country.
func New(countryCode, path string) ([]Hday, error) {
	p := fmt.Sprintf("%s/%s.txt", path, countryCode)
	return parseFile(p)
}

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

// DaysBetween returns days between dates.
func DaysBetween(from, to time.Time) int {
	from = midnight(from)
	to = midnight(to)
	// convert diff hours to days
	d := to.Sub(from).Hours() / 24
	return int(math.Abs(d))
}

// NextFestivus returns date of next Festivus.
func NextFestivus(today time.Time) time.Time {
	year := today.Year()
	fd := time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	if today.After(fd) {
		year++
		fd = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	}
	return fd
}

// DaysToFestivus returns number of days from to today to festivus
func DaysToFestivus(today time.Time) int {
	nf := NextFestivus(today)
	return DaysBetween(today, nf)
}

// ByYear filter holidays by year.
func ByYear(hlds []Hday, today time.Time) []Hday {
	year := today.Year()
	all := []Hday{}
	for _, h := range hlds {
		hy := h.Date().Year()
		if hy < year {
			continue
		}
		if hy > year {
			break
		}
		all = append(all, h)
	}
	return all
}
