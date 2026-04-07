package scalar

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Time time.Time

func (t Time) MarshalGQL(w io.Writer) {
	timeStr := time.Time(t).Format(time.RFC3339)
	io.WriteString(w, strconv.Quote(timeStr))
}

func (t *Time) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("time must be a string")
	}
	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}
	*t = Time(parsed)
	return nil
}
