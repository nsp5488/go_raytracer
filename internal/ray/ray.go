package ray

import (
	"fmt"

	"github.com/nsp5488/go_raytracer/internal/vec"
)

// A ray in 3D space
type Ray struct {
	origin    *vec.Vec3
	direction *vec.Vec3
	time      float64
}

// Creates a new ray with the given origin and direction
func New(origin, direction *vec.Vec3) *Ray {
	return &Ray{origin: origin, direction: direction, time: 0}
}
func NewWithTime(origin, direction *vec.Vec3, time float64) *Ray {
	return &Ray{origin: origin, direction: direction, time: time}
}
func (r *Ray) Origin() *vec.Vec3 {
	return r.origin
}

func (r *Ray) Direction() *vec.Vec3 {
	return r.direction
}
func (r *Ray) Time() float64 {
	return r.time
}

// Returns the point at time t along the ray
func (r *Ray) At(t float64) *vec.Vec3 {
	// P(t) = A + Bt, let A = origin, B = direction, t = time
	return r.origin.Add(r.direction.Scale(t))
}

func (r *Ray) String() string {
	return fmt.Sprintf("%s + %st", r.origin, r.direction)
}
