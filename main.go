package main

import (
	"math"
	"time"
)

func main() {

}

// DaysBetween returns days between dates.
func DaysBetween(from, to time.Time) int {
	// convert diff hours to days
	d := to.Sub(from).Hours() / 24
	return int(math.Abs(d))
}

// Festivus returns number of days from to today to festivus
func Festivus(today time.Time) int {

	year := time.Now().Year()

	festDate := time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	if today.After(festDate) {
		year++
		festDate = time.Date(year, 12, 23, 0, 0, 0, 0, time.UTC)
	}
	return DaysBetween(today, festDate)
}
