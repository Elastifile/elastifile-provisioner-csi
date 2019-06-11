package retry

import (
	"testing"
	"time"
)

func withDuration(dur time.Duration, t *testing.T) {
	s := Sigmoid{
		Limit:   dur / 10,
		Retries: 10,
	}
	sr := &sigmoidRetrier{
		step:  s.Retries / -2,
		upper: s.Limit,
	}
	sr.retries = s.Retries
	total := sr.totalTimeout()
	err := time.Millisecond
	if total-dur < err {
		t.Fatalf("Expected total %s - %s to be < %s", total, dur, err)
	}
}

func TestTotalTime(t *testing.T) {
	durations := []time.Duration{
		5 * time.Second,
		5 * time.Minute,
		5 * time.Hour,
	}
	for _, d := range durations {
		withDuration(d, t)
	}
}
