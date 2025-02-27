package hittable

import (
	"log"
	"math/rand"

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

func (hr *HitRecord) U() float64 {
	return hr.u
}

func (hr *HitRecord) V() float64 {
	return hr.v
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
	PdfValue(origin, direction *vec.Vec3) float64
	Random(origin *vec.Vec3) *vec.Vec3
}

type defaultPdfImpl struct{}

func (d defaultPdfImpl) PdfValue(origin, direction *vec.Vec3) float64 {
	log.Fatal("hit an invalid PDF function")
	return 0.0
}
func (d defaultPdfImpl) Random(origin *vec.Vec3) *vec.Vec3 {
	return vec.New(1, 0, 0)
}

// A container struct for a list of hittable objects. Effectively a scene.
type hittableList struct {
	objects []Hittable

	bbox *aabb.AABB
}

func NewHittableList(startSize int) *hittableList {
	hl := &hittableList{}
	hl.init(startSize)
	return hl
}
func (hl *hittableList) PdfValue(origin, direction *vec.Vec3) float64 {
	weight := 1.0 / float64(len(hl.objects))
	sum := 0.0
	for _, obj := range hl.objects {
		sum += weight * obj.PdfValue(origin, direction)
	}
	return sum

}
func (hl *hittableList) Random(origin *vec.Vec3) *vec.Vec3 {
	if len(hl.objects) <= 0 {
		return vec.Random()
	}
	return hl.objects[rand.Intn(len(hl.objects))].Random(origin)
}

func (hl *hittableList) init(startSize int) {
	hl.objects = make([]Hittable, 0, startSize)
	hl.bbox = aabb.EmptyBBox()
}
func (hl *hittableList) Clear() {
	hl.objects = make([]Hittable, 10)
}

func (hl *hittableList) Add(obj Hittable) {
	hl.objects = append(hl.objects, obj)
	hl.bbox = aabb.FromBBoxes(hl.bbox, obj.BBox())
}
func (hl *hittableList) BBox() *aabb.AABB {
	return hl.bbox
}

// Checks if a ray hits any of the objects in a scene
func (hl *hittableList) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	hitRecord := &HitRecord{}
	// tmp := &HitRecord{}
	hitAny := false
	closest_so_far := rayT.Max
	interval := interval.New(rayT.Min, closest_so_far)

	for _, obj := range hl.objects {
		if obj.Hit(r, *interval, hitRecord) {
			hitAny = true
			closest_so_far = hitRecord.t
			interval.Max = closest_so_far
			*record = *hitRecord
		}
	}
	return hitAny
}
