/*
 * Copyright contributors to the Galasa project
 */

package utils

import (
	"time"
)

type TimeService interface {
	Sleep(time.Duration)
	Now() time.Time
}

type timeService struct {
}

func NewRealTimeService() TimeService {
	return &timeService{}
}

func (ts *timeService) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

func (ts *timeService) Now() time.Time {
	return time.Now()
}
