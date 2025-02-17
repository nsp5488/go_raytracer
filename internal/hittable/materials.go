package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type Material interface {
	Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool
}

type Lambertian struct {
	Albedo vec.Vec3
}

func (l *Lambertian) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	direction := record.Normal().Add(vec.RandomUnitVector())
	if direction.NearZero() {
		direction = record.Normal()
	}
	*rayOut = *ray.New(record.P(), direction)
	*attenuation = l.Albedo
	return true
}

type Metal struct {
	Albedo vec.Vec3
	Fuzz   float64
}

func (m *Metal) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	reflected := rayIn.Direction().Reflect(record.normal)
	reflected = reflected.UnitVector().Add(vec.RandomUnitVector().Scale(m.Fuzz))
	*rayOut = *ray.New(record.P(), reflected)
	*attenuation = m.Albedo
	return rayOut.Direction().Dot(record.Normal()) > 0
}

type Dielectric struct {
	RefractionIndex float64
}

func (d *Dielectric) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	*attenuation = *vec.New(1, 1, 1)

	var ri float64
	if record.FrontFace() {
		ri = 1.0 / d.RefractionIndex
	} else {
		ri = d.RefractionIndex
	}

	unitDirection := rayIn.Direction().UnitVector()
	cosineTheta := math.Min(unitDirection.Negate().Dot(record.normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosineTheta*cosineTheta)
	cannotRefract := ri*sinTheta > 1.0

	var direction *vec.Vec3
	if cannotRefract || d.reflectance(cosineTheta) > rand.Float64() {
		direction = unitDirection.Reflect(record.normal)
	} else {
		direction = unitDirection.Refract(record.normal, ri)
	}

	*rayOut = *ray.New(record.P(), direction)
	return true
}
func (d *Dielectric) reflectance(cosine float64) float64 {
	r0 := (1.0 - d.RefractionIndex) / (1.0 + d.RefractionIndex)
	r0 *= r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
