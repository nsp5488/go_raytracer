package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Material interface defines the behavior of a material when a ray hits it.
type Material interface {
	Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool
}

// Lambertian (matte) material.
type Lambertian struct {
	tex Texture
}

// Creates a new matte material
func NewLambertian(albedo *vec.Vec3) *Lambertian {
	return &Lambertian{tex: NewSolidColor(albedo)}
}

// A lambertian with an externally defined texture
func NewTexturedLambertian(tex Texture) *Lambertian {
	return &Lambertian{tex: tex}
}

// Scatter implements the Lambertian material's scattering behavior.
func (l Lambertian) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	direction := record.Normal().Add(vec.RandomUnitVector())
	if direction.NearZero() {
		direction = record.Normal()
	}
	*rayOut = *ray.NewWithTime(record.P(), direction, rayIn.Time())
	*attenuation = *l.tex.Value(record.u, record.v, record.p)
	return true
}

// Metal material.
type Metal struct {
	Albedo vec.Vec3
	Fuzz   float64
}

// Scatter implements the metal material's scattering behavior.
func (m Metal) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	reflected := rayIn.Direction().Reflect(record.normal)
	reflected = reflected.UnitVector().Add(vec.RandomUnitVector().Scale(m.Fuzz))
	*rayOut = *ray.NewWithTime(record.P(), reflected, rayIn.Time())
	*attenuation = m.Albedo
	return rayOut.Direction().Dot(record.Normal()) > 0
}

// Dielectric material.
type Dielectric struct {
	RefractionIndex float64
}

// Scatter implements the dielectric material's scattering behavior.
func (d Dielectric) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
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

	*rayOut = *ray.NewWithTime(record.P(), direction, rayIn.Time())
	return true
}

// Helper function to calculate the reflectance of a dielectric material
func (d Dielectric) reflectance(cosine float64) float64 {
	r0 := (1.0 - d.RefractionIndex) / (1.0 + d.RefractionIndex)
	r0 *= r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
