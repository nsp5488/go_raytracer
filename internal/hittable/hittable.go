package hittable

import (
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type HitRecord struct {
	p         *vec.Vec3
	normal    *vec.Vec3
	t         float64
	frontFace bool
	Material  Material
}

func (hr *HitRecord) setFaceNormal(r *ray.Ray, normal *vec.Vec3) {
	hr.frontFace = r.Direction().Dot(normal) < 0
	if hr.frontFace {
		hr.normal = normal
	} else {
		hr.normal = normal.Negate()
	}
}
func (hr *HitRecord) Normal() *vec.Vec3 {
	return hr.normal
}

func (hr *HitRecord) P() *vec.Vec3 {
	return hr.p
}
func (hr *HitRecord) FrontFace() bool {
	return hr.frontFace
}

type Hittable interface {
	Hit(r *ray.Ray, rayT *interval.Interval, record *HitRecord) bool
}

type HittableList struct {
	objects []Hittable
}

func (hl *HittableList) Init(startSize int) {
	hl.objects = make([]Hittable, 0, startSize)
}
func (hl *HittableList) Clear() {
	hl.objects = make([]Hittable, 10)
}

func (hl *HittableList) Add(obj Hittable) {
	hl.objects = append(hl.objects, obj)
}

func (hl *HittableList) Hit(r *ray.Ray, rayT *interval.Interval, record *HitRecord) bool {
	hitRecord := HitRecord{}
	tmp := &HitRecord{}
	hitAny := false
	closest_so_far := rayT.Max
	interval := interval.New(rayT.Min, closest_so_far)

	for _, obj := range hl.objects {
		if obj.Hit(r, interval, &hitRecord) {
			hitAny = true
			closest_so_far = hitRecord.t
			interval.Max = closest_so_far
			tmp = &hitRecord
		}
	}
	*record = *tmp
	return hitAny
}
