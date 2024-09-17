package fininf

import (
	"errors"
	"strings"
	"time"
)

type Interval struct {
	Duration time.Duration
}

func (i *Interval) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return err
	}

	var t time.Duration
	switch strings.ToLower(buf) {
	case "daily":
		t = 24 * time.Hour
	case "weekly":
		t = 7 * 24 * time.Hour
	case "fortnightly", "bi-weekly":
		t = 14 * 24 * time.Hour
	case "4-weekly":
		t = 4 * 7 * 24 * time.Hour
	case "monthly":
		t = (365.25 * 24 * time.Hour) / 12
	case "quarterly":
		t = (365.25 * 24 * time.Hour) / 4
	case "yearly", "annually":
		t = 365.25 * 24 * time.Hour
	default:
		return errors.New("invalid interval")
	}

	i.Duration = t
	return nil
}
