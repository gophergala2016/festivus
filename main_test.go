package main

import (
	"testing"
	"time"
)

func TestDaysBetween(t *testing.T) {

	since, err := time.Parse("02.01.2006.", "01.12.2016.")
	if err != nil {
		t.Fatalf("Date conversion failed with err %q", err)
	}

	to, err := time.Parse("02.01.2006.", "23.12.2016.")
	if err != nil {
		t.Fatalf("Date conversion failed with err %q", err)
	}

	want := 22
	got := DaysBetween(since, to)

	if got != want {
		t.Errorf("DaysBetween(%v, %v) = %v; want %v", since, to, got, want)
	}

}

func TestFestivus(t *testing.T) {

	today, err := time.Parse("02.01.2006", "01.12.2016")
	if err != nil {
		t.Fatalf("Date conversion failed with err %q", err)
	}

	want := 22
	got := Festivus(today)

	if got != want {
		t.Errorf("Festivus (today: %v) = %v; want %v", today, got, want)
	}
}
