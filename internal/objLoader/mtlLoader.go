package objLoader

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nsp5488/go_raytracer/internal/hittable"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// MtlMaterial represents a material defined in an MTL file
type MtlMaterial struct {
	Name       string
	Ambient    *vec.Vec3 // Ka
	Diffuse    *vec.Vec3 // Kd
	Specular   *vec.Vec3 // Ks
	Emission   *vec.Vec3 // Ke
	Tf         *vec.Vec3
	SpecExp    float64           // Ns (specular exponent)
	Dissolve   float64           // d (transparency)
	Refraction float64           // Ni (index of refraction)
	Illum      int               // illumination model
	MapKd      string            // diffuse texture map
	MapKa      string            // ambient texture map
	MapKs      string            // specular texture map  TODO
	MapNs      string            // specular highlight map  TODO
	MapBump    string            // bump map TODO
	Material   hittable.Material // Converted raytracer material
}

// MaterialLibrary stores materials loaded from an MTL file
type MaterialLibrary struct {
	Materials map[string]*MtlMaterial
	BaseDir   string // Directory of MTL file for resolving texture paths
	Debug     bool
}

// NewMaterialLibrary creates a new empty material library
func NewMaterialLibrary(debug bool) *MaterialLibrary {
	return &MaterialLibrary{
		Materials: make(map[string]*MtlMaterial),
		Debug:     debug,
	}
}

// LoadMTL loads an MTL file and parses its materials
func LoadMTL(mtlPath string, debug bool) (*MaterialLibrary, error) {
	lib := NewMaterialLibrary(debug)
	lib.BaseDir = filepath.Dir(mtlPath)

	file, err := os.Open(mtlPath)
	if err != nil {
		return nil, fmt.Errorf("could not open MTL file %s: %v", mtlPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var currentMaterial *MtlMaterial

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
		case "newmtl":
			if len(parts) < 2 {
				log.Printf("Line %d: Malformed newmtl directive: %s", lineNum, line)
				continue
			}

			name := parts[1]
			currentMaterial = &MtlMaterial{
				Name:       name,
				Ambient:    vec.New(0.2, 0.2, 0.2), // Default values
				Diffuse:    vec.New(0.8, 0.8, 0.8),
				Specular:   vec.New(0.0, 0.0, 0.0),
				Emission:   vec.New(0.0, 0.0, 0.0),
				Tf:         vec.New(0, 0, 0),
				SpecExp:    0.0,
				Dissolve:   1.0, // Fully opaque by default
				Refraction: 1.0, // Air by default
				Illum:      2,   // Default illumination model
			}
			lib.Materials[name] = currentMaterial

		case "Ka":
			if currentMaterial == nil || len(parts) < 4 {
				continue
			}
			r, _ := strconv.ParseFloat(parts[1], 64)
			g, _ := strconv.ParseFloat(parts[2], 64)
			b, _ := strconv.ParseFloat(parts[3], 64)
			currentMaterial.Ambient = vec.New(r, g, b)

		case "Kd":
			if currentMaterial == nil || len(parts) < 4 {
				continue
			}
			r, _ := strconv.ParseFloat(parts[1], 64)
			g, _ := strconv.ParseFloat(parts[2], 64)
			b, _ := strconv.ParseFloat(parts[3], 64)
			currentMaterial.Diffuse = vec.New(r, g, b)

		case "Ks":
			if currentMaterial == nil || len(parts) < 4 {
				continue
			}
			r, _ := strconv.ParseFloat(parts[1], 64)
			g, _ := strconv.ParseFloat(parts[2], 64)
			b, _ := strconv.ParseFloat(parts[3], 64)
			currentMaterial.Specular = vec.New(r, g, b)

		case "Ke":
			if currentMaterial == nil || len(parts) < 4 {
				continue
			}
			r, _ := strconv.ParseFloat(parts[1], 64)
			g, _ := strconv.ParseFloat(parts[2], 64)
			b, _ := strconv.ParseFloat(parts[3], 64)
			currentMaterial.Emission = vec.New(r, g, b)

		case "Ns":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			val, _ := strconv.ParseFloat(parts[1], 64)
			currentMaterial.SpecExp = val

		case "d":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			val, _ := strconv.ParseFloat(parts[1], 64)
			currentMaterial.Dissolve = val

		case "Ni":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			val, _ := strconv.ParseFloat(parts[1], 64)
			currentMaterial.Refraction = val
		case "Tf":
			if currentMaterial == nil || len(parts) < 4 {
				continue
			}
			r, _ := strconv.ParseFloat(parts[1], 64)
			g, _ := strconv.ParseFloat(parts[2], 64)
			b, _ := strconv.ParseFloat(parts[3], 64)
			// Average the transparency values
			currentMaterial.Tf = vec.New(r, g, b)
			currentMaterial.Dissolve = ((r + g + b) / 3.0)
		case "illum":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			val, _ := strconv.Atoi(parts[1])
			currentMaterial.Illum = val

		case "map_Kd":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			currentMaterial.MapKd = strings.Join(parts[1:], " ")

		case "map_Ka":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			currentMaterial.MapKa = strings.Join(parts[1:], " ")

		case "map_Ks":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			currentMaterial.MapKs = strings.Join(parts[1:], " ")

		case "map_Ns":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			currentMaterial.MapNs = strings.Join(parts[1:], " ")

		case "map_bump", "bump":
			if currentMaterial == nil || len(parts) < 2 {
				continue
			}
			currentMaterial.MapBump = strings.Join(parts[1:], " ")
		}
	}

	// Convert all materials to raytracer materials
	for _, mtl := range lib.Materials {
		mtl.Material = ConvertToRaytracerMaterial(mtl)
	}

	if lib.Debug {
		fmt.Printf("=== MTL SUMMARY ===\n")
		fmt.Printf("Loaded %d materials from %s\n", len(lib.Materials), mtlPath)
		for name, mtl := range lib.Materials {
			fmt.Printf("  Material '%s':\n", name)
			fmt.Printf("    Diffuse: [%f, %f, %f]\n", mtl.Diffuse.X(), mtl.Diffuse.Y(), mtl.Diffuse.Z())
			if mtl.MapKd != "" {
				fmt.Printf("    Diffuse Map: %s\n", mtl.MapKd)
			}
			if mtl.Dissolve < 1.0 {
				fmt.Printf("    Transparency: %f\n", 1.0-mtl.Dissolve)
			}
			if mtl.Refraction > 1.0 {
				fmt.Printf("    Refraction Index: %f\n", mtl.Refraction)
			}
		}
	}

	return lib, nil
}

// ConvertToRaytracerMaterial converts an MTL material to a raytracer material
func ConvertToRaytracerMaterial(mtl *MtlMaterial) hittable.Material {
	// First, handle special case materials

	// 1. Handle completely transparent or refractive materials (glass, water, etc.)
	if mtl.Dissolve < 0.95 && mtl.Refraction > 1.0 || mtl.Illum == 4 || mtl.Illum == 6 || mtl.Illum == 7 {
		// Dielectric material with proper refraction index
		refractiveIndex := mtl.Refraction
		if refractiveIndex <= 1.01 { // If Ni is suspiciously close to air
			refractiveIndex = 1.5 // Default glass
		}
		return hittable.NewDielectric(refractiveIndex)
	}

	// For isotropic (volume scattering) materials
	if mtl.Dissolve < 0.95 {
		return hittable.NewIsotropic(mtl.Diffuse)
	}

	// 3. Handle emissive materials (light sources)
	emissiveIntensity := mtl.Emission.X() + mtl.Emission.Y() + mtl.Emission.Z()
	if emissiveIntensity > 0.1 {
		if mtl.MapKd != "" {
			tex := hittable.NewImageTexture(mtl.MapKd)
			return hittable.NewDiffuseLightTextured(tex)
		} else if mtl.MapKa != "" {
			tex := hittable.NewImageTexture(mtl.MapKa)
			return hittable.NewDiffuseLightTextured(tex)
		}
		return hittable.NewDiffuseLight(mtl.Emission)
	}
	// 4. Handle metal-like materials
	// Detect metals by analyzing specular and diffuse properties
	specIntensity := mtl.Specular.X() + mtl.Specular.Y() + mtl.Specular.Z()
	diffuseIntensity := mtl.Diffuse.X() + mtl.Diffuse.Y() + mtl.Diffuse.Z()

	isMetallic := specIntensity > 0.1 && specIntensity > diffuseIntensity*0.5

	if isMetallic {
		// Convert specular exponent to roughness (inverse relationship)
		// Typical specular exponents range from 0 to 1000, with higher values being smoother
		roughness := 0.0
		if mtl.SpecExp <= 0.0 {
			roughness = 1.0
		} else if mtl.SpecExp >= 1000.0 {
			roughness = 0.0 // Perfectly smooth
		} else {
			// Non-linear mapping from spec exponent to roughness
			roughness = math.Pow(1.0-mtl.SpecExp/1000.0, 2.0)
			roughness = math.Max(0.0, math.Min(1.0, roughness))
		}

		// Use specular color as the metal color, but if it's too dark,
		// blend with diffuse to get a more reasonable metal color
		metalColor := mtl.Specular
		if specIntensity < 0.2 {
			// Blend with diffuse for better visual results
			blend := 1.0 - (specIntensity / 0.2)
			metalColor = vec.New(
				(1.0-blend)*mtl.Specular.X()+blend*mtl.Diffuse.X(),
				(1.0-blend)*mtl.Specular.Y()+blend*mtl.Diffuse.Y(),
				(1.0-blend)*mtl.Specular.Z()+blend*mtl.Diffuse.Z(),
			)
		}

		return hittable.NewMetal(metalColor, roughness)
	}

	// 5. Handle illumination model specific conversions
	switch mtl.Illum {
	case 0, 1, 2:
		// For standard diffuse materials
		if mtl.MapKd != "" {
			return hittable.NewTexturedLambertian(hittable.NewImageTexture(mtl.MapKd))
		} else if mtl.MapKa != "" {
			// Use ambient map as fallback
			return hittable.NewTexturedLambertian(hittable.NewImageTexture(mtl.MapKa))
		}
		return hittable.NewLambertian(mtl.Diffuse)

	case 3, 4, 5:
		roughness := 0.3 // Medium roughness as default
		return hittable.NewMetal(mtl.Specular, roughness)

	default:
		// Fallback on diffuse
		if mtl.MapKd != "" {
			return hittable.NewTexturedLambertian(hittable.NewImageTexture(mtl.MapKd))
		} else if mtl.MapKa != "" {
			// Use ambient map as fallback
			return hittable.NewTexturedLambertian(hittable.NewImageTexture(mtl.MapKa))
		}
		return hittable.NewLambertian(mtl.Diffuse)
	}
}
