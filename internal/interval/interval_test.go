package interval_test

import (
	"testing"

	"github.com/nsp5488/go_raytracer/internal/interval"
)

func TestIntervalContains(t *testing.T) {
	i := interval.New(0, 1)
	if !i.Contains(0.5) {
		t.Error("Expected true, got false")
	}
}

func TestIntervalContainsEdge(t *testing.T) {
	i := interval.New(0, 1)
	if !i.Contains(0) || !i.Contains(1) {
		t.Error("Expected true, got false")
	}
}

func TestIntervalContainsOutside(t *testing.T) {
	i := interval.New(0, 1)
	if i.Contains(-1) || i.Contains(2) {
		t.Error("Expected false, got true")
	}
}

func TestIntervalSurrounds(t *testing.T) {
	i := interval.New(0, 1)
	if !i.Surrounds(0.5) {
		t.Error("Expected true, got false")
	}
}

func TestIntervalSurroundsEdge(t *testing.T) {
	i := interval.New(0, 1)
	if !i.Surrounds(0.001) || !i.Surrounds(0.99) {
		t.Error("Expected true, got false")
	}
}

func TestIntervalSurroundsOutside(t *testing.T) {
	i := interval.New(0, 1)
	if i.Surrounds(-1) || i.Surrounds(1.1) {
		t.Error("Expected false, got true")
	}
}

func TestIntervalClamp(t *testing.T) {
	i := interval.New(0, 1)
	if i.Clamp(0.5) != 0.5 {
		t.Error("Expected interval.New(0.5, 0.5), got", i.Clamp(0.5))
	}
	if i.Clamp(-1) != 0 {
		t.Error("Expected interval.New(0, 1), got", i.Clamp(-1))
	}
	if i.Clamp(2) != 1 {
		t.Error("Expected interval.New(0, 1), got", i.Clamp(2))
	}
}

func TestIntervalClampEdge(t *testing.T) {
	i := interval.New(0, 1)
	if i.Clamp(0) != 0 {
		t.Error("Expected interval.New(0, 0), got", i.Clamp(0))
	}
	if i.Clamp(1) != 1 {
		t.Error("Expected interval.New(1, 1), got", i.Clamp(1))
	}
}
