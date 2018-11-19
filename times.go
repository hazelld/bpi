package bpi

import (
	"time"
)

func CompareDay(t1 time.Time, t2 time.Time) bool {
	return (t1.Year() == t2.Year() &&
		t1.Month() == t2.Month() &&
		t1.Day() == t2.Day())
}
