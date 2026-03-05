package clock

import "time"

// RealClock implements port.Clock using the system clock.
type RealClock struct{}

// Now returns the current UTC time.
func (RealClock) Now() time.Time {
	return time.Now().UTC()
}
