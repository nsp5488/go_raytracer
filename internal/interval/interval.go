package interval

import "math"

type Interval struct {
	Min float64
	Max float64
}

func New(min, max float64) *Interval {
	return &Interval{min, max}
}

func (i *Interval) Size() float64 {
	return i.Max - i.Min
}
func (i *Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}
func (i *Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}
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
