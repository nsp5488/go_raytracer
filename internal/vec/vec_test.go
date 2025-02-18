package vec_test

import (
	"math"
	"testing"

	"github.com/nsp5488/go_raytracer/internal/vec"
)

func checkVec(t *testing.T, v *vec.Vec3, expected *vec.Vec3) {
	if !v.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, v)
	}
}

func checkVecFloat(t *testing.T, v float64, expected float64) {
	if v != expected {
		t.Errorf("Expected %v, got %v", expected, v)
	}
}

func TestVecAdd(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(4, 5, 6)
	expected := vec.New(5, 7, 9)
	checkVec(t, v1.Add(v2), expected)
}

func TestVecSub(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(4, 5, 6)
	expected := vec.New(-3, -3, -3)
	checkVec(t, v1.Add(v2.Negate()), expected)
}

func TestVecMul(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(4, 5, 6)
	expected := vec.New(4, 10, 18)
	checkVec(t, v1.Multiply(v2), expected)
}

func TestVecDiv(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(4, 5, 6)
	expected := vec.New(0.25, 0.4, 0.5)
	checkVec(t, v1.Divide(v2), expected)
}

func TestVecNegate(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := vec.New(-1, -2, -3)
	checkVec(t, v.Negate(), expected)
}

func TestVecEquals(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(1, 2, 3)
	v3 := vec.New(1, 2, 4)
	checkVec(t, v1, v2)
	if v1.Equals(v3) {
		t.Error("Expected false, got true")
	}
}

func TestVecDot(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(4, 5, 6)
	expected := 32.0
	checkVecFloat(t, v1.Dot(v2), expected)
}

func TestVecCross(t *testing.T) {
	v1 := vec.New(1, 2, 3)
	v2 := vec.New(4, 5, 6)
	expected := vec.New(-3, 6, -3)
	checkVec(t, v1.Cross(v2), expected)
}

func TestVecScale(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := vec.New(2, 4, 6)
	checkVec(t, v.Scale(2), expected)
}

func TestVecScaleInplace(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := vec.New(2, 4, 6)
	v.ScaleInplace(2)
	checkVec(t, v, expected)
}

func TestVecAddInplace(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := vec.New(4, 5, 6)
	v.AddInplace(vec.New(3, 3, 3))
	checkVec(t, v, expected)
}

func TestVecLengthSquared(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := 14.0
	checkVecFloat(t, v.LengthSquared(), expected)
}

func TestVecLength(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := math.Sqrt(14)
	checkVecFloat(t, v.Length(), expected)
}

func TestVecUnitVector(t *testing.T) {
	v := vec.New(1, 2, 3)
	expected := vec.New(1/math.Sqrt(14), 2/math.Sqrt(14), 3/math.Sqrt(14))
	checkVec(t, v.UnitVector(), expected)
}

func TestVecNearZero(t *testing.T) {
	v := vec.New(1e-9, 1e-9, 1e-9)
	if !v.NearZero() {
		t.Error("Expected true, got false")
	}
	v = vec.New(1e-7, 1e-7, 1e-7)
	if v.NearZero() {
		t.Error("Expected false, got true")
	}
}
