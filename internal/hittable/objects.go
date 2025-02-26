package hittable

import (
	"math"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Represents a sphere in 3D space
type sphere struct {
	Center *ray.Ray // Using a ray to represent motion

	Radius   float64
	Material Material
	bbox     *aabb.AABB
}

// Creates a new sphere
func NewSphere(center *vec.Vec3, radius float64, material Material) *sphere {
	rvec := vec.New(radius, radius, radius)
	bbox := aabb.FromPoints(center.Sub(rvec), center.Add(rvec))
	return &sphere{Center: ray.New(center, vec.Empty()), Radius: radius, Material: material, bbox: bbox}
}

// Creates a new sphere with motion blur
func NewMotionSphere(center1, center2 *vec.Vec3, radius float64, material Material) *sphere {
	rvec := vec.New(radius, radius, radius)
	center := *ray.New(center1, center2.Sub(center1))
	bbox1 := aabb.FromPoints(center.At(0).Sub(rvec), center.At(0).Add(rvec))
	bbox2 := aabb.FromPoints(center.At(1).Sub(rvec), center.At(1).Add(rvec))

	return &sphere{Center: &center, Radius: radius, Material: material, bbox: aabb.FromBBoxes(bbox1, bbox2)}
}
func (s *sphere) BBox() *aabb.AABB {
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
func (s *sphere) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
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

type quad struct {
	Q      *vec.Vec3 // One corner of the plane
	u      *vec.Vec3 // u,v are vectors that point from Q to two other corners
	v      *vec.Vec3
	normal *vec.Vec3 // normal = unit(u x v)
	w      *vec.Vec3
	D      float64 // D = Ax + By + Cz = dot(Q, normal)

	bbox     *aabb.AABB
	material Material
}

func NewQuad(Q, u, v *vec.Vec3, material Material) *quad {
	q := &quad{Q: Q, u: u, v: v, material: material}

	n := u.Cross(v)

	q.normal = n.UnitVector()
	q.D = q.normal.Dot(Q)
	q.w = n.Scale(1 / n.Dot(n))

	q.setBBox()
	return q
}

func (q *quad) setBBox() {
	diag1 := aabb.FromPoints(q.Q, q.Q.Add(q.u).Add(q.v))
	diag2 := aabb.FromPoints(q.Q.Add(q.u), q.Q.Add(q.v))
	q.bbox = aabb.FromBBoxes(diag1, diag2)
}

func (q *quad) BBox() *aabb.AABB {
	return q.bbox
}
func (q *quad) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	denom := q.normal.Dot(r.Direction())

	// Low values in denominator -> ray is parallel to the plane
	if math.Abs(denom) < 1e-8 {
		return false
	}

	// If t lands outside of our search interval
	t := (q.D - q.normal.Dot(r.Origin())) / denom
	if !rayT.Contains(t) {
		return false
	}

	intersection := r.At(t)

	// Check that the ray intersects the quad itself, not just the plane
	planarHitpoint := intersection.Sub(q.Q)
	alpha := q.w.Dot(planarHitpoint.Cross(q.v))
	beta := q.w.Dot(q.u.Cross(planarHitpoint))
	if !isInterior(alpha, beta, record) {
		return false
	}

	record.t = t
	record.p = intersection
	record.Material = q.material
	record.setFaceNormal(r, q.normal)
	return true
}

func isInterior(alpha, beta float64, record *HitRecord) bool {
	if !interval.Unit().Contains(alpha) || !interval.Unit().Contains(beta) {
		return false
	}

	record.u = alpha
	record.v = beta
	return true
}

func NewBox(a, b *vec.Vec3, mat Material) Hittable {
	sides := NewHittableList(6)

	minVec := vec.New(
		min(a.X(), b.X()),
		min(a.Y(), b.Y()),
		min(a.Z(), b.Z()),
	)
	maxVec := vec.New(
		max(a.X(), b.X()),
		max(a.Y(), b.Y()),
		max(a.Z(), b.Z()),
	)

	dx := vec.New(maxVec.X()-minVec.X(), 0, 0)
	dy := vec.New(0, maxVec.Y()-minVec.Y(), 0)
	dz := vec.New(0, 0, maxVec.Z()-minVec.Z())

	// front
	sides.Add(NewQuad(vec.New(minVec.X(), minVec.Y(), maxVec.Z()), dx, dy, mat))
	// right
	sides.Add(NewQuad(vec.New(maxVec.X(), minVec.Y(), maxVec.Z()), dz.Negate(), dy, mat))
	// back
	sides.Add(NewQuad(vec.New(maxVec.X(), minVec.Y(), minVec.Z()), dx.Negate(), dy, mat))
	// left
	sides.Add(NewQuad(vec.New(minVec.X(), minVec.Y(), minVec.Z()), dz, dy, mat))
	// top
	sides.Add(NewQuad(vec.New(minVec.X(), maxVec.Y(), maxVec.Z()), dx, dz.Negate(), mat))
	// bottom
	sides.Add(NewQuad(vec.New(minVec.X(), minVec.Y(), minVec.Z()), dx, dz, mat))

	return BuildBVH(sides)
}
