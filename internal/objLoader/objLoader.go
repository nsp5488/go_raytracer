package objLoader

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/nsp5488/go_raytracer/internal/hittable"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// LoadObjOptions provides configuration options for loading OBJ files
type LoadObjOptions struct {
	ScaleFactor   float64
	FlipYZ        bool
	Debug         bool
	IgnoreNormals bool
	Center        bool
	FlipFaces     bool
	Position      *vec.Vec3
}

// DefaultLoadOptions provides reasonable defaults
func DefaultLoadOptions() LoadObjOptions {
	return LoadObjOptions{
		ScaleFactor:   1.0,
		FlipYZ:        false,
		Debug:         false,
		IgnoreNormals: false,
		Center:        true,
		FlipFaces:     false,
		Position:      vec.New(0, 0, 0),
	}
}

func fixIndex(i, length int) int {
	if i < 0 {
		i = length + i // Convert negative index
	} else {
		i = i - 1 // Convert 1-based to 0-based
	}

	// Safety check to avoid out-of-bounds access
	if i < 0 || i >= length {
		log.Printf("Warning: Index %d out of bounds (0-%d), clamping", i, length-1)
		i = int(math.Max(0, math.Min(float64(i), float64(length-1))))
	}

	return i
}

// LoadObj loads a 3D model from an OBJ file with default options
func LoadObj(filename string, mat hittable.Material) hittable.Hittable {
	return LoadObjWithOptions(filename, DefaultLoadOptions(), mat)
}

// LoadObjWithOptions loads a 3D model from an OBJ file with custom options
func LoadObjWithOptions(filename string, options LoadObjOptions, mat hittable.Material) hittable.Hittable {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open file %s: %v", filename, err)
	}
	defer file.Close()

	// Store raw vertices separately for manipulation
	var rawVertices []*vec.Vec3
	var vertices []*vec.Vec3 // These will be the processed vertices
	var normals []*vec.Vec3
	var triangles []*hittable.Triangle

	// For computing bounds
	minBounds := [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	maxBounds := [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// First pass: read vertices and compute bounds
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		if parts[0] == "v" { // Vertex
			if len(parts) < 4 {
				log.Printf("Line %d: Malformed vertex, expected at least 3 coordinates: %s", lineNum, line)
				continue
			}

			x, errX := strconv.ParseFloat(parts[1], 64)
			y, errY := strconv.ParseFloat(parts[2], 64)
			z, errZ := strconv.ParseFloat(parts[3], 64)

			if errX != nil || errY != nil || errZ != nil {
				log.Printf("Line %d: Invalid vertex coordinates: %s", lineNum, line)
				continue
			}

			// Apply scale but store this raw value
			x *= options.ScaleFactor
			y *= options.ScaleFactor
			z *= options.ScaleFactor

			if options.FlipYZ {
				y, z = z, y
			}

			vertex := vec.New(x, y, z)
			rawVertices = append(rawVertices, vertex)

			// Update bounds for centering calculation
			minBounds[0] = math.Min(minBounds[0], x)
			minBounds[1] = math.Min(minBounds[1], y)
			minBounds[2] = math.Min(minBounds[2], z)
			maxBounds[0] = math.Max(maxBounds[0], x)
			maxBounds[1] = math.Max(maxBounds[1], y)
			maxBounds[2] = math.Max(maxBounds[2], z)
		}
	}

	// Calculate center before further processing
	center := vec.New(
		(minBounds[0]+maxBounds[0])/2,
		(minBounds[1]+maxBounds[1])/2,
		(minBounds[2]+maxBounds[2])/2,
	)

	// Reset file for second pass
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	lineNum = 0

	// Debug info about model bounds
	if options.Debug {
		fmt.Printf("=== OBJ MODEL DIMENSIONS ===\n")
		fmt.Printf("Min bounds: [%f, %f, %f]\n", minBounds[0], minBounds[1], minBounds[2])
		fmt.Printf("Max bounds: [%f, %f, %f]\n", maxBounds[0], maxBounds[1], maxBounds[2])
		fmt.Printf("Center: [%f, %f, %f]\n", center.X(), center.Y(), center.Z())

		width := maxBounds[0] - minBounds[0]
		height := maxBounds[1] - minBounds[1]
		depth := maxBounds[2] - minBounds[2]
		fmt.Printf("Dimensions: width=%f, height=%f, depth=%f\n", width, height, depth)

		diag := math.Sqrt(width*width + height*height + depth*depth)
		fmt.Printf("Diagonal length: %f\n", diag)
	}

	// Process vertices with centering if requested
	for _, v := range rawVertices {
		transformedVertex := vec.New(v.X(), v.Y(), v.Z())

		// Apply centering if requested
		if options.Center {
			transformedVertex.AddInplace(center.Negate())

			// Apply desired position offset after centering
			transformedVertex.AddInplace(options.Position)
		}

		vertices = append(vertices, transformedVertex)
	}

	// Verify centering worked
	if options.Debug && options.Center {
		// Calculate new bounds
		newMinBounds := [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
		newMaxBounds := [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

		for _, v := range vertices {
			newMinBounds[0] = math.Min(newMinBounds[0], v.X())
			newMinBounds[1] = math.Min(newMinBounds[1], v.Y())
			newMinBounds[2] = math.Min(newMinBounds[2], v.Z())
			newMaxBounds[0] = math.Max(newMaxBounds[0], v.X())
			newMaxBounds[1] = math.Max(newMaxBounds[1], v.Y())
			newMaxBounds[2] = math.Max(newMaxBounds[2], v.Z())
		}

		newCenter := vec.New(
			(newMinBounds[0]+newMaxBounds[0])/2,
			(newMinBounds[1]+newMaxBounds[1])/2,
			(newMinBounds[2]+newMaxBounds[2])/2,
		)

		fmt.Printf("=== AFTER CENTERING ===\n")
		fmt.Printf("New min bounds: [%f, %f, %f]\n", newMinBounds[0], newMinBounds[1], newMinBounds[2])
		fmt.Printf("New max bounds: [%f, %f, %f]\n", newMaxBounds[0], newMaxBounds[1], newMaxBounds[2])
		fmt.Printf("New center: [%f, %f, %f]\n", newCenter.X(), newCenter.Y(), newCenter.Z())

		if options.Position.X() != 0 || options.Position.Y() != 0 || options.Position.Z() != 0 {
			fmt.Printf("Shifted to requested position: [%f, %f, %f]\n",
				options.Position.X(), options.Position.Y(), options.Position.Z())
		}
	}

	// Second pass: read normals, texture coords, and faces
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "vn": // Normal
			if len(parts) < 4 {
				continue
			}

			nx, errX := strconv.ParseFloat(parts[1], 64)
			ny, errY := strconv.ParseFloat(parts[2], 64)
			nz, errZ := strconv.ParseFloat(parts[3], 64)

			if errX != nil || errY != nil || errZ != nil {
				continue
			}

			if options.FlipYZ {
				ny, nz = nz, ny
			}

			normal := vec.New(nx, ny, nz)
			// Normalize the normal vector
			length := math.Sqrt(nx*nx + ny*ny + nz*nz)
			if length > 0 {
				normal.ScaleInplace(length)
			}

			normals = append(normals, normal)

		case "f": // Face
			if len(parts) < 4 {
				continue
			}

			// Parse vertex/texture/normal indices
			var faceVertices []*vec.Vec3
			var faceNormals []*vec.Vec3

			for i := 1; i < len(parts); i++ {
				// Handle v/vt/vn format
				indices := strings.Split(parts[i], "/")

				// Vertex index (required)
				if len(indices) > 0 && indices[0] != "" {
					idx, err := strconv.Atoi(indices[0])
					if err != nil {
						continue
					}

					vIdx := fixIndex(idx, len(vertices))
					if vIdx >= 0 && vIdx < len(vertices) {
						faceVertices = append(faceVertices, vertices[vIdx])
					} else {
						continue
					}
				}

				// Normal index (optional)
				if len(indices) > 2 && indices[2] != "" && len(normals) > 0 && !options.IgnoreNormals {
					idx, err := strconv.Atoi(indices[2])
					if err == nil {
						nIdx := fixIndex(idx, len(normals))
						if nIdx >= 0 && nIdx < len(normals) {
							faceNormals = append(faceNormals, normals[nIdx])
						}
					}
				}
			}

			// Create triangles for the face (triangulate if needed)
			if len(faceVertices) >= 3 {
				// For a face with more than 3 vertices, we need to triangulate it
				for i := 2; i < len(faceVertices); i++ {
					v1, v2, v3 := faceVertices[0], faceVertices[i-1], faceVertices[i]

					// Optionally flip the winding order
					if options.FlipFaces {
						v2, v3 = v3, v2
					}

					// Create triangle with appropriate material
					if len(faceNormals) >= 3 && i < len(faceNormals) && !options.IgnoreNormals {
						// Use corresponding normals for the triangle vertices
						n1Idx := 0
						n2Idx := i - 1
						n3Idx := i

						// Safely get normals
						if n1Idx < len(faceNormals) && n2Idx < len(faceNormals) && n3Idx < len(faceNormals) {
							n1, n2, n3 := faceNormals[n1Idx], faceNormals[n2Idx], faceNormals[n3Idx]
							if options.FlipFaces {
								n2, n3 = n3, n2
							}

							// Create a triangle with custom normals
							triangles = append(triangles, hittable.NewTriangleWithNormals(
								[3]*vec.Vec3{v1, v2, v3},
								[3]*vec.Vec3{n1, n2, n3},
								mat,
							))
						} else {
							// Fall back to creating a triangle without custom normals
							triangles = append(triangles, hittable.NewTriangle(
								[3]*vec.Vec3{v1, v2, v3},
								mat,
							))
						}
					} else {
						// Create a triangle without custom normals
						triangles = append(triangles, hittable.NewTriangle(
							[3]*vec.Vec3{v1, v2, v3},
							mat,
						))
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file %s: %v", filename, err)
	}

	if options.Debug {
		fmt.Printf("=== MODEL SUMMARY ===\n")
		fmt.Printf("Loaded %d vertices, %d normals,  %d triangles\n",
			len(vertices), len(normals), len(triangles))
	}

	// Create a hittable list and add all triangles
	model := hittable.NewHittableList(len(triangles))
	for _, triangle := range triangles {
		model.Add(triangle)
	}

	// Build a Bounding Volume Hierarchy for faster ray intersection tests
	bvh := hittable.BuildBVH(model)

	// Final bbox check to verify positioning
	if options.Debug {
		bbox := bvh.BBox()
		if bbox != nil {
			fmt.Printf("=== FINAL BVH BOUNDS ===\n")
			fmt.Printf("X: %f to %f\n", bbox.AxisInterval(0).Min, bbox.AxisInterval(0).Max)
			fmt.Printf("Y: %f to %f\n", bbox.AxisInterval(1).Min, bbox.AxisInterval(1).Max)
			fmt.Printf("Z: %f to %f\n", bbox.AxisInterval(2).Min, bbox.AxisInterval(2).Max)

			// Calculate bbox center
			bboxCenter := vec.New(
				(bbox.AxisInterval(0).Min+bbox.AxisInterval(0).Max)/2,
				(bbox.AxisInterval(1).Min+bbox.AxisInterval(1).Max)/2,
				(bbox.AxisInterval(2).Min+bbox.AxisInterval(2).Max)/2,
			)
			fmt.Printf("BVH center: [%f, %f, %f]\n", bboxCenter.X(), bboxCenter.Y(), bboxCenter.Z())
		} else {
			fmt.Printf("Warning: BVH returned nil bbox\n")
		}
	}

	return bvh
}
