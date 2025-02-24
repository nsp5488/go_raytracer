package hittable

import (
	"math"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/util"
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

type translate struct {
	object Hittable
	offset *vec.Vec3
	bbox   *aabb.AABB
}

func Translate(object Hittable, offset *vec.Vec3) *translate {
	bbox := object.BBox().VecOffset(offset)
	return &translate{object: object, offset: offset, bbox: bbox}

}
func (t *translate) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	offsetRay := ray.NewWithTime(r.Origin().Sub(t.offset), r.Direction(), r.Time())

	if !t.object.Hit(offsetRay, rayT, record) {
		return false
	}

	record.p.AddInplace(t.offset)
	return true
}

func (t *translate) BBox() *aabb.AABB {
	return t.bbox
}

type rotateY struct {
	object   Hittable
	sinTheta float64
	cosTheta float64
	bbox     *aabb.AABB
}

func RotateY(object Hittable, theta float64) *rotateY {
	ry := &rotateY{object: object}
	radians := util.DegressToRadians(theta)
	ry.sinTheta = math.Sin(radians)
	ry.cosTheta = math.Cos(radians)
	bbox := object.BBox()

	min := [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)}
	max := [3]float64{math.Inf(-1), math.Inf(-1), math.Inf(-1)}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*bbox.AxisInterval(0).Max + float64(1-i)*bbox.AxisInterval(0).Min
				y := float64(j)*bbox.AxisInterval(1).Max + float64(1-j)*bbox.AxisInterval(1).Min
				z := float64(k)*bbox.AxisInterval(2).Max + float64(1-k)*bbox.AxisInterval(2).Min

				newX := ry.cosTheta*x + ry.sinTheta*z
				newZ := -ry.sinTheta*x + ry.cosTheta*z
				test := vec.New(newX, y, newZ)
				for c := 0; c < 3; c++ {
					min[c] = math.Min(min[c], test.Get(c))
					max[c] = math.Max(max[c], test.Get(c))
				}
			}
		}
	}
	ry.bbox = aabb.FromPoints(vec.New(min[0], min[1], min[2]), vec.New(max[0], max[1], max[2]))
	return ry
}

func (ry *rotateY) rayTranslationHelper(vector *vec.Vec3) *vec.Vec3 {
	return vec.New(
		ry.cosTheta*vector.X()-ry.sinTheta*vector.Z(),
		vector.Y(),
		ry.sinTheta*vector.X()+ry.cosTheta*vector.Z(),
	)
}

func (ry *rotateY) recordTranslationHelper(vector *vec.Vec3) *vec.Vec3 {
	return vec.New(
		ry.cosTheta*vector.X()+ry.sinTheta*vector.Z(),
		vector.Y(),
		-ry.sinTheta*vector.X()+ry.cosTheta*vector.Z(),
	)
}
func (ry *rotateY) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	origin := ry.rayTranslationHelper(r.Origin())
	direction := ry.rayTranslationHelper(r.Direction())

	rotatedRay := ray.NewWithTime(origin, direction, r.Time())

	if !ry.object.Hit(rotatedRay, rayT, record) {
		return false
	}
	record.p = ry.recordTranslationHelper(record.p)
	record.normal = ry.recordTranslationHelper(record.normal)

	return true
}
func (ry *rotateY) BBox() *aabb.AABB {
	return ry.bbox
}
