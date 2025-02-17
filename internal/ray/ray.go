package ray

import "github.com/nsp5488/go_raytracer/internal/vec"

type Ray struct {
	origin    *vec.Vec3
	direction *vec.Vec3
}

func New(origin, direction *vec.Vec3) *Ray {
	return &Ray{origin: origin, direction: direction}
}
func (r *Ray) Origin() vec.Vec3 {
	return *r.origin
}
func (r *Ray) Direction() vec.Vec3 {
	return *r.direction
}

func (r *Ray) At(t float64) vec.Vec3 {
	// P(t) = A + Bt, let A = origin, B = direction, t = time
	return *r.origin.Add(r.direction.Scale(t))
}

func (r *Ray) Color() *vec.Vec3 {
	unitDirection := r.direction.UnitVector()
	a := 0.5 * (unitDirection.Y() + 1.0)
	c1 := vec.New(1, 1, 1).Scale(1.0 - a)
	c2 := vec.New(0.5, 0.7, 1).Scale(a)
	return c1.Add(c2)
}
