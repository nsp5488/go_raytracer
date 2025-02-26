package hittable

import (
	"math"

	"github.com/nsp5488/go_raytracer/internal/vec"
)

type orthonormalBasis struct {
	axis [3]*vec.Vec3
}

func NewONB(n *vec.Vec3) *orthonormalBasis {
	onb := &orthonormalBasis{}
	onb.axis[2] = n.UnitVector()
	var a *vec.Vec3
	if math.Abs(n.X()) > .9 {
		a = vec.New(0, 1, 0)
	} else {
		a = vec.New(1, 0, 0)
	}
	onb.axis[1] = n.Cross(a).UnitVector()
	onb.axis[0] = n.Cross(onb.axis[1]).UnitVector()
	return onb
}

func (onb *orthonormalBasis) U() *vec.Vec3 {
	return onb.axis[0]
}
func (onb *orthonormalBasis) V() *vec.Vec3 {
	return onb.axis[1]
}

func (onb *orthonormalBasis) W() *vec.Vec3 {
	return onb.axis[2]
}

func (onb *orthonormalBasis) Transform(v *vec.Vec3) *vec.Vec3 {
	// ux + vy + zw = O'
	return onb.axis[0].Scale(v.X()).
		Add(onb.axis[1].Scale(v.Y())).
		Add(onb.axis[2].Scale(v.Z()))
}
