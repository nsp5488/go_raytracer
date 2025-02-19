package aabb

import (
	"fmt"

	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Axis Aligned Bounding Box
type AABB struct {
	x *interval.Interval
	y *interval.Interval
	z *interval.Interval
}

func EmptyBBox() *AABB {
	return NewAABB(&interval.EMPTY, &interval.EMPTY, &interval.EMPTY)
}

// Constructs a new AABB
func NewAABB(x, y, z *interval.Interval) *AABB {
	return &AABB{x: x, y: y, z: z}
}

// Constructs an AABB from two points
func FromPoints(a, b *vec.Vec3) *AABB {
	var x, y, z *interval.Interval

	if a.X() < b.X() {
		x = interval.New(a.X(), b.X())
	} else {
		x = interval.New(b.X(), a.X())
	}
	if a.Y() < b.Y() {
		y = interval.New(a.Y(), b.Y())
	} else {
		y = interval.New(b.Y(), a.Y())
	}
	if a.Z() < b.Z() {
		z = interval.New(a.Z(), b.Z())
	} else {
		z = interval.New(b.Z(), a.Z())
	}

	return NewAABB(x, y, z)
}
func FromBBoxes(a, b *AABB) *AABB {
	x := interval.Combine(a.x, b.x)
	y := interval.Combine(a.y, b.y)
	z := interval.Combine(a.z, b.z)
	return NewAABB(x, y, z)
}

func (bb *AABB) AxisInterval(n int) *interval.Interval {
	if n == 2 {
		return bb.z
	}
	if n == 1 {
		return bb.y
	}
	return bb.x
}
func (bb *AABB) LongestAxis() int {
	if bb.x.Size() > bb.y.Size() {
		if bb.x.Size() > bb.z.Size() {
			return 0
		} else {
			return 2
		}
	} else {
		if bb.y.Size() > bb.z.Size() {
			return 1
		} else {
			return 2
		}
	}
}

func (bb *AABB) Hit(r *ray.Ray, rayT interval.Interval) bool {
	direction := r.Direction()
	origin := r.Origin()

	for axis := 0; axis < 3; axis++ {
		ax := bb.AxisInterval(axis)
		invD := 1 / direction.Get(axis)

		t0 := (ax.Min - origin.Get(axis)) * invD
		t1 := (ax.Max - origin.Get(axis)) * invD

		if invD < 0 {
			t0, t1 = t1, t0
		}
		rayT.Min = max(t0, rayT.Min)
		rayT.Max = min(t1, rayT.Max)

		if rayT.Max <= rayT.Min {
			return false
		}
	}

	return true
}
func (bb *AABB) String() string {
	return fmt.Sprintf("X: (%f,%f), Y: (%f,%f), Z: (%f, %f)", bb.x.Min, bb.x.Max, bb.y.Min, bb.y.Max, bb.z.Min, bb.z.Max)
}
