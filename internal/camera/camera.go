package camera

import (
	"fmt"
	"io"
	"math"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/progress"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"

	"github.com/nsp5488/go_raytracer/internal/hittable"
)

type Camera struct {
	AspectRatio     float64
	Width           int
	Out             io.Writer
	SamplesPerPixel int
	MaxDepth        int

	imageHeight       int
	center            *vec.Vec3
	pixel00Loc        *vec.Vec3
	pixelDeltaU       *vec.Vec3
	pixelDeltaV       *vec.Vec3
	pixelSamplesScale float64

	progressBar *tea.Program
}

func (c *Camera) Render(world *hittable.HittableList) {
	c.initialize()

	io.WriteString(c.Out, fmt.Sprintf("P3\n%d %d\n255\n", c.Width, c.imageHeight))
	// Run the processing in a separate goroutine

	go func() {
		for i := range c.imageHeight {
			for j := range c.Width {

				pixelColor := vec.Empty()
				for range c.SamplesPerPixel {
					r := c.getRay(j, i)
					pixelColor.AddInplace(c.rayColor(r, world, c.MaxDepth))

				}

				pixelColor.Scale(c.pixelSamplesScale).PrintColor(c.Out)
			}
			c.progressBar.Send(i + 1)
		}
	}()

	if _, err := c.progressBar.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		return
	}

}
func (c *Camera) initialize() {

	// calculate image height given aspect ratio, clamped to >=1
	c.imageHeight = max(1, int(float64(c.Width)/c.AspectRatio))
	c.pixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)
	// define camera information
	focalLength := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * (float64(c.Width) / float64(c.imageHeight))
	c.center = vec.Empty()

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := vec.New(viewportWidth, 0, 0)
	viewportV := vec.New(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	c.pixelDeltaU = viewportU.Scale(1.0 / float64(c.Width))
	c.pixelDeltaV = viewportV.Scale(1.0 / float64(c.imageHeight))

	// Calculate the location of the upper left pixel.
	viewportTopLeft := c.center.Add(vec.New(0, 0, focalLength).Negate()).Add(viewportU.Scale(0.5).Negate()).Add(viewportV.Scale(0.5).Negate())
	c.pixel00Loc = viewportTopLeft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Scale(0.5))

	c.progressBar = progress.InitBar(c.imageHeight)
}

func (c *Camera) getRay(i, j int) *ray.Ray {
	offset := sampleSquare()
	pixelSample := c.pixel00Loc.
		Add(c.pixelDeltaU.Scale(float64(i) + offset.X())).
		Add(c.pixelDeltaV.Scale(float64(j) + offset.Y()))
	ray_direction := pixelSample.Add(c.center.Negate())
	return ray.New(c.center, ray_direction)
}

func sampleSquare() *vec.Vec3 {
	return vec.New(rand.Float64()-0.5, rand.Float64()-0.5, 0)
}
func (c *Camera) rayColor(r *ray.Ray, world *hittable.HittableList, depth int) *vec.Vec3 {
	if depth < 0 {
		return vec.Empty()
	}
	rec := hittable.HitRecord{}

	if world.Hit(r, *interval.New(0.001, math.Inf(1)), &rec) {
		scattered := &ray.Ray{}
		attenuation := &vec.Vec3{}
		if rec.Material.Scatter(r, scattered, &rec, attenuation) {
			return attenuation.Multiply(c.rayColor(scattered, world, depth-1))
		}
		return vec.Empty()
	}
	unitDirection := r.Direction().UnitVector()
	a := 0.5 * (unitDirection.Y() + 1.0)
	c1 := vec.New(1, 1, 1).Scale(1.0 - a)
	c2 := vec.New(0.5, 0.7, 1).Scale(a)
	return c1.Add(c2)
}
