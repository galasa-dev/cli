/*
 * Copyright contributors to the Galasa project
 */
package runs

type TimeUnit struct {
	name             string
	minuteMultiplier int
}

func (unit *TimeUnit) getName() string {
	return unit.name
}

func (unit *TimeUnit) getMinuteMultiplier() int {
	return unit.minuteMultiplier
}
