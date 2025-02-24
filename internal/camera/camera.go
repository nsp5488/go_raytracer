package camera

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/progress"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/util"
	"github.com/nsp5488/go_raytracer/internal/vec"

	"github.com/nsp5488/go_raytracer/internal/hittable"
)

// Defines a camera which represents the viewpoint from which the scene is rendered.
type Camera struct {
	// public members
	AspectRatio     float64
	Width           int
	Out             io.Writer
	SamplesPerPixel int
	MaxDepth        int
	MaxThreads      int
	VerticalFOV     float64
	DefocusAngle    float64
	FocusDistance   float64
	Background      *vec.Vec3

	// private members
	groupSize         chan struct{}
	waitGroup         *sync.WaitGroup
	imageHeight       int
	center            *vec.Vec3
	pixel00Loc        *vec.Vec3
	pixelDeltaU       *vec.Vec3
	pixelDeltaV       *vec.Vec3
	pixelSamplesScale float64
	defocusDiskU      *vec.Vec3
	defocusDiskV      *vec.Vec3

	u        *vec.Vec3
	v        *vec.Vec3
	w        *vec.Vec3
	lookFrom *vec.Vec3
	lookAt   *vec.Vec3
	vup      *vec.Vec3

	// progress bar state
	progressBar *tea.Program
	pbarMutex   sync.Mutex
}

// PositionCamera positions the camera with the given parameters.
func (c *Camera) PositionCamera(lookFrom, lookAt, vup *vec.Vec3) {
	if lookFrom != nil {
		c.lookFrom = lookFrom
	} else {
		c.lookFrom = vec.Empty()
	}
	if lookAt != nil {
		c.lookAt = lookAt
	} else {
		c.lookAt = vec.New(0, 0, -1)
	}
	if vup != nil {
		c.vup = vup
	} else {
		c.vup = vec.New(0, 1, 0)
	}
}

// A struct for representing a row of pixel data.
type rowData struct {
	index int
	data  *bytes.Buffer
}

// calculates the pixel data for one row of the image utilizing a thread pool.
func (c *Camera) renderRow(world *hittable.HittableList, buf *rowData) {
	defer c.waitGroup.Done()
	c.groupSize <- struct{}{}
	for j := range c.Width {
		pixelColor := vec.Empty()
		for range c.SamplesPerPixel {
			r := c.getRay(j, buf.index)
			pixelColor.AddInplace(c.rayColor(r, world, c.MaxDepth))

		}

		pixelColor.Scale(c.pixelSamplesScale).PrintColor(buf.data)
	}
	<-c.groupSize
	c.pbarMutex.Lock()
	c.progressBar.Send(1)
	c.pbarMutex.Unlock()
}

// A threaded variant of the renderer.
func (c *Camera) threadedRenderer(world *hittable.HittableList) {
	buffers := make([]rowData, c.imageHeight, c.imageHeight)
	for i := range c.imageHeight {
		buffers[i].data = &bytes.Buffer{}
		buffers[i].index = i
	}

	for i := range c.imageHeight {
		c.waitGroup.Add(1)
		go c.renderRow(world, &buffers[i])
	}
	c.waitGroup.Wait()

	sort.Slice(buffers, func(i, j int) bool {
		return buffers[i].index < buffers[j].index
	})
	for _, buf := range buffers {
		buf.data.WriteTo(c.Out)
	}
	c.progressBar.Send(1)
}

// A synchronous variant of the renderer.
func (c *Camera) syncRenderer(world *hittable.HittableList) {
	for i := range c.imageHeight {
		for j := range c.Width {
			pixelColor := vec.Empty()
			for range c.SamplesPerPixel {
				r := c.getRay(j, i)
				pixelColor.AddInplace(c.rayColor(r, world, c.MaxDepth))

			}

			pixelColor.Scale(c.pixelSamplesScale).PrintColor(c.Out)
		}
		c.progressBar.Send(1)
	}
	c.progressBar.Send(1)
}

// Render the provided scene using the camera's settings.
func (c *Camera) Render(world *hittable.HittableList) {
	c.initialize()

	io.WriteString(c.Out, fmt.Sprintf("P3\n%d %d\n255\n", c.Width, c.imageHeight))

	// Run the processing in a separate goroutine
	if c.MaxThreads <= 1 {
		// use a low-overhead synchronous renderer if we're only alloted one thread.
		go c.syncRenderer(world)
	} else {
		go c.threadedRenderer(world)
	}

	if _, err := c.progressBar.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		c.progressBar.ReleaseTerminal()
		return
	}

}

// initialize the camera's settings.
func (c *Camera) initialize() {
	// Ensure defaults:
	if c.AspectRatio == 0 {
		c.AspectRatio = 1.0
	}
	if c.Width == 0 {
		c.Width = 100
	}
	if c.Out == nil {
		log.Fatal("Must specify an output")
	}
	if c.SamplesPerPixel == 0 {
		c.SamplesPerPixel = 100
	}
	if c.MaxDepth == 0 {
		c.MaxDepth = 10
	}
	if c.VerticalFOV == 0 {
		c.VerticalFOV = 90
	}
	if c.FocusDistance == 0 {
		c.FocusDistance = 10
	}
	if c.DefocusAngle == 0 {
		c.DefocusAngle = 0
	}

	// calculate image height given aspect ratio, clamped to >=1
	c.imageHeight = max(1, int(float64(c.Width)/c.AspectRatio))
	c.pixelSamplesScale = 1.0 / float64(c.SamplesPerPixel)

	// define camera information
	c.center = c.lookFrom

	theta := util.DegressToRadians(c.VerticalFOV)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h * c.FocusDistance
	viewportWidth := viewportHeight * (float64(c.Width) / float64(c.imageHeight))

	// calculate camera basis vectors
	c.w = c.lookFrom.Sub(c.lookAt).UnitVector()
	c.u = c.vup.Cross(c.w).UnitVector()
	c.v = c.w.Cross(c.u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := c.u.Scale(viewportWidth)
	viewportV := c.v.Negate().Scale(viewportHeight)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	c.pixelDeltaU = viewportU.Scale(1.0 / float64(c.Width))
	c.pixelDeltaV = viewportV.Scale(1.0 / float64(c.imageHeight))

	// Calculate the location of the upper left pixel.
	viewportTopLeft := c.center.Sub(
		c.w.Scale(c.FocusDistance)).
		Sub(viewportU.Scale(0.5)).
		Sub(viewportV.Scale(0.5))
	c.pixel00Loc = viewportTopLeft.Add(c.pixelDeltaU.Add(c.pixelDeltaV).Scale(0.5))

	// calculate defocus disk basis vectors
	defocusRadius := c.FocusDistance * math.Tan(util.DegressToRadians(c.DefocusAngle/2.0))
	c.defocusDiskU = c.u.Scale(defocusRadius)
	c.defocusDiskV = c.v.Scale(defocusRadius)

	c.waitGroup = &sync.WaitGroup{}
	c.groupSize = make(chan struct{}, c.MaxThreads)

	// initialize the progress bar
	c.progressBar = progress.InitBar(c.imageHeight + 1)
}

// getRay returns a ray from the camera with some amount of defocus and sampling to offset. This creates a smoother image and simulates depth of field.
func (c *Camera) getRay(i, j int) *ray.Ray {
	offset := c.sampleSquare()
	pixelSample := c.pixel00Loc.
		Add(c.pixelDeltaU.Scale(float64(i) + offset.X())).
		Add(c.pixelDeltaV.Scale(float64(j) + offset.Y()))
	var rayOrigin *vec.Vec3
	if c.DefocusAngle <= 0 {
		rayOrigin = c.center
	} else {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.Sub(rayOrigin)
	rayTime := rand.Float64()
	return ray.NewWithTime(rayOrigin, rayDirection, rayTime)
}

// Returns a random offset within a 1x1 square
func (c *Camera) sampleSquare() *vec.Vec3 {
	return vec.New(rand.Float64()-0.5, rand.Float64()-0.5, 0)
}

// Returns a randomly offset center point for the ray to simulate depth of field
func (c *Camera) defocusDiskSample() *vec.Vec3 {
	p := vec.RandomUnitDisk()
	return c.center.
		Add(c.defocusDiskU.Scale(p.X())).
		Add(c.defocusDiskV.Scale(p.Y()))
}

// Calculates the color of a ray after it has been traced through the scene.
func (c *Camera) rayColor(r *ray.Ray, world *hittable.HittableList, depth int) *vec.Vec3 {
	if depth < 0 {
		return vec.Empty()
	}
	rec := hittable.HitRecord{}
	if !world.Hit(r, *interval.New(0.001, math.Inf(1)), &rec) {
		return c.Background
	}

	emitColor := vec.Empty()
	if emit, ok := rec.Material.(hittable.EmissiveMaterial); ok {
		emitColor = emit.Emitted(rec.U(), rec.V(), rec.P())
	}

	scattered := &ray.Ray{}
	attenuation := &vec.Vec3{}
	scatterColor := vec.Empty()
	if !rec.Material.Scatter(r, scattered, &rec, attenuation) {
		return emitColor
	}
	scatterColor = attenuation.Multiply(c.rayColor(scattered, world, depth-1))

	return emitColor.Add(scatterColor)
}
