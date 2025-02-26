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
	ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64
}

type EmissiveMaterial interface {
	Material
	Emitted(u, v float64, point *vec.Vec3) *vec.Vec3
}

// Lambertian (matte) material.
type lambertian struct {
	tex Texture
}

// Creates a new matte material
func NewLambertian(albedo *vec.Vec3) *lambertian {
	return &lambertian{tex: NewSolidColor(albedo)}
}

// A lambertian with an externally defined texture
func NewTexturedLambertian(tex Texture) *lambertian {
	return &lambertian{tex: tex}
}

// Scatter implements the Lambertian material's scattering behavior.
func (l *lambertian) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	direction := vec.RandomOnHemisphere(record.normal)
	if direction.NearZero() {
		direction = record.Normal()
	}
	*rayOut = *ray.NewWithTime(record.P(), direction, rayIn.Time())
	*attenuation = *l.tex.Value(record.u, record.v, record.p)
	return true
}
func (l *lambertian) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 1 / (2 * math.Pi)
	cosTheta := record.Normal().Dot(rayOut.Direction().UnitVector())
	if cosTheta < 0 {
		return 0
	}
	return cosTheta / math.Pi
}

// Metal material.
type metal struct {
	Albedo *vec.Vec3
	Fuzz   float64
}

func NewMetal(albedo *vec.Vec3, fuzz float64) *metal {
	return &metal{Albedo: albedo, Fuzz: fuzz}
}

// Scatter implements the metal material's scattering behavior.
func (m *metal) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	reflected := rayIn.Direction().Reflect(record.normal)
	reflected = reflected.UnitVector().Add(vec.RandomUnitVector().Scale(m.Fuzz))
	*rayOut = *ray.NewWithTime(record.P(), reflected, rayIn.Time())
	*attenuation = *m.Albedo
	return rayOut.Direction().Dot(record.Normal()) > 0
}
func (m *metal) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

// Dielectric material.
type dielectric struct {
	RefractionIndex float64
}

func NewDielectric(refractionIndex float64) *dielectric {
	return &dielectric{RefractionIndex: refractionIndex}
}

// Scatter implements the dielectric material's scattering behavior.
func (d dielectric) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
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
func (d dielectric) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

// Helper function to calculate the reflectance of a dielectric material
func (d dielectric) reflectance(cosine float64) float64 {
	r0 := (1.0 - d.RefractionIndex) / (1.0 + d.RefractionIndex)
	r0 *= r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}

type diffuseLight struct {
	tex Texture
}

func NewDiffuseLight(color *vec.Vec3) *diffuseLight {
	return &diffuseLight{tex: NewSolidColor(color)}
}
func newDiffuseLightTextured(tex Texture) *diffuseLight {
	return &diffuseLight{tex: tex}
}
func (dl diffuseLight) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

func (dl diffuseLight) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	*attenuation = *vec.New(1, 1, 1)
	return false
}

func (dl diffuseLight) Emitted(u, v float64, point *vec.Vec3) *vec.Vec3 {
	return dl.tex.Value(u, v, point)
}

type isotropic struct {
	tex Texture
}

func (i isotropic) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

func NewIsotropicTexture(tex Texture) *isotropic {
	return &isotropic{tex: tex}
}
func NewIsotropic(albedo *vec.Vec3) *isotropic {
	return &isotropic{tex: NewSolidColor(albedo)}
}

func (i *isotropic) Scatter(rayIn, rayOut *ray.Ray, record *HitRecord, attenuation *vec.Vec3) bool {
	*rayOut = *ray.NewWithTime(record.p, vec.RandomUnitVector(), rayIn.Time())
	*attenuation = *i.tex.Value(record.u, record.v, record.p)
	return true
}
