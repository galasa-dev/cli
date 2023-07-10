/*
 * Copyright contributors to the Galasa project
 */
package runs

const (
	TIME_UNIT_WEEKS   = "weeks"
	TIME_UNIT_DAYS    = "days"
	TIME_UNIT_HOURS   = "hours"
	TIME_UNIT_MINUTES = "minutes"
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
