package types

import (
	"encoding/json"
	"errors"
	"time"
)

type CustomTime time.Time

const (
	timeFormat = time.RFC3339
)

func (ct CustomTime) String() string {
	return time.Time(ct).Format(timeFormat)
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
		t, err = time.Parse(timeFormat, string(src))
	case string:
		t, err = time.Parse(timeFormat, src)
	default:
		return errors.New("invalid type")
	}

	if err != nil {
		return errors.New("invalid time format")
	}

	*ct = CustomTime(t)
	return nil
}

func (ct CustomTime) Value() (interface{}, error) {
	return ct.String(), nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	val := time.Time(ct).Format(timeFormat)
	return json.Marshal(val)
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(timeFormat, string(b))
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}
