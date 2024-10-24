/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"log"
	"time"

	"github.com/galasa-dev/cli/pkg/spi"
)

type timeService struct {
	interruptEventChannel chan string
}

func NewRealTimeService() spi.TimeService {
	service := timeService{
		// The interrupt channel has enough capacity for 100 events before anything blocks.
		interruptEventChannel: make(chan string, 100),
	}
	log.Printf("timeService: %v created\n", service)

	// service.logStackTrace()

	return &service
}

// func (ts *timeService) logStackTrace() {
// 	// Print the stack trace.
// 	buf := make([]byte, 1<<16)
// 	bytesInStackTrace := runtime.Stack(buf, true)
// 	log.Printf("timeService: %v stack trace : %s", *ts, buf[:bytesInStackTrace])
// }

const (
	SECONDS_BETWEEN_INTERRUPTED_CHECKS  = 3
	DURATION_BETWEEN_INTERRUPTED_CHECKS = SECONDS_BETWEEN_INTERRUPTED_CHECKS * time.Second
)

// Interrupts any timer sleeping.
func (ts *timeService) Interrupt(message string) {
	log.Printf("timeService: %v Interrupting the timing service sleeping. %s\n", *ts, message)
	// ts.logStackTrace()
	ts.interruptEventChannel <- "INTERRUPT: " + message
}

// Sleep for a bit. Waking up if anything calls the interrupt method.
func (ts *timeService) Sleep(duration time.Duration) {

	log.Printf("timeService: %v : sleep entered\n", *ts)
	timer := time.After(duration)

	select {
	case msg := <-ts.interruptEventChannel:
		log.Printf("timeService: %v : received interrupt message %s\n", *ts, msg)
	case <-timer:
		log.Printf("timeService: %v : sleep timed out\n", *ts)
		// ts.logStackTrace()
	}
	// log.Printf("timeService: %v : sleep exiting\n", *ts)
}

// Retrieves the current time, with the location set to UTC.
func (ts *timeService) Now() time.Time {
	return time.Now().UTC()
}
