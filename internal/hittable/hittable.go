package hittable

import (
	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Records information about a ray hitting a surface (hittable)
type HitRecord struct {
	p         *vec.Vec3
	normal    *vec.Vec3
	t         float64
	frontFace bool

	u float64
	v float64

	Material Material
}

// Sets the face normal based on the ray direction and the normal vector
func (hr *HitRecord) setFaceNormal(r *ray.Ray, normal *vec.Vec3) {
	hr.frontFace = r.Direction().Dot(normal) < 0
	if hr.frontFace {
		hr.normal = normal
	} else {
		hr.normal = normal.Negate()
	}
}

// Returns the normal vector of the hit surface
func (hr *HitRecord) Normal() *vec.Vec3 {
	return hr.normal
}

// Returns the point of intersection of the ray with the surface
func (hr *HitRecord) P() *vec.Vec3 {
	return hr.p
}

// Returns whether the ray hit the surface from the front or back
func (hr *HitRecord) FrontFace() bool {
	return hr.frontFace
}

// Defines the behavior of a hittable object
type Hittable interface {
	Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool
	BBox() *aabb.AABB
}

// A container struct for a list of hittable objects. Effectively a scene.
type HittableList struct {
	objects []Hittable

	bbox *aabb.AABB
}

func NewHittableList(objects []Hittable) *HittableList {
	hl := &HittableList{}
	hl.Init(len(objects))
	copy(hl.objects, objects)
	return hl
}

func (hl *HittableList) Init(startSize int) {
	hl.objects = make([]Hittable, 0, startSize)
	hl.bbox = aabb.EmptyBBox()
}
func (hl *HittableList) Clear() {
	hl.objects = make([]Hittable, 10)
}

func (hl *HittableList) Add(obj Hittable) {
	hl.objects = append(hl.objects, obj)
	hl.bbox = aabb.FromBBoxes(hl.bbox, obj.BBox())
}
func (hl *HittableList) BBox() *aabb.AABB {
	return hl.bbox
}

// Checks if a ray hits any of the objects in a scene
func (hl *HittableList) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	hitRecord := HitRecord{}
	tmp := &HitRecord{}
	hitAny := false
	closest_so_far := rayT.Max
	interval := interval.New(rayT.Min, closest_so_far)

	for _, obj := range hl.objects {
		if obj.Hit(r, *interval, &hitRecord) {
			hitAny = true
			closest_so_far = hitRecord.t
			interval.Max = closest_so_far
			tmp = &hitRecord
		}
	}
	*record = *tmp
	return hitAny
}
