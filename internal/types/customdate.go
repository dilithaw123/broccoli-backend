package types

import (
	"encoding/json"
	"errors"
	"time"
)

type CustomTime time.Time

func (ct CustomTime) String() string {
	return time.Time(ct).Format(time.DateOnly)
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime(time.Time{})
		return nil
	}

	var t time.Time
	var err error

	switch src := value.(type) {
	case time.Time:
		t = src
	case []byte:
		t, err = time.Parse(time.DateOnly, string(src))
	case string:
		if len(src) < 10 {
			return errors.New("invalid time format")
		}
		t, err = time.Parse(time.DateOnly, src[:10])
	default:
		return errors.New("invalid type")
	}

	if err != nil {
		return errors.New("invalid time format")
	}

	*ct = CustomTime(t)
	return nil
}

func (ct *CustomTime) Value() (interface{}, error) {
	return time.Time(*ct), nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	val := time.Time(ct).Format(time.DateOnly)
	return json.Marshal(val)
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(time.DateOnly, string(b))
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}
