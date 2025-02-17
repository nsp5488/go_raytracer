package hittable

import (
	"math"

	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type Sphere struct {
	Center   vec.Vec3
	Radius   float64
	Material Material
}

func (s *Sphere) Init(center vec.Vec3, radius float64, material Material) {
	s.Center = center
	s.Radius = radius
	s.Material = material
}
func (s *Sphere) Hit(r *ray.Ray, rayT *interval.Interval, record *HitRecord) bool {
	oc := s.Center.Add(r.Origin().Negate())

	// a = direction * direction = len(direction)^2
	a := r.Direction().LengthSquared()
	// h = r * oc
	h := r.Direction().Dot(oc)
	// c = oc * oc - radius^2 = len(oc)^2 - radius^2
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	record.t = root
	record.p = r.At(root)
	outward_normal := record.p.Add(s.Center.Negate()).Scale(1 / s.Radius)
	record.setFaceNormal(r, outward_normal)
	record.Material = s.Material
	return true
}
