package vec

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/util"
)

type Vec3 struct {
	e [3]float64 // elements
}

func Empty() *Vec3 {
	return &Vec3{[3]float64{0, 0, 0}}
}
func New(x, y, z float64) *Vec3 {
	return &Vec3{[3]float64{x, y, z}}
}
func Random() *Vec3 {
	return &Vec3{[3]float64{rand.Float64(), rand.Float64(), rand.Float64()}}
}
func RangeRandom(min, max float64) *Vec3 {
	return &Vec3{[3]float64{util.RangeRange(min, max), util.RangeRange(min, max), util.RangeRange(min, max)}}
}
func (v *Vec3) X() float64 {
	return v.e[0]
}
func (v *Vec3) Y() float64 {
	return v.e[1]
}
func (v *Vec3) Z() float64 {
	return v.e[2]
}

func (v *Vec3) Get(idx int) float64 {
	return v.e[idx]
}
func (v *Vec3) Negate() *Vec3 {
	return New(-v.X(), -v.Y(), -v.Z())
}

func (v *Vec3) AddInplace(other *Vec3) {
	v.e[0] += other.e[0]
	v.e[1] += other.e[1]
	v.e[2] += other.e[2]
}
func (v *Vec3) ScaleInplace(t float64) {
	v.e[0] *= t
	v.e[1] *= t
	v.e[2] *= t
}

func (v *Vec3) Scale(t float64) *Vec3 {
	return New(v.e[0]*t, v.e[1]*t, v.e[2]*t)
}

func (v *Vec3) Add(other *Vec3) *Vec3 {
	return New(v.e[0]+other.e[0], v.e[1]+other.e[1], v.e[2]+other.e[2])
}
func (v *Vec3) Multiply(other *Vec3) *Vec3 {
	return New(v.e[0]*other.e[0], v.e[1]*other.e[1], v.e[2]*other.e[2])
}
func (v *Vec3) Divide(other *Vec3) *Vec3 {
	return New(v.e[0]/other.e[0], v.e[1]/other.e[1], v.e[2]/other.e[2])
}

func (v *Vec3) LengthSquared() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vec3) Dot(other *Vec3) float64 {
	return v.e[0]*other.e[0] + v.e[1]*other.e[1] + v.e[2]*other.e[2]
}
func (v *Vec3) Cross(other *Vec3) *Vec3 {
	return New(
		v.e[1]*other.e[2]-v.e[2]*other.e[1],
		v.e[2]*other.e[0]-v.e[0]*other.e[2],
		v.e[0]*other.e[1]-v.e[1]*other.e[0],
	)
}

func (v *Vec3) UnitVector() *Vec3 {
	return v.Scale(1 / v.Length())
}
func (v *Vec3) NearZero() bool {
	s := 1e-8
	return math.Abs(v.e[0]) < s && math.Abs(v.e[1]) < s && math.Abs(v.e[2]) < s
}

func (v *Vec3) Reflect(normal *Vec3) *Vec3 {
	return v.Add(normal.Scale(normal.Dot(v) * 2).Negate())
}
func (v *Vec3) Refract(normal *Vec3, etaIOverEtaT float64) *Vec3 {
	cosineTheta := math.Min(v.Negate().Dot(normal), 1.0)
	rPerp := v.Add(normal.Scale(cosineTheta)).Scale(etaIOverEtaT)
	rParallel := normal.Scale(-math.Sqrt(math.Abs(1.0 - rPerp.LengthSquared())))
	return rPerp.Add(rParallel)
}

func RandomUnitDisk() *Vec3 {
	for {
		p := New(util.RangeRange(-1, 1), util.RangeRange(-1, 1), 0)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func RandomUnitVector() *Vec3 {
	for {
		p := Random()
		lenSq := p.LengthSquared()
		if 1e-160 < lenSq && lenSq <= 1 {
			return p.Scale(math.Sqrt(lenSq))
		}
	}
}

func RandomOnHemisphere(normal *Vec3) *Vec3 {
	random := RandomUnitVector()
	if random.Dot(normal) > 0 {
		return random
	}
	return random.Negate()
}

func (v *Vec3) String() string {
	return fmt.Sprintf("%f %f %f", v.X(), v.Y(), v.Z())
}
