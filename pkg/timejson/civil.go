package timejson

import (
	"strings"
	"time"
)

type CivilTime time.Time

const YMS = "2006-01-02"

func (c *CivilTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(YMS, value)
	if err != nil {
		return err
	}

	*c = CivilTime(t)
	return nil
}

func (c CivilTime) MarshalJSON() ([]byte, error) {
	if time.Time(c).IsZero() {
		return nil, nil
	}
	return []byte(`"` + time.Time(c).Format(YMS) + `"`), nil
}
