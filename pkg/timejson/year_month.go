package timejson

import (
	"strings"
	"time"
)

type YearMonthTime time.Time

const YearMonthOnly = "2006-01"

func (y *YearMonthTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(YearMonthOnly, value)
	if err != nil {
		return err
	}

	*y = YearMonthTime(t)
	return nil
}

func (y YearMonthTime) MarshalJSON() ([]byte, error) {
	if time.Time(y).IsZero() {
		return nil, nil
	}
	return []byte(`"` + time.Time(y).Format(YearMonthOnly) + `"`), nil
}
