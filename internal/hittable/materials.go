package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

type ScatterRecord struct {
	Attenuation *vec.Vec3
	Pdf         Pdf
	SkipPdf     bool
	SkipPdfRay  *ray.Ray
}

// Material interface defines the behavior of a material when a ray hits it.
type Material interface {
	Scatter(rayIn *ray.Ray, record *HitRecord, srecord *ScatterRecord) bool
	ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64
}

type EmissiveMaterial interface {
	Material
	Emitted(record *HitRecord) *vec.Vec3
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
func (l *lambertian) Scatter(rayIn *ray.Ray, record *HitRecord, srecord *ScatterRecord) bool {
	srecord.Attenuation = l.tex.Value(record.u, record.v, record.p)
	srecord.Pdf = CosinePdf(record.normal)
	srecord.SkipPdf = false
	return true
}
func (l *lambertian) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
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
func (m *metal) Scatter(rayIn *ray.Ray, record *HitRecord, srecord *ScatterRecord) bool {
	reflected := rayIn.Direction().Reflect(record.normal)
	reflected = reflected.UnitVector().Add(vec.RandomUnitVector().Scale(m.Fuzz))

	srecord.Attenuation = m.Albedo
	srecord.Pdf = nil
	srecord.SkipPdf = true
	srecord.SkipPdfRay = ray.NewWithTime(record.p, reflected, rayIn.Time())
	return true
}
func (m *metal) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

// Dielectric material.
type Dielectric struct {
	RefractionIndex float64
}

func NewDielectric(refractionIndex float64) *Dielectric {
	return &Dielectric{RefractionIndex: refractionIndex}
}

// Scatter implements the dielectric material's scattering behavior.
func (d Dielectric) Scatter(rayIn *ray.Ray, record *HitRecord, srecord *ScatterRecord) bool {
	srecord.Attenuation = vec.New(1, 1, 1)
	srecord.Pdf = nil
	srecord.SkipPdf = true

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

	srecord.SkipPdfRay = ray.NewWithTime(record.P(), direction, rayIn.Time())
	return true
}
func (d Dielectric) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

// Helper function to calculate the reflectance of a dielectric material
func (d Dielectric) reflectance(cosine float64) float64 {
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
func NewDiffuseLightTextured(tex Texture) *diffuseLight {
	return &diffuseLight{tex: tex}
}
func (dl diffuseLight) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 0
}

func (dl diffuseLight) Scatter(rayIn *ray.Ray, record *HitRecord, srecord *ScatterRecord) bool {
	return false
}

func (dl diffuseLight) Emitted(record *HitRecord) *vec.Vec3 {
	if !record.frontFace {
		return vec.Empty()
	}
	return dl.tex.Value(record.u, record.v, record.p)
}

type isotropic struct {
	tex Texture
}

func (i isotropic) ScatteringPdf(rayIn, rayOut *ray.Ray, record *HitRecord) float64 {
	return 1 / (4 * math.Pi)
}

func NewIsotropicTexture(tex Texture) *isotropic {
	return &isotropic{tex: tex}
}
func NewIsotropic(albedo *vec.Vec3) *isotropic {
	return &isotropic{tex: NewSolidColor(albedo)}
}

func (i *isotropic) Scatter(rayIn *ray.Ray, record *HitRecord, srecord *ScatterRecord) bool {
	srecord.Attenuation = i.tex.Value(record.u, record.v, record.p)
	srecord.Pdf = &SpherePdf{}
	srecord.SkipPdf = false
	return true
}
