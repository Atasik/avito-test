package timejson

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type YearMonthTime struct {
	time.Time
}

const YearMonthOnly = "2006-01"

func (t *YearMonthTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`) // remove quotes
	if s == "null" || s == "" {
		return
	}
	t.Time, err = time.Parse(YearMonthOnly, s)
	if err != nil {
		return err
	}
	return nil
}

func (t YearMonthTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(YearMonthOnly))), nil
}

func (t *YearMonthTime) Scan(v interface{}) error {
	ti, ok := v.(time.Time)
	if !ok {
		return errors.New("invalid time format")
	}
	*t = YearMonthTime{Time: ti}
	return nil
}

func (t YearMonthTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
