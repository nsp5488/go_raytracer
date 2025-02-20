package vec

import (
	"fmt"
	"io"
	"math"

	"github.com/nsp5488/go_raytracer/internal/interval"
)

var intensity = interval.New(0, 0.99999)

// Converts a linear color component to gamma2 space
func linearToGamma(linearComponent float64) float64 {
	if linearComponent <= 0 {
		return 0
	}
	return math.Sqrt(linearComponent)
}

// There is no API prevention on calling this for any given vec3. I may refactor this into a Color struct at some point
// Prints the color components of the vector to the given writer
func (v *Vec3) PrintColor(out io.Writer) {
	r := linearToGamma(v.e[0])
	g := linearToGamma(v.e[1])
	b := linearToGamma(v.e[2])
	rB := int(intensity.Clamp(r) * 256)
	gB := int(intensity.Clamp(g) * 256)
	bB := int(intensity.Clamp(b) * 256)

	io.WriteString(out, fmt.Sprintf("%d %d %d\n", rB, gB, bB))
}
