package date

import (
	"time"
)

type intervalGap struct {
	isOn bool
	last time.Time
	next time.Time
}

type IntervalStep struct {
	reference time.Time
	interval  time.Duration

	from  time.Time
	cache map[time.Time]*intervalGap
}

var intervalSteps = map[time.Time]map[time.Duration]*IntervalStep{}

func NewIntervalStep(reference time.Time, interval time.Duration, from time.Time, to time.Time) *IntervalStep {
	if interval == 0 {
		panic("interval cannot be zero")
	}

	is := &IntervalStep{
		reference: Date(reference),
		interval:  interval,
		from:      Date(from),
		cache:     buildIntervalCache(reference, interval, from, to),
	}

	return is
}

func buildIntervalCache(reference time.Time, interval time.Duration, from time.Time, to time.Time) map[time.Time]*intervalGap {
	cache := map[time.Time]*intervalGap{}

	if reference == (time.Time{}) {
		panic("reference cannot be zero")
	}
	if interval == 0 {
		panic("interval cannot be zero")
	}

	for d := Date(from); d.Compare(to) <= 0; d = Date(d.AddDate(0, 0, 1)) {
		cache[d] = &intervalGap{
			isOn: false,
		}
	}

	for d := Date(reference); d.Compare(from) >= 0; d = Date(d.Add(0 - interval)) {
		cache[d].isOn = true
	}
	for d := Date(reference); d.Compare(to) <= 0; d = Date(d.Add(interval)) {
		cache[d].isOn = true
	}

	lastD := time.Time{}
	for d := Date(from); d.Compare(to) <= 0; d = Date(d.AddDate(0, 0, 1)) {
		cache[d].last = lastD
		if cache[d].isOn {
			lastD = d
		}
	}
	lastD = time.Time{}
	for d := Date(to); d.Compare(from) >= 0; d = Date(d.AddDate(0, 0, -1)) {
		cache[d].next = lastD
		if cache[d].isOn {
			lastD = d
		}
	}

	return cache
}

func (s *IntervalStep) IsOn(d time.Time) bool {
	return s.cache[Date(d)].isOn
}

func (s *IntervalStep) FirstOnOrAfter(d time.Time) time.Time {
	d = Date(d)
	if s.IsOn(d) {
		return d
	}
	return s.FirstAfter(d)
}

func (s *IntervalStep) FirstAfter(d time.Time) time.Time {
	return s.cache[Date(d)].next
}
