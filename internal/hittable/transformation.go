package hittable

import (
	"math"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/util"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

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
