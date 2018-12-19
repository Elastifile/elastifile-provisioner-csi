package etime

// TODO: Move this to infra/optional.

import "time"

type NilableTime *time.Time

func NewNilableTime(t time.Time) NilableTime {
	return &t
}
