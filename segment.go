package geo

import "github.com/busyster996/geo/util"

// Segment represents a line segment defined by two endpoints
type Segment struct {
	A, B Coord // Start and end points of the segment
}

// NewSegment creates a new line segment
func NewSegment(a, b Coord) Segment {
	return Segment{
		A: a,
		B: b,
	}
}

// ToVector converts segment to vector representation
func (s *Segment) ToVector() Vector {
	return NewVector(s.A, s.B)
}

// CalCoordDst calculates the distance from a point to the line segment
func (s *Segment) CalCoordDst(coord Coord) float64 {
	// Distance from point to segment endpoints
	a := CalDstCoordToCoord(coord, s.A)
	b := CalDstCoordToCoord(coord, s.B)
	dst := min(a, b)

	// Check if perpendicular from point intersects with segment
	ab := NewVector(s.A, s.B)
	ap := NewVector(s.A, coord)
	if lab := ab.Length(); util.AC.Greater(ap.Dot(&ab)/lab, lab) {
		return dst
	}

	// Distance from point to line
	vec := s.ToVector()
	c := vec.CalCoordDst(s.A, coord)

	return min(dst, c)
}

// Pan translates the segment parallel by given distance
// Left-hand coordinate system, moves in direction of positive cross product
func (s *Segment) Pan(dst int32, positive bool) Segment {
	v := s.ToVector()
	// Normal vector in direction of positive cross product
	normalV := NewVectorByCoord(Coord{X: -v.Z, Z: v.X})
	if !positive {
		normalV = NewVectorByCoord(Coord{X: v.Z, Z: -v.X})
	}
	ratio := float64(dst) / v.Length()
	newV := normalV.Trunc(ratio)

	return NewSegment(newV.ToCoord(s.A), newV.ToCoord(s.B))
}

// CrossCircle checks if segment intersects with circle, returns first intersection point if exists
func (s *Segment) CrossCircle(circle Circle) (Coord, bool) {
	return GetLineCrossCircle(s.A, s.B, circle.Center, circle.Radius)
}

// IsRectCross performs quick rejection test for line segment intersection
func IsRectCross(p0, p1, q0, q1 Coord) bool {
	ret := min(p0.X, p1.X) <= max(q0.X, q1.X) &&
		min(q0.X, q1.X) <= max(p0.X, p1.X) &&
		min(p0.Z, p1.Z) <= max(q0.Z, q1.Z) &&
		min(q0.Z, q1.Z) <= max(p0.Z, p1.Z)
	return ret
}

// IsLineSegmentCross performs straddle test for line segment intersection
func IsLineSegmentCross(p0, p1, q0, q1 Coord) bool {
	// q0q1 X q0p0
	b1 := cross(q1, p0, q0)
	// q0q1 X q0p1
	b2 := cross(q1, p1, q0)

	// Cross product equals 0 means one point is collinear with the other segment
	if b1 == 0 || b2 == 0 {
		return true
	}

	// p0p1 X p0q0
	a1 := cross(p1, q0, p0)
	// p0p1 X p0q1
	a2 := cross(p1, q1, p0)

	if a1 == 0 || a2 == 0 {
		return true
	}

	return ((b1 < 0) != (b2 < 0)) && ((a1 < 0) != (a2 < 0))
}

// GetCrossCoord calculates intersection point between line segments P1P2 and Q1Q2
// Reference: https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect/565282#
func GetCrossCoord(p0, p1, q0, q1 Coord) (Coord, bool) {
	v1 := NewVector(p0, p1)
	v2 := NewVector(q0, q1)
	// Parallel lines
	if v1.Cross(&v2) == 0 {
		return Coord{}, false
	}

	if IsRectCross(p0, p1, q0, q1) {
		if IsLineSegmentCross(p0, p1, q0, q1) {
			// Calculate intersection point
			s1X := float64(p1.X - p0.X)
			s1Z := float64(p1.Z - p0.Z)
			s2X := float64(q1.X - q0.X)
			s2Z := float64(q1.Z - q0.Z)

			t := (s2X*float64(p0.Z-q0.Z) - s2Z*float64(p0.X-q0.X)) / (-s2X*s1Z + s1X*s2Z)

			return Coord{X: p0.X + int32(t*s1X), Z: p0.Z + int32(t*s1Z)}, true
		}
	}
	return Coord{}, false
}
