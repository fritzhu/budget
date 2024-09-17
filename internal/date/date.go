package date

import (
	"time"

	"gopkg.in/yaml.v3"
)

type DateValue struct {
	Date time.Time
}

func Today() time.Time {
	return Date(time.Now().Local())
}

func (d *DateValue) UnmarshalYAML(node *yaml.Node) error {
	dateString := ""
	err := node.Decode(&dateString)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("2006-01-02", dateString, time.Local)
	if err != nil {
		return err
	}

	d.Date = Date(t)
	return nil
}

func Date(t time.Time) time.Time {
	y, m, da := t.Date()
	return time.Date(y, m, da, 12, 0, 0, 0, time.Local)
}

func EarliestDateAfter(after time.Time, reference time.Time, interval time.Duration) time.Time {
	d := Date(reference)
	for d.After(after) {
		d = Date(d.Add(-interval))
	}
	for d.Before(after) {
		d = Date(d.Add(interval))
	}
	return d
}
