package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time().IsZero() {
		return []byte("{}"), nil
	}
	stamp := fmt.Sprintf("\"%s\"", t.Time().Format(time.DateTime))
	return []byte(stamp), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if data == nil || string(data) == "{}" {
		return nil
	}
	time, err := time.Parse(time.DateTime, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}
	*t = Time(time)
	return nil
}

func (t *Time) UnmarshalBinary(data []byte) error {
	if data == nil {
		return nil
	}
	time, err := time.Parse(time.DateTime, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}
	*t = Time(time)
	return nil
}

func (t *Time) UnmarshalText(data []byte) error {
	if data == nil {
		return nil
	}
	time, err := time.Parse(time.DateTime, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}
	*t = Time(time)
	return nil
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) String() string {
	if t.Time().IsZero() {
		return ""
	}
	return t.Time().Format(time.DateTime)
}

// Scan helper to retrieve duration data from postgres
func (t *Time) Scan(value interface{}) error {
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	case time.Time:
		*t = Time(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into time", value)
	}

	var parsed Time
	err := json.Unmarshal([]byte(s), &parsed)
	if err != nil {
		return fmt.Errorf("time.Parse(%q): %w", s, err)
	}
	*t = parsed
	return nil
}

// Value helper to insert duration data into postgres
func (t Time) Value() (driver.Value, error) {
	return t.String(), nil
}
