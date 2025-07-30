package geo

import (
	"math"

	"github.com/busyster996/geo/util"
)

// Circle represents a 2D circle with center coordinate and radius
type Circle struct {
	Center Coord // Center coordinate of the circle
	Radius int32 // Radius of the circle
}

// NewCirCle creates a new circle with given center and radius
func NewCirCle(center Coord, radius int32) Circle {
	return Circle{
		Center: center,
		Radius: radius,
	}
}

// GetLocationToBorder returns the relative position between circle and border
func (c *Circle) GetLocationToBorder(b *Border) LocationState {
	minX, minZ, maxX, maxZ := c.ToRect()
	return b.RectLocation(minX, minZ, maxX, maxZ)
}

// ToRect converts circle to bounding rectangle
// Returns the minimum and maximum X,Z coordinates of the bounding rectangle
func (c *Circle) ToRect() (minX, minZ, maxX, maxZ int32) {
	minX = c.Center.X - c.Radius
	minZ = c.Center.Z - c.Radius
	maxX = c.Center.X + c.Radius
	maxZ = c.Center.Z + c.Radius
	return
}

// GetIntersectCoord calculates the intersection point between circle and line from external point to center
// p: external point outside the circle
// Returns the intersection point on the circle
func (c *Circle) GetIntersectCoord(p Coord) Coord {
	// Return center directly if center and p coincide
	if c.Center == p {
		return c.Center
	}
	vec := NewVector(c.Center, p)
	length := vec.Length()
	ratio := float64(c.Radius) / length
	vec = vec.Trunc(ratio)
	return vec.ToCoord(c.Center)
}

// IsIntersect checks if line segment intersects with circle
func (c *Circle) IsIntersect(s *Segment) bool {
	_, ok := c.GetLineCross(s)
	return ok
}

// GetLineCross calculates intersection point between line segment and circle
// Returns the first intersection point if exists, otherwise returns false
func (c *Circle) GetLineCross(s *Segment) (Coord, bool) {
	var coord1 *Coord
	var coord2 *Coord
	fDis := CalDstCoordToCoord(s.A, s.B)

	dx := float64(s.B.X-s.A.X) / fDis
	dz := float64(s.B.Z-s.A.Z) / fDis

	ex := float64(c.Center.X - s.A.X)
	ez := float64(c.Center.Z - s.A.Z)

	a := ex*dx + ez*dz
	a2 := a * a
	e2 := ex*ex + ez*ez
	r2 := float64(c.Radius * c.Radius)
	if util.AC.Smaller(r2-e2+a2, 0) {
		return Coord{}, false
	}
	f := math.Sqrt(r2 - e2 + a2)
	t := a - f
	if t > -util.Epsilon && (t-fDis) < util.Epsilon {
		coord1 = &Coord{
			X: s.A.X + int32(t*dx),
			Z: s.A.Z + int32(t*dz),
		}
	}
	t = a + f
	if t > -util.Epsilon && (t-fDis) < util.Epsilon {
		coord2 = &Coord{
			X: s.A.X + int32(t*dx),
			Z: s.A.Z + int32(t*dz),
		}
	}

	if coord1 == nil {
		coord1 = coord2
	}
	if coord1 == nil {
		return Coord{}, false
	}
	return *coord1, true
}

// IsInterPolygon checks if circle intersects with polygon
// Returns true if intersected, false if separated
// Intersection cases: circle inside polygon, polygon inside circle, partial intersection
// This intersection detection can be used for collision detection
// circle: the circle to check
// vectors: polygon vector array in counter-clockwise order
// Reference: https://bitlush.com/blog/circle-vs-polygon-collision-detection-in-c-sharp
func (c *Circle) IsInterPolygon(vectors []Vector) bool {
	radiusSquared := float64(c.Radius * c.Radius)

	vertex := vectors[len(vectors)-1]
	center := NewVectorByCoord(c.Center)

	nearestDistance := math.MaxFloat64
	nearestIsInside := false
	nearestVertex := -1
	lastIsInside := false

	for i := 0; i < len(vectors); i++ {
		nextVertex := vectors[i]
		axis := center.Minus(&vertex)
		distance := axis.LengthSquared() - radiusSquared

		if util.AC.SmallerOrEqual(distance, 0) {
			return true
		}

		isInside := false
		edge := nextVertex.Minus(&vertex)
		edgeLengthSquared := edge.LengthSquared()
		if !util.AC.Equal(edgeLengthSquared, 0) {
			dot := edge.Dot(&axis)
			if util.AC.GreaterOrEqual(dot, 0) && util.AC.SmallerOrEqual(dot, edgeLengthSquared) {
				projection := edge.Trunc(dot / edgeLengthSquared)
				projection = vertex.Add(&projection)

				axis = projection.Minus(&center)
				if util.AC.SmallerOrEqual(axis.LengthSquared(), radiusSquared) {
					return true
				}

				if !isInsideEdge(&edge, &axis) {
					return false
				}

				isInside = true
			}
		}

		if util.AC.Smaller(distance, nearestDistance) {
			nearestDistance = distance
			nearestIsInside = isInside || lastIsInside
			nearestVertex = i
		}

		vertex = nextVertex
		lastIsInside = isInside
	}

	if nearestVertex == 0 {
		return nearestIsInside || lastIsInside
	}
	return nearestIsInside
}

// isInsideEdge checks if point is inside the edge for polygon intersection detection
func isInsideEdge(edge, axis *Vector) bool {
	switch {
	case edge.X > 0 && axis.Z > 0:
		return false
	case edge.X < 0 && axis.Z < 0:
		return false
	case edge.X == 0 && edge.Z > 0 && axis.X < 0:
		return false
	case edge.X == 0 && edge.Z <= 0 && axis.X > 0:
		return false
	}
	return true
}
