package hittable

import (
	"math"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Represents a sphere in 3D space
type Sphere struct {
	Center ray.Ray // Using a ray to represent motion

	Radius   float64
	Material Material
	bbox     *aabb.AABB
}

// Creates a new sphere
func NewSphere(center vec.Vec3, radius float64, material Material) *Sphere {
	rvec := vec.New(radius, radius, radius)
	bbox := aabb.FromPoints(center.Sub(rvec), center.Add(rvec))
	return &Sphere{Center: *ray.New(&center, vec.Empty()), Radius: radius, Material: material, bbox: bbox}
}

// Creates a new sphere with motion blur
func NewMotionSphere(center1, center2 vec.Vec3, radius float64, material Material) *Sphere {
	rvec := vec.New(radius, radius, radius)
	center := *ray.New(&center1, center2.Sub(&center1))
	bbox1 := aabb.FromPoints(center.At(0).Sub(rvec), center.At(0).Add(rvec))
	bbox2 := aabb.FromPoints(center.At(1).Sub(rvec), center.At(1).Add(rvec))

	return &Sphere{Center: center, Radius: radius, Material: material, bbox: aabb.FromBBoxes(bbox1, bbox2)}
}
func (s *Sphere) BBox() *aabb.AABB {
	return s.bbox
}

// Calculates the UV values of the ray intersection of a given sphere
// and stores them in (u, v)
func calculateSphereUV(point *vec.Vec3, u, v *float64) {
	theta := math.Acos(-point.Y())
	phi := math.Atan2(-point.Z(), point.X()) + math.Pi

	*u = phi / (2 * math.Pi)
	*v = theta / math.Pi
}

// Hit checks if a ray intersects with the sphere.
func (s *Sphere) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	curCenter := s.Center.At(r.Time())
	oc := curCenter.Sub(r.Origin())

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
	outward_normal := record.p.Sub(curCenter).Scale(1 / s.Radius)
	record.setFaceNormal(r, outward_normal)
	record.Material = s.Material
	calculateSphereUV(outward_normal, &record.u, &record.v)
	return true
}
