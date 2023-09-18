/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"log"
	"time"
)

type TimeService interface {
	Sleep(duration time.Duration)
	Now() time.Time
	Interrupt(message string)
}

type timeService struct {
	interruptEventChannel chan string
}

func NewRealTimeService() TimeService {
	service := timeService{
		// The interrupt channel has enough capacity for 100 events before anything blocks.
		interruptEventChannel: make(chan string, 100),
	}
	return &service
}

const (
	SECONDS_BETWEEN_INTERRUPTED_CHECKS  = 3
	DURATION_BETWEEN_INTERRUPTED_CHECKS = SECONDS_BETWEEN_INTERRUPTED_CHECKS * time.Second
)

// Interrupts any timer sleeping.
func (ts *timeService) Interrupt(message string) {
	log.Printf("Interrupting the timing service sleeping. %s\n", message)
	ts.interruptEventChannel <- "INTERRUPT: " + message
}

// Sleep for a bit. Waking up if anything calls the interrupt method.
func (ts *timeService) Sleep(duration time.Duration) {

	// Clear any interruptions which occurred before we went to sleep.
	for ts.hasBeenInterrupted() {
		// Do nothing here. Checking is enough to clear an interruption.
	}

	isInterrupted := false
	isDone := false
	for !isDone {

		if duration < DURATION_BETWEEN_INTERRUPTED_CHECKS {
			// Only a bit of time left to sleep for. Do it.
			time.Sleep(duration)
			duration = 0
		} else {

			// Have we been interrupted ?
			isInterrupted = ts.hasBeenInterrupted()

			if !isInterrupted {
				time.Sleep(DURATION_BETWEEN_INTERRUPTED_CHECKS)
				duration -= DURATION_BETWEEN_INTERRUPTED_CHECKS
			}
		}

		if duration <= 0 || isInterrupted {
			isDone = true
		}
	}
}

// Check to see if anything has interrupted the timer service.
// Called by a sleeping go routine, to check between polls.
func (ts *timeService) hasBeenInterrupted() bool {
	isInterrupted := false
	select {
	case msg := <-ts.interruptEventChannel:
		log.Printf("TimeService: received interrupt message %s\n", msg)
		// There was an interrupt message received. So don't sleep for any longer.
		isInterrupted = true
	default:
		// log.Printf("no interrupt message received\n")
	}
	return isInterrupted
}

// Retrieves the current time, with the location set to UTC.
func (ts *timeService) Now() time.Time {
	return time.Now().UTC()
}
