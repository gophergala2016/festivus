package holidays

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Hday holds holiday's data.
type Hday struct {
	date    time.Time
	endDate time.Time
	name    string
}

// Date returns holiday's start date.
func (h *Hday) Date() time.Time {
	return h.date
}

// EndDate returns holiday's end date.
func (h *Hday) EndDate() time.Time {
	return h.endDate
}

// Name returns holiday's name.
func (h *Hday) Name() string {
	return h.name
}

// String returns holiday's string representation.
func (h *Hday) String() string {
	return fmt.Sprintf(
		"%s %s %s",
		h.Date().Format(ISO8601DateFormat),
		h.EndDate().Format(ISO8601DateFormat),
		h.Name(),
	)
}

func NewHday(in string) (*Hday, error) {
	fs := strings.Split(in, "\t")
	if len(fs) != 5 {
		return nil, errors.New("hday: invalid data row")
	}

	// start date
	d, err := parseDate(fs[0])
	if err != nil {
		return nil, errors.New("hday: invalid start date")
	}

	// end date
	ed, err := parseDate(fs[2])
	if err != nil {
		return nil, errors.New("hday: invalid end date")
	}

	// name
	n := fs[4]

	// valid holiday
	h := &Hday{
		date:    d,
		endDate: ed,
		name:    n,
	}

	return h, nil
}
