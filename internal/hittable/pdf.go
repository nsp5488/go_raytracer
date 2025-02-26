package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/vec"
)

type Pdf interface {
	Value(direction *vec.Vec3) float64
	Generate() *vec.Vec3
}

type SpherePdf struct{}

func (s *SpherePdf) Value(direction *vec.Vec3) float64 {
	return 1 / (4 * math.Pi)
}

func (s *SpherePdf) Generate() *vec.Vec3 {
	return vec.RandomUnitVector()
}

type cosinePdf struct {
	onb *orthonormalBasis
}

func CosinePdf(normal *vec.Vec3) *cosinePdf {
	return &cosinePdf{onb: NewONB(normal)}
}

func (c *cosinePdf) Value(direction *vec.Vec3) float64 {
	cosTheta := direction.UnitVector().Dot(c.onb.W())
	return math.Max(0, cosTheta/math.Pi)
}

func (c *cosinePdf) Generate() *vec.Vec3 {
	return c.onb.Transform(vec.RandomCosineDirection())
}

type hittablePdf struct {
	object Hittable
	origin *vec.Vec3
}

func HittablePdf(origin *vec.Vec3, object Hittable) *hittablePdf {
	return &hittablePdf{object: object, origin: origin}
}

func (hp *hittablePdf) Value(direction *vec.Vec3) float64 {
	return hp.object.PdfValue(hp.origin, direction)
}
func (hp *hittablePdf) Generate() *vec.Vec3 {
	return hp.object.Random(hp.origin)
}

type mixturePdf struct {
	p [2]Pdf
}

func MixturePdf(p0, p1 Pdf) *mixturePdf {
	return &mixturePdf{p: [2]Pdf{p0, p1}}
}
func (mp *mixturePdf) Value(direction *vec.Vec3) float64 {
	return 0.5*mp.p[0].Value(direction) + 0.5*mp.p[1].Value(direction)
}

func (mp *mixturePdf) Generate() *vec.Vec3 {
	if rand.Float64() < 0.5 {
		return mp.p[0].Generate()
	}
	return mp.p[1].Generate()
}
