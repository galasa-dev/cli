/*
 * Copyright contributors to the Galasa project
 */
package runs

const (
	TIME_UNIT_WEEKS_LONG   = "weeks"
	TIME_UNIT_DAYS_LONG    = "days"
	TIME_UNIT_HOURS_LONG   = "hours"
	TIME_UNIT_MINUTES_LONG = "minutes"

	TIME_UNIT_WEEKS_SHORT   = "w"
	TIME_UNIT_DAYS_SHORT    = "d"
	TIME_UNIT_HOURS_SHORT   = "h"
	TIME_UNIT_MINUTES_SHORT = "m"
)

type TimeUnit struct {
	name             string
	minuteMultiplier int
}

func newTimeUnit(name string, multiplier int) *TimeUnit {
	return &TimeUnit{name: name, minuteMultiplier: multiplier}
}

func (unit *TimeUnit) getName() string {
	return unit.name
}

func (unit *TimeUnit) getMinuteMultiplier() int {
	return unit.minuteMultiplier
}
