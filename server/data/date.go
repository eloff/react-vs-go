package data

import (
	"fmt"
	"time"
)

// Date a compact, naturally orderable numeric representation of a date
// packs year / month / day into a 32bit integer
type Date uint32

const (
	monthShift = 5
	yearShift  = monthShift + 4
	monthMask  = 0xf
	dayMask    = 0x1f

	// Time layout string for ShortName() Time.Format call
	shortNameLayout = "Mon Jan 2"
)

func NewDate(year int, month time.Month, day int) Date {
	return Date(uint32(year)<<yearShift | uint32(month)<<monthShift | uint32(day))
}

func DateFromTime(t time.Time) Date {
	return NewDate(t.Year(), t.Month(), t.Day())
}

func (date Date) After(other Date) bool {
	return date > other
}

func (date Date) Year() int {
	return int(date >> yearShift)
}

func (date Date) Month() time.Month {
	return time.Month((date >> monthShift) & monthMask)
}

func (date Date) Day() int {
	return int(date & dayMask)
}

func (date Date) Time() time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

func (date Date) Weekday() time.Weekday {
	return date.Time().Weekday()
}

func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(date.String()), nil
}

func (date Date) String() string {
	return fmt.Sprintf(`"%d-%d-%d"`, date.Year(), date.Month(), date.Day())
}

func (date Date) ShortName() string {
	return date.Time().Format(shortNameLayout)
}
