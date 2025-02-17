package vec

import (
	"fmt"
	"io"
	"math"
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
		v.e[0]*other.e[2]-v.e[2]*other.e[0],
		v.e[1]*other.e[0]-v.e[0]*other.e[1],
	)
}

func (v *Vec3) UnitVector() *Vec3 {
	return v.Scale(1 / v.Length())
}

func (v *Vec3) String() string {
	return fmt.Sprintf("%f %f %f", v.X(), v.Y(), v.Z())
}

// There is no API prevention on calling this for any given vec3. I may refactor this into a Color struct at some point
func (v *Vec3) PrintColor(out io.Writer) {
	r := int(v.e[0] * 255.999)
	g := int(v.e[1] * 255.999)
	b := int(v.e[2] * 255.999)

	io.WriteString(out, fmt.Sprintf("%d %d %d\n", r, g, b))
}
