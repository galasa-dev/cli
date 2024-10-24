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

// Retrieves the current time, with the location set to UTC.
func (ts *timeService) Now() time.Time {
	return time.Now().UTC()
}
