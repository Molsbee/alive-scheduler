package parser

import "time"

// Schedule is a struct representation of the cron schedule
// Example: */15 * * * * *
type Schedule struct {
	Minute     uint64
	Hour       uint64
	DayOfMonth uint64
	Month      uint64
	DayOfWeek  uint64
}

// Next determines the next time a schedule should execute
func (s *Schedule) Next(t time.Time) time.Time {
	t = t.Add(1 * time.Second)

	for 1<<uint(t.Month())&s.Month == 0 {
		t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		t = t.AddDate(0, 1, 0)
		return s.Next(t)
	}

	for !dayMatches(s, t) {
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		t = t.AddDate(0, 0, 1)
		return s.Next(t)
	}

	for 1<<uint(t.Hour())&s.Hour == 0 {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		t = t.Add(1 * time.Hour)
		return s.Next(t)
	}

	for 1<<uint(t.Minute())&s.Minute == 0 {
		t = t.Add(1 * time.Minute)
		return s.Next(t)
	}

	return t
}

// dayMatches returns true if the schedule's day-of-week and day-of-month
// restrictions are satisfied by the given time.
func dayMatches(s *Schedule, t time.Time) bool {
	var (
		domMatch bool = 1<<uint(t.Day())&s.DayOfMonth > 0
		dowMatch bool = 1<<uint(t.Weekday())&s.DayOfWeek > 0
	)

	if s.DayOfMonth > 0 || s.DayOfWeek > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}
