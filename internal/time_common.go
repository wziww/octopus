// +build !linux

package internal

import "time"

// Now ...
func Now() time.Time {
	return time.Now()
}
