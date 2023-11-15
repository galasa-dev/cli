/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"time"
)

type MockTimeService struct {
	MockNow time.Time
}

func NewMockTimeServiceAsMock(now time.Time) *MockTimeService {
	return &MockTimeService{MockNow: now}
}

func NewMockTimeService() TimeService {
	return NewMockTimeServiceAsMock(time.Now())
}

func NewOverridableMockTimeService(now time.Time) TimeService {
	return NewMockTimeServiceAsMock(now)
}

func (ts *MockTimeService) Interrupt(message string) {
	// The mock timing service doesn't know how to be interrupted.
}

func (ts *MockTimeService) Sleep(duration time.Duration) {
	// Do not sleep. Just advance the mock now time.
	ts.MockNow.Add(duration)
}

func (ts *MockTimeService) Now() time.Time {
	return ts.MockNow
}
