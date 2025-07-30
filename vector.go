package geo

import "math"

// Vector represents a position vector
// A vector from the coordinate origin to a point's position is called a position vector
// Reference: https://zh.wikipedia.org/wiki/%E4%BD%8D%E7%BD%AE%E5%90%91%E9%87%8F
type Vector struct {
	X, Z int32 // X and Z components of the vector
}

// NewVector creates a new vector from start to end coordinate
func NewVector(start, end Coord) Vector {
	return Vector{
		X: end.X - start.X,
		Z: end.Z - start.Z,
	}
}

// NewVectorByCoord creates a new vector from a coordinate (treated as position vector)
func NewVectorByCoord(p Coord) Vector {
	return Vector(p)
}

// Add performs vector addition
func (v *Vector) Add(vec *Vector) Vector {
	return Vector{
		X: v.X + vec.X,
		Z: v.Z + vec.Z,
	}
}

// Minus performs vector subtraction
func (v *Vector) Minus(vec *Vector) Vector {
	return Vector{
		X: v.X - vec.X,
		Z: v.Z - vec.Z,
	}
}

// Dot calculates the dot product of two vectors
func (v *Vector) Dot(vec *Vector) float64 {
	result := int64(v.X) * int64(vec.X)
	result += int64(v.Z) * int64(vec.Z)
	return float64(result)
}

// Cross calculates the cross product of two vectors (2D cross product returns scalar)
// result > 0: vec is on the left side of v
// result < 0: vec is on the right side of v
// result = 0: vectors are collinear
func (v *Vector) Cross(vec *Vector) int64 {
	return int64(v.X)*int64(vec.Z) - int64(v.Z)*int64(vec.X)
}

// LengthSquared returns the squared length of the vector
func (v *Vector) LengthSquared() float64 {
	return float64(int64(v.X)*int64(v.X) + int64(v.Z)*int64(v.Z))
}

// Length returns the length (magnitude) of the vector
func (v *Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Trunc scales the vector by given ratio
func (v *Vector) Trunc(ratio float64) Vector {
	return Vector{
		X: int32(ratio * float64(v.X)),
		Z: int32(ratio * float64(v.Z)),
	}
}

// TruncEdge truncates edge to a unit vector with length 1000
func TruncEdge(start, end Coord) Coord {
	// Generate inflection point array
	vec0 := NewVector(start, end)
	// Unit vector with length 1000
	vec0 = vec0.Trunc(1000 / vec0.Length())

	return vec0.ToCoord(start)
}

// ToCoord converts vector to coordinate by adding to start coordinate
func (v Vector) ToCoord(start Coord) Coord {
	return Coord{
		X: start.X + v.X,
		Z: start.Z + v.Z,
	}
}

// Rotate rotates the vector by given angle and returns new vector
// Left-hand coordinate system, default rotation is to the left
// For any two different points A and B, A rotated θ angle around B results in:
// (Δx*cosθ - Δy*sinθ + xB, Δy*cosθ + Δx*sinθ + yB)
// Note: xB, yB are coordinates of point B
// Reference: https://blog.csdn.net/u013445530/article/details/44904017
func (v *Vector) Rotate(angle float64) Vector {
	x0 := float64(v.X)
	z0 := float64(v.Z)

	cos := math.Cos(angle)
	sin := math.Sin(angle)

	x := x0*cos - z0*sin
	z := x0*sin + z0*cos

	return Vector{
		X: int32(x),
		Z: int32(z),
	}
}

// CalCoordDst calculates the distance from a point to the vector (line)
func (v *Vector) CalCoordDst(start, target Coord) float64 {
	vec := NewVector(start, target)

	angle := v.GetAngle(&vec)
	return vec.Length() * math.Sin(angle)
}

// GetAngle calculates the angle between two vectors
func (v *Vector) GetAngle(vec *Vector) float64 {
	a := v.Dot(vec)
	b := v.Length() * vec.Length()
	t := a / b
	angle := math.Acos(t)
	if math.IsNaN(angle) {
		if t > 0 {
			return 0
		}
		return math.Pi
	}

	return angle
}

// cross calculates cross product for three points
// Result = p3p1 X p3p2
func cross(p1, p2, p3 Coord) int64 {
	s := int64(p1.X-p3.X)*int64((p2.Z-p3.Z)) - int64(p2.X-p3.X)*int64(p1.Z-p3.Z)
	return s
}

// CalCoordByRatio calculates coordinate by ratio
// ratio represents the ratio between newVec and Vec lengths
// Returns the endpoint coordinate of newVec
func CalCoordByRatio(startCoord, endCoord Coord, ratio float64) Coord {
	vec := NewVector(startCoord, endCoord)
	return vec.Trunc(ratio).ToCoord(startCoord)
}
