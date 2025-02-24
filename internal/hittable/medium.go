package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type constantMedium struct {
	boundary               Hittable
	negativeInverseDensity float64
	phaseFunction          Material
}

func ConstantMediumTexture(boundary Hittable, density float64, tex Texture) *constantMedium {
	return &constantMedium{boundary: boundary, negativeInverseDensity: -1 / density, phaseFunction: NewIsotropicTexture(tex)}
}
func ConstantMedium(boundary Hittable, density float64, albedo *vec.Vec3) *constantMedium {
	return &constantMedium{boundary: boundary, negativeInverseDensity: -1 / density, phaseFunction: NewIsotropic(albedo)}
}

func (cm *constantMedium) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	hr1 := &HitRecord{}
	if !cm.boundary.Hit(r, *interval.Universe(), hr1) {
		return false
	}

	hr2 := &HitRecord{}
	if !cm.boundary.Hit(r, *interval.New(hr1.t+.0001, math.Inf(1)), hr2) {
		return false
	}
	hr1.t = max(hr1.t, rayT.Min)
	hr2.t = min(hr2.t, rayT.Max)
	if hr1.t >= hr2.t {
		return false
	}

	hr1.t = max(0, hr1.t)

	rayLength := r.Direction().Length()
	distanceInsideBoundary := (hr2.t - hr1.t) * rayLength
	hitDistance := cm.negativeInverseDensity * math.Log(rand.Float64())
	if hitDistance > distanceInsideBoundary {
		return false
	}

	record.t = hr1.t + hitDistance/rayLength
	record.p = r.At(record.t)
	record.normal = vec.New(1, 0, 0)
	record.frontFace = true
	record.Material = cm.phaseFunction
	return true
}

func (cm *constantMedium) BBox() *aabb.AABB {
	return cm.boundary.BBox()
}
