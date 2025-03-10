package hittable

import (
	"math"
	"math/rand"

	"github.com/nsp5488/go_raytracer/internal/aabb"
	"github.com/nsp5488/go_raytracer/internal/interval"
	"github.com/nsp5488/go_raytracer/internal/ray"
	"github.com/nsp5488/go_raytracer/internal/vec"
)

// Represents a sphere in 3D space
type sphere struct {
	Center *ray.Ray // Using a ray to represent motion

	Radius   float64
	Material Material
	bbox     *aabb.AABB
}

// Creates a new sphere
func NewSphere(center *vec.Vec3, radius float64, material Material) *sphere {
	rvec := vec.New(radius, radius, radius)
	bbox := aabb.FromPoints(center.Sub(rvec), center.Add(rvec))
	return &sphere{Center: ray.New(center, vec.Empty()), Radius: radius, Material: material, bbox: bbox}
}

// Creates a new sphere with motion blur
func NewMotionSphere(center1, center2 *vec.Vec3, radius float64, material Material) *sphere {
	rvec := vec.New(radius, radius, radius)
	center := *ray.New(center1, center2.Sub(center1))
	bbox1 := aabb.FromPoints(center.At(0).Sub(rvec), center.At(0).Add(rvec))
	bbox2 := aabb.FromPoints(center.At(1).Sub(rvec), center.At(1).Add(rvec))

	return &sphere{Center: &center, Radius: radius, Material: material, bbox: aabb.FromBBoxes(bbox1, bbox2)}
}
func (s *sphere) BBox() *aabb.AABB {
	return s.bbox
}

// Calculates the UV values of the ray intersection of a given sphere
// and stores them in (u, v)
func calculateSphereUV(point *vec.Vec3, u, v *float64) {
	theta := math.Acos(-point.Y())
	phi := math.Atan2(-point.Z(), point.X()) + math.Pi

	*u = phi / (2 * math.Pi)
	*v = theta / math.Pi
}

func (s *sphere) PdfValue(origin, direction *vec.Vec3) float64 {
	rec := &HitRecord{}
	if !s.Hit(ray.New(origin, direction), *interval.New(.0001, math.Inf(1)), rec) {
		return 0
	}
	distSquared := s.Center.At(0).Sub(origin).LengthSquared()
	cosThetaMax := math.Sqrt(1 - s.Radius*s.Radius/distSquared)
	solidAngle := 2 * math.Pi * (1 - cosThetaMax)

	return 1 / solidAngle
}
func (s *sphere) Random(origin *vec.Vec3) *vec.Vec3 {
	direction := s.Center.At(0).Sub(origin)
	distSquared := direction.LengthSquared()
	onb := NewONB(direction)

	return onb.Transform(randomToSphere(s.Radius, distSquared))
}
func randomToSphere(radius, distSquared float64) *vec.Vec3 {
	r1 := rand.Float64()
	r2 := rand.Float64()
	z := 1 + r2*(math.Sqrt(1-radius*radius/distSquared)-1)
	phi := 2 * math.Pi * r1

	t := math.Sqrt(1 - z*z)
	x := math.Cos(phi) * t
	y := math.Sin(phi) * t
	return vec.New(x, y, z)
}

// Hit checks if a ray intersects with the sphere.
func (s *sphere) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	curCenter := s.Center.At(r.Time())
	oc := curCenter.Sub(r.Origin())

	// a = direction * direction = len(direction)^2
	a := r.Direction().LengthSquared()
	// h = r * oc
	h := r.Direction().Dot(oc)
	// c = oc * oc - radius^2 = len(oc)^2 - radius^2
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}

	sqrtd := math.Sqrt(discriminant)
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}

	record.t = root
	record.p = r.At(root)
	outward_normal := record.p.Sub(curCenter).Scale(1 / s.Radius)
	record.setFaceNormal(r, outward_normal)
	record.Material = s.Material
	calculateSphereUV(outward_normal, &record.u, &record.v)
	return true
}

type quad struct {
	Q        *vec.Vec3 // One corner of the plane
	u        *vec.Vec3 // u,v are vectors that point from Q to two other corners
	v        *vec.Vec3
	normal   *vec.Vec3 // normal = unit(u x v)
	w        *vec.Vec3
	D        float64 // D = Ax + By + Cz = dot(Q, normal)
	area     float64
	bbox     *aabb.AABB
	material Material
}

func NewQuad(Q, u, v *vec.Vec3, material Material) *quad {
	q := &quad{Q: Q, u: u, v: v, material: material}

	n := u.Cross(v)
	q.area = n.Length()
	q.normal = n.UnitVector()
	q.D = q.normal.Dot(Q)
	q.w = n.Scale(1 / n.Dot(n))

	q.setBBox()
	return q
}

func (q *quad) setBBox() {
	diag1 := aabb.FromPoints(q.Q, q.Q.Add(q.u).Add(q.v))
	diag2 := aabb.FromPoints(q.Q.Add(q.u), q.Q.Add(q.v))
	q.bbox = aabb.FromBBoxes(diag1, diag2)
}

func (q *quad) BBox() *aabb.AABB {
	return q.bbox
}

func (q *quad) PdfValue(origin, direction *vec.Vec3) float64 {
	record := &HitRecord{}
	if !q.Hit(ray.New(origin, direction), *interval.New(0.001, math.Inf(1)), record) {
		return 0
	}
	distSquared := record.t * record.t * direction.LengthSquared()
	cosine := math.Abs(direction.Dot(record.normal) / direction.Length())
	return distSquared / (cosine * q.area)
}
func (q *quad) Random(origin *vec.Vec3) *vec.Vec3 {
	p := q.Q.Add(q.u.Scale(rand.Float64())).Add(q.v.Scale(rand.Float64()))

	return p.Sub(origin)
}

func (q *quad) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	denom := q.normal.Dot(r.Direction())

	// Low values in denominator -> ray is parallel to the plane
	if math.Abs(denom) < 1e-8 {
		return false
	}

	// If t lands outside of our search interval
	t := (q.D - q.normal.Dot(r.Origin())) / denom
	if !rayT.Contains(t) {
		return false
	}

	intersection := r.At(t)

	// Check that the ray intersects the quad itself, not just the plane
	planarHitpoint := intersection.Sub(q.Q)
	alpha := q.w.Dot(planarHitpoint.Cross(q.v))
	beta := q.w.Dot(q.u.Cross(planarHitpoint))
	if !isInterior(alpha, beta, record) {
		return false
	}

	record.t = t
	record.p = intersection
	record.Material = q.material
	record.setFaceNormal(r, q.normal)
	return true
}

func isInterior(alpha, beta float64, record *HitRecord) bool {
	if !interval.Unit().Contains(alpha) || !interval.Unit().Contains(beta) {
		return false
	}

	record.u = alpha
	record.v = beta
	return true
}

func NewBox(a, b *vec.Vec3, mat Material) Hittable {
	sides := NewHittableList(6)

	minVec := vec.New(
		min(a.X(), b.X()),
		min(a.Y(), b.Y()),
		min(a.Z(), b.Z()),
	)
	maxVec := vec.New(
		max(a.X(), b.X()),
		max(a.Y(), b.Y()),
		max(a.Z(), b.Z()),
	)

	dx := vec.New(maxVec.X()-minVec.X(), 0, 0)
	dy := vec.New(0, maxVec.Y()-minVec.Y(), 0)
	dz := vec.New(0, 0, maxVec.Z()-minVec.Z())

	// front
	sides.Add(NewQuad(vec.New(minVec.X(), minVec.Y(), maxVec.Z()), dx, dy, mat))
	// right
	sides.Add(NewQuad(vec.New(maxVec.X(), minVec.Y(), maxVec.Z()), dz.Negate(), dy, mat))
	// back
	sides.Add(NewQuad(vec.New(maxVec.X(), minVec.Y(), minVec.Z()), dx.Negate(), dy, mat))
	// left
	sides.Add(NewQuad(vec.New(minVec.X(), minVec.Y(), minVec.Z()), dz, dy, mat))
	// top
	sides.Add(NewQuad(vec.New(minVec.X(), maxVec.Y(), maxVec.Z()), dx, dz.Negate(), mat))
	// bottom
	sides.Add(NewQuad(vec.New(minVec.X(), minVec.Y(), minVec.Z()), dx, dz, mat))

	return BuildBVH(sides)
}

type Triangle struct {
	defaultPdfImpl
	Vertices [3]*vec.Vec3
	Normals  [3]*vec.Vec3 // Vertex normals from the OBJ file
	normal   *vec.Vec3    // Face normal (calculated from vertices)
	area     float64
	bbox     *aabb.AABB
	Material Material

	texCoords        [3][2]float64 // Contains the texture coordinates for each vertex, if hasUV is true, else nil
	hasUV            bool
	hasVertexNormals bool // Whether this triangle uses per-vertex normals
}

// NewTriangle creates a triangle with a calculated face normal
func NewTriangle(vertices [3]*vec.Vec3, material Material) *Triangle {
	t := &Triangle{
		Vertices:         vertices,
		Material:         material,
		hasVertexNormals: false,
	}
	// Calc area
	edge1 := t.Vertices[1].Sub(t.Vertices[0])
	edge2 := t.Vertices[2].Sub(t.Vertices[0])
	crossProduct := edge1.Cross(edge2)
	t.area = crossProduct.Length() / 2.0

	// Calculate face normal from vertices
	e0 := t.Vertices[1].Sub(t.Vertices[0])
	e1 := t.Vertices[2].Sub(t.Vertices[0])
	t.normal = e0.Cross(e1).UnitVector()
	t.hasUV = false

	t.SetBbox()
	return t
}

// NewTriangleWithNormals creates a triangle with custom vertex normals
func NewTriangleWithNormals(vertices [3]*vec.Vec3, normals [3]*vec.Vec3, material Material) *Triangle {
	t := &Triangle{
		Vertices:         vertices,
		Normals:          normals,
		Material:         material,
		hasVertexNormals: true,
	}
	// Calculate area
	edge1 := t.Vertices[1].Sub(t.Vertices[0])
	edge2 := t.Vertices[2].Sub(t.Vertices[0])
	crossProduct := edge1.Cross(edge2)
	t.area = crossProduct.Length() / 2.0

	// Still calculate face normal for fallback/bbox calculations
	e0 := t.Vertices[1].Sub(t.Vertices[0])
	e1 := t.Vertices[2].Sub(t.Vertices[0])
	t.normal = e0.Cross(e1).UnitVector()
	t.hasUV = false
	t.SetBbox()
	return t
}

// NewTexturedTriangle creates a new triangle with texture coordinates
func NewTexturedTriangle(vertices [3]*vec.Vec3, texCoords [3][2]float64, material Material) *Triangle {
	t := NewTriangle(vertices, material)
	t.texCoords = texCoords
	t.hasUV = true
	return t
}

// NewTexturedTriangleWithNormals creates a new triangle with custom normals and texture coordinates
func NewTexturedTriangleWithNormals(vertices [3]*vec.Vec3, normals [3]*vec.Vec3, texCoords [3][2]float64, material Material) *Triangle {
	t := NewTriangleWithNormals(vertices, normals, material)
	t.texCoords = texCoords
	t.hasUV = true
	return t
}
func (t *Triangle) SetBbox() {
	minX := math.Inf(1)
	maxX := math.Inf(-1)
	minY := math.Inf(1)
	maxY := math.Inf(-1)
	minZ := math.Inf(1)
	maxZ := math.Inf(-1)

	// identify the bounding interval (min,max) across all 3 dimensions
	for vert := range 3 {
		minX = min(t.Vertices[vert].X(), minX)
		maxX = max(t.Vertices[vert].X(), maxX)
		minY = min(t.Vertices[vert].Y(), minY)
		maxY = max(t.Vertices[vert].Y(), maxY)
		minZ = min(t.Vertices[vert].Z(), minZ)
		maxZ = max(t.Vertices[vert].Z(), maxZ)
	}

	// Add a small epsilon to avoid degenerate boxes
	const epsilon = 1e-8
	if maxX-minX < epsilon {
		maxX += epsilon
		minX -= epsilon
	}
	if maxY-minY < epsilon {
		maxY += epsilon
		minY -= epsilon
	}
	if maxZ-minZ < epsilon {
		maxZ += epsilon
		minZ -= epsilon
	}

	xInt := interval.New(minX, maxX)
	yInt := interval.New(minY, maxY)
	zInt := interval.New(minZ, maxZ)
	t.bbox = aabb.NewAABB(xInt, yInt, zInt)
}

func (t *Triangle) PdfValue(origin, direction *vec.Vec3) float64 {
	record := &HitRecord{}
	if !t.Hit(ray.New(origin, direction), *interval.New(0.001, math.Inf(1)), record) {
		return 0
	}

	// Calculate the area of the triangle
	distSquared := record.t * record.t * direction.LengthSquared()
	cosine := math.Abs(direction.Dot(record.normal) / direction.Length())

	return distSquared / (cosine * t.area)
}

func (t *Triangle) Random(origin *vec.Vec3) *vec.Vec3 {
	// Use barycentric coordinates for random point generation
	// Generate two random values r1, r2 where r1 + r2 <= 1
	r1 := rand.Float64()
	r2 := rand.Float64() * (1 - r1)

	// Barycentric coordinates (1-r1-r2, r1, r2)
	a := 1 - r1 - r2
	b := r1
	c := r2

	// Calculate the random point on the triangle
	p := t.Vertices[0].Scale(a).Add(t.Vertices[1].Scale(b)).Add(t.Vertices[2].Scale(c))

	// Return direction from origin to point
	return p.Sub(origin)
}

// interpolateNormal calculates the interpolated normal at the hit point
// using barycentric coordinates (u,v)
func (t *Triangle) interpolateNormal(u, v float64) *vec.Vec3 {
	if !t.hasVertexNormals {
		return t.normal
	}

	// Barycentric coordinates: u, v, and w=1-u-v
	w := 1.0 - u - v

	// Interpolate normals using barycentric coordinates
	nx := w*t.Normals[0].X() + u*t.Normals[1].X() + v*t.Normals[2].X()
	ny := w*t.Normals[0].Y() + u*t.Normals[1].Y() + v*t.Normals[2].Y()
	nz := w*t.Normals[0].Z() + u*t.Normals[1].Z() + v*t.Normals[2].Z()

	// Return normalized interpolated normal
	interpolated := vec.New(nx, ny, nz)
	return interpolated.UnitVector()
}

// Muller-Trumbore implementation
func (t *Triangle) Hit(r *ray.Ray, rayT interval.Interval, record *HitRecord) bool {
	e0 := t.Vertices[1].Sub(t.Vertices[0])
	e1 := t.Vertices[2].Sub(t.Vertices[0])

	pvec := r.Direction().Cross(e1)
	det := e0.Dot(pvec)

	if math.Abs(det) < 1e-8 {
		return false // Ray is parallel to the triangle
	}

	invDet := 1.0 / det
	tvec := r.Origin().Sub(t.Vertices[0])
	u := tvec.Dot(pvec) * invDet
	if u < 0 || u > 1 {
		return false
	}

	qvec := tvec.Cross(e0)
	v := r.Direction().Dot(qvec) * invDet
	if v < 0 || (u+v) > 1 {
		return false
	}

	tl := e1.Dot(qvec) * invDet
	if tl < rayT.Min || tl > rayT.Max {
		return false
	}

	if t.hasUV {
		w := (1 - u - v)
		// Interpolate texture coordinates using barycentric coordinates
		record.u = w*t.texCoords[0][0] + u*t.texCoords[1][0] + v*t.texCoords[2][0]
		record.v = w*t.texCoords[0][1] + u*t.texCoords[1][1] + v*t.texCoords[2][1]
	} else {
		// Default UVs, these may cause weird artifacts
		record.u = u
		record.v = v
	}
	record.t = tl
	record.p = r.At(tl)

	// Use the interpolated normal if we have vertex normals
	if t.hasVertexNormals {
		interpolatedNormal := t.interpolateNormal(u, v)
		record.setFaceNormal(r, interpolatedNormal)
	} else {
		record.setFaceNormal(r, t.normal)
	}

	record.Material = t.Material

	return true
}

func (t *Triangle) BBox() *aabb.AABB {
	return t.bbox
}
