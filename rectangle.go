package geo

import "math/rand"

// Rectangle represents a rectangle defined by bottom-left corner and dimensions
type Rectangle struct {
	Coord        // Bottom-left corner coordinate
	Width  int32 // Rectangle width
	Height int32 // Rectangle height
}

// NewRectangle creates a new rectangle with given position and dimensions
func NewRectangle(x, z, width, height int32) Rectangle {
	return Rectangle{
		Coord: Coord{
			X: x,
			Z: z,
		},
		Width:  width,
		Height: height,
	}
}

// RandCoord generates a random coordinate within the rectangle
func (rec *Rectangle) RandCoord() Coord {
	return Coord{
		X: rec.X + rand.Int31n(rec.Width),
		Z: rec.Z + rand.Int31n(rec.Height),
	}
}

// GetVerticeCoords returns the 4 vertex coordinates in counter-clockwise order
func (rec *Rectangle) GetVerticeCoords() [4]Coord {
	var p [4]Coord
	p[0] = Coord{X: rec.Coord.X, Z: rec.Coord.Z}                          // Bottom-left
	p[1] = Coord{X: rec.Coord.X + rec.Width, Z: rec.Coord.Z}              // Bottom-right
	p[2] = Coord{X: rec.Coord.X + rec.Width, Z: rec.Coord.Z + rec.Height} // Top-right
	p[3] = Coord{X: rec.Coord.X, Z: rec.Coord.Z + rec.Height}             // Top-left
	return p
}

// GetVectors returns 4 edge vectors in counter-clockwise order
func (rec *Rectangle) GetVectors() [4]Vector {
	coords := rec.GetVerticeCoords()
	return [4]Vector{
		NewVectorByCoord(coords[0]),
		NewVectorByCoord(coords[1]),
		NewVectorByCoord(coords[2]),
		NewVectorByCoord(coords[3]),
	}
}

// GetLocationToBorder returns the positional relationship between rectangle and given border
func (rec *Rectangle) GetLocationToBorder(b *Border) LocationState {
	minX := rec.X
	maxX := rec.X + rec.Width
	minZ := rec.Z
	maxZ := rec.Z + rec.Height
	return b.RectLocation(minX, minZ, maxX, maxZ)
}

// IsCoordInside checks if point is inside the rectangle
// Rectangle vectors are in counter-clockwise order, so when cross product
// between point p and rectangle vectors are all >= 0, point p is inside rectangle
// Note: collinear cases are also considered as inside, so we need to include =0
func (rec *Rectangle) IsCoordInside(p Coord) bool {
	pts := rec.GetVerticeCoords()

	pa := NewVector(p, pts[0])
	pb := NewVector(p, pts[1])
	pc := NewVector(p, pts[2])

	b1 := pa.Cross(&pb) >= 0
	b2 := pb.Cross(&pc) >= 0
	if b1 != b2 {
		return false
	}

	pd := NewVector(p, pts[3])
	b3 := pc.Cross(&pd) >= 0
	if b2 != b3 {
		return false
	}

	b4 := pd.Cross(&pa) >= 0
	return b3 == b4
}
