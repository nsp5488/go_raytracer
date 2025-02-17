package ray

import (
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type Ray struct {
	origin    *vec.Vec3
	direction *vec.Vec3
}

func New(origin, direction *vec.Vec3) *Ray {
	return &Ray{origin: origin, direction: direction}
}
func (r *Ray) Origin() *vec.Vec3 {
	return r.origin
}
func (r *Ray) Direction() *vec.Vec3 {
	return r.direction
}

func (r *Ray) At(t float64) *vec.Vec3 {
	// P(t) = A + Bt, let A = origin, B = direction, t = time
	return r.origin.Add(r.direction.Scale(t))
}
