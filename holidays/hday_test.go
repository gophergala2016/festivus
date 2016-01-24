package holidays

import (
	"reflect"
	"testing"
	"time"
)

func fakeHday(d, ed time.Time, n string) Hday {
	return Hday{
		date:    d,
		endDate: ed,
		name:    n,
	}
}

func TestNewHday(t *testing.T) {
	tests := []struct {
		in   string
		want Hday
	}{
		{"2016-01-01	00:00	2016-01-02	00:00	New Year's Day",
			fakeHday(
				fakeDate(2016, 1, 1),
				fakeDate(2016, 1, 2),
				"New Year's Day",
			),
		},
		{"2016-06-01	00:00	2016-07-14	00:00	Very long holiday",
			fakeHday(
				fakeDate(2016, 6, 1),
				fakeDate(2016, 7, 14),
				"Very long holiday",
			),
		},
	}
	for _, tt := range tests {
		got, err := NewHday(tt.in)
		if !reflect.DeepEqual(got, tt.want) || err != nil {
			t.Errorf("NewHday(%q) = %v, %v; want %v, nil",
				tt.in,
				got,
				err,
				tt.want,
			)
		}
	}
}

func TestNewHday_errors(t *testing.T) {

	tests := []struct {
		in string
	}{
		// wrong dates
		{"wrong	00:00	2016-01-02	00:00	Wrong start date"},
		{"2016-04-02	00:00	wrong	00:00	Wrong end date"},
		// fields !=5
		{"2016-01-01	00:00	2016-01-02	00:00"},
		// empty
		{""},
	}
	for _, tt := range tests {
		if _, err := NewHday(tt.in); err == nil {
			t.Errorf("NewHday(%q) = _, %v; want _, error",
				tt.in,
				err,
			)
		}
	}
}

func TestHday_String(t *testing.T) {
	h := fakeHday(
		fakeDate(2016, 1, 1),
		fakeDate(2016, 1, 2),
		"Fake Holiday")

	want := "2016-01-01 2016-01-02 Fake Holiday"
	if got := h.String(); got != want {
		t.Errorf("String() = %v; want %v",
			got,
			want)
	}
}
