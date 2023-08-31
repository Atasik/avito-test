package timejson

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CivilTime struct {
	time.Time
}

func (t *CivilTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" || s == "" {
		return
	}
	t.Time, err = time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	return nil
}

func (t CivilTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(time.DateOnly))), nil
}

func (t *CivilTime) Scan(v interface{}) error {
	ti := v.(time.Time)
	*t = CivilTime{Time: ti}
	return nil
}

func (t CivilTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
