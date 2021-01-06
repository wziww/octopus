// +build linux

package internal

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
	secondsPerWeek   = 7 * secondsPerDay
	daysPer400Years  = 365*400 + 97
	daysPer100Years  = 365*100 + 24
	daysPer4Years    = 365*4 + 1
)
const (
	hasMonotonic         = 1 << 63
	maxWall              = wallToInternal + (1<<33 - 1) // year 2157
	minWall              = wallToInternal               // year 1885
	nsecMask             = 1<<30 - 1
	wallToInternal int64 = (1884*365 + 1884/4 - 1884/100 + 1884/400) * secondsPerDay
	nsecShift            = 30
)
const (
	unixToInternal int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
	internalToUnix int64 = -unixToInternal
)

func walltime() (sec int64, nsec int32)

type bTime struct {
	wall uint64
}

// Now ...
func Now() bTime {
	sec, nsec := walltime()
	sec += unixToInternal - minWall
	return bTime{
		hasMonotonic | uint64(sec)<<nsecShift | uint64(nsec),
	}
}
func (t bTime) UnixNano() int64 {
	return (t.unixSec())*1e9 + int64(t.nsec())
}
func (t *bTime) sec() int64 {
	if t.wall&hasMonotonic != 0 {
		return wallToInternal + int64(t.wall<<1>>(nsecShift+1))
	}
	panic("todo")
}
func (t *bTime) unixSec() int64 { return t.sec() + internalToUnix }
func (t *bTime) nsec() int32 {
	return int32(t.wall & nsecMask)
}
