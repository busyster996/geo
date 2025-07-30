package geo

import "math"

// Coord represents a 2D coordinate with X and Z components
type Coord struct {
	X int32 `bson:"x"` // X coordinate
	Z int32 `bson:"z"` // Z coordinate
}

// NewCoord creates a new coordinate with given x and z values
func NewCoord(x, z int32) Coord {
	return Coord{
		X: x,
		Z: z,
	}
}

// IsEqual checks if this coordinate equals the target coordinate
func (c *Coord) IsEqual(target Coord) bool {
	return *c == target
}

// GetLocationToBorder returns the positional relationship between coordinate and border
func (c *Coord) GetLocationToBorder(b *Border) LocationState {
	return b.CoordLocation(*c)
}

// CalDstCoordToCoord calculates the distance between two coordinates
func CalDstCoordToCoord(coord1, coord2 Coord) float64 {
	dx := int64(coord1.X - coord2.X)
	dz := int64(coord1.Z - coord2.Z)

	squared := dx*dx + dz*dz
	return math.Sqrt(float64(squared))
}

// CalDstCoordToCoordWithoutSqrt calculates the squared distance between two coordinates (without square root)
func CalDstCoordToCoordWithoutSqrt(coord1, coord2 Coord) float64 {
	dx := int64(coord1.X - coord2.X)
	dz := int64(coord1.Z - coord2.Z)

	squared := dx*dx + dz*dz
	return float64(squared)
}
