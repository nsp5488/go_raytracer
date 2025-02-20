package vec

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/util"
)

// A vector in 3D space
type Vec3 struct {
	e [3]float64 // elements
}

// Creates a 0 vector
func Empty() *Vec3 {
	return &Vec3{[3]float64{0, 0, 0}}
}

// Creates a new vector with the given components
func New(x, y, z float64) *Vec3 {
	return &Vec3{[3]float64{x, y, z}}
}

// Creates a random vector with components in the range [0, 1)
func Random() *Vec3 {
	return &Vec3{[3]float64{rand.Float64(), rand.Float64(), rand.Float64()}}
}

// Creates a random vector with components in the range [min, max)
func RangeRandom(min, max float64) *Vec3 {
	return &Vec3{[3]float64{util.RangeRange(min, max), util.RangeRange(min, max), util.RangeRange(min, max)}}
}

// Returns the x component of the vector
func (v *Vec3) X() float64 {
	return v.e[0]
}

// Returns the y component of the vector
func (v *Vec3) Y() float64 {
	return v.e[1]
}

// Returns the z component of the vector
func (v *Vec3) Z() float64 {
	return v.e[2]
}

// Returns the component at the given index
func (v *Vec3) Get(idx int) float64 {
	return v.e[idx]
}

// Returns the negation of the vector
func (v *Vec3) Negate() *Vec3 {
	return New(-v.X(), -v.Y(), -v.Z())
}

// Adds the given vector to the current vector in place
func (v *Vec3) AddInplace(other *Vec3) {
	v.e[0] += other.e[0]
	v.e[1] += other.e[1]
	v.e[2] += other.e[2]
}

// Scales the vector by the given factor in place
func (v *Vec3) ScaleInplace(t float64) {
	v.e[0] *= t
	v.e[1] *= t
	v.e[2] *= t
}

// Scales the vector by the given factor and returns the result
func (v *Vec3) Scale(t float64) *Vec3 {
	return New(v.e[0]*t, v.e[1]*t, v.e[2]*t)
}

// Adds the given vector to the current vector and returns the result
func (v *Vec3) Add(other *Vec3) *Vec3 {
	return New(v.e[0]+other.e[0], v.e[1]+other.e[1], v.e[2]+other.e[2])
}

// Subtracts thee given vector from the current vector and returns the result
func (v *Vec3) Sub(other *Vec3) *Vec3 {
	return New(v.e[0]-other.e[0], v.e[1]-other.e[1], v.e[2]-other.e[2])
}

// Multiplies the given vector with the current vector and returns the result
func (v *Vec3) Multiply(other *Vec3) *Vec3 {
	return New(v.e[0]*other.e[0], v.e[1]*other.e[1], v.e[2]*other.e[2])
}

// Divides the current vector by the given vector and returns the result
func (v *Vec3) Divide(other *Vec3) *Vec3 {
	return New(v.e[0]/other.e[0], v.e[1]/other.e[1], v.e[2]/other.e[2])
}

// Calculates the length squared of the vector
func (v *Vec3) LengthSquared() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

// Calculates the length of the vector
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Calculates the dot product of the current vector with the given vector
func (v *Vec3) Dot(other *Vec3) float64 {
	return v.e[0]*other.e[0] + v.e[1]*other.e[1] + v.e[2]*other.e[2]
}

// Calculates the cross product of the current vector with the given vector
func (v *Vec3) Cross(other *Vec3) *Vec3 {
	return New(
		v.e[1]*other.e[2]-v.e[2]*other.e[1],
		v.e[2]*other.e[0]-v.e[0]*other.e[2],
		v.e[0]*other.e[1]-v.e[1]*other.e[0],
	)
}

// Returns a unit vector in the direction of the current vector
func (v *Vec3) UnitVector() *Vec3 {
	return v.Scale(1 / v.Length())
}

// Checks if the vector is near zero
func (v *Vec3) NearZero() bool {
	s := 1e-8
	return math.Abs(v.e[0]) < s && math.Abs(v.e[1]) < s && math.Abs(v.e[2]) < s
}

// Reflects the vector about the given normal vector
func (v *Vec3) Reflect(normal *Vec3) *Vec3 {
	return v.Sub(normal.Scale(normal.Dot(v) * 2))
}

// Refracts the vector through the given normal vector with the given etaIOverEtaT ratio
func (v *Vec3) Refract(normal *Vec3, etaIOverEtaT float64) *Vec3 {
	cosineTheta := math.Min(v.Negate().Dot(normal), 1.0)
	rPerp := v.Add(normal.Scale(cosineTheta)).Scale(etaIOverEtaT)
	rParallel := normal.Scale(-math.Sqrt(math.Abs(1.0 - rPerp.LengthSquared())))
	return rPerp.Add(rParallel)
}

// Generates a random unit vector in the unit disk
func RandomUnitDisk() *Vec3 {
	for {
		p := New(util.RangeRange(-1, 1), util.RangeRange(-1, 1), 0)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

// Generates a random unit vector in the unit sphere
func RandomUnitVector() *Vec3 {
	for {
		p := Random()
		lenSq := p.LengthSquared()
		if 1e-160 < lenSq && lenSq <= 1 {
			return p.Scale(math.Sqrt(lenSq))
		}
	}
}

// Generates a random unit vector on the hemisphere with the given normal vector
func RandomOnHemisphere(normal *Vec3) *Vec3 {
	random := RandomUnitVector()
	if random.Dot(normal) > 0 {
		return random
	}
	return random.Negate()
}

// Returns a string representation of the vector
func (v *Vec3) String() string {
	return fmt.Sprintf("%f %f %f", v.X(), v.Y(), v.Z())
}

func (v *Vec3) Equals(other *Vec3) bool {
	return v.X() == other.X() && v.Y() == other.Y() && v.Z() == other.Z()
}
