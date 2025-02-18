package interval

import "math"

// Interval represents a closed interval [min, max]
type Interval struct {
	Min float64
	Max float64
}

// New creates a new interval with the given min and max values
func New(min, max float64) *Interval {
	return &Interval{min, max}
}

// Size returns the size of the interval
func (i *Interval) Size() float64 {
	return i.Max - i.Min
}

// Contains checks if the interval contains the given value
func (i *Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

// Surrounds checks if the interval surrounds the given value (the value is strictly inside the interval)
func (i *Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

// Clamp clamps the given value to the interval
func (i *Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

var EMPTY = Interval{math.Inf(1), math.Inf(-1)}
var UNIVERSE = Interval{math.Inf(-1), math.Inf(1)}

func Empty() Interval {
	return EMPTY
}
func Universe() Interval {
	return UNIVERSE
}
