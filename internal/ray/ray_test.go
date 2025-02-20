package ray_test

import (
	"fmt"
	"testing"

	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

func TestRayAt(t *testing.T) {
	ray := ray.New(vec.New(0, 0, 0), vec.New(0, 1, 0))
	exp := vec.New(0, 1, 0)
	act := ray.At(1)
	if !exp.Equals(act) {
		t.Errorf("Expected %s, but got %s", exp.String(), act.String())
	}
}
func TestString(t *testing.T) {
	ray := ray.New(vec.New(1, 2, 3), vec.New(4, 5, 6))
	exp := fmt.Sprintf("(%f, %f, %f) + (%f, %f, %f)t", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0)
	act := ray.String()
	if exp != act {
		t.Errorf("Expected %s, but got %s", exp, act)
	}
}
