package sql

import (
	"fmt"
	"strings"
	"time"
)

type Date time.Time

func (t Date) MarshalJSON() ([]byte, error) {
	if t.Time().IsZero() {
		return nil, nil
	}
	date := fmt.Sprintf("\"%s\"", t.Time().Format(time.DateOnly))
	return []byte(date), nil
}

func (t *Date) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == "{}" {
		return nil
	}
	time, err := time.Parse(time.DateOnly, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}
	*t = Date(time)
	return nil
}

func (t *Date) UnmarshalBinary(data []byte) error {
	if data == nil {
		return nil
	}
	time, err := time.Parse(time.DateOnly, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}
	*t = Date(time)
	return nil
}

func (t *Date) UnmarshalText(data []byte) error {
	if data == nil {
		return nil
	}
	time, err := time.Parse(time.DateOnly, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}
	*t = Date(time)
	return nil
}

func (t Date) Time() time.Time {
	return time.Time(t)
}

func (t Date) String() string {
	if t.Time().IsZero() {
		return ""
	}
	return t.Time().Format(time.DateOnly)
}

func (t Date) AddDate(years int, months int, days int) Date {
	return Date(t.Time().AddDate(years, months, days))
}
