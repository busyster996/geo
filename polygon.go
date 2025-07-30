package geo

// Polygon defines the interface for polygon operations
type Polygon interface {
	IsCoordInside(p Coord) bool             // Check if point is inside polygon
	GetVectors() []Vector                   // Get vector representation of polygon
	ToRect() (minX, minZ, maxX, maxZ int32) // Get bounding rectangle
	GetIndex() int32                        // Get polygon index
	GetEdgeIDs() []int32                    // Get edge identifiers
	GetEdgeMidCoords() []Coord              // Get midpoints of edges
	GetVertices() []Vertice                 // Get vertices of polygon
}

// CrossProduct calculates the cross product of three vertices
// Returns 1 for counter-clockwise, -1 for clockwise, 0 for collinear
func CrossProduct(p1, p2, p3 Vertice) int32 {
	ax := p2.Coord.X - p1.Coord.X
	ay := p2.Coord.Z - p1.Coord.Z
	bx := p3.Coord.X - p2.Coord.X
	by := p3.Coord.Z - p2.Coord.Z
	cp := ax*by - ay*bx
	if cp > 0 {
		return 1
	} else if cp < 0 {
		return -1
	} else {
		return 0
	}
}

// IsConvex checks if the given vertices form a convex polygon
func IsConvex(vertices []Vertice) bool {
	numPoints := len(vertices)
	negativeFlag := false
	positiveFlag := false
	for i := 0; i < numPoints; i++ {
		curvec := NewVector(vertices[i].Coord, vertices[(i+1)%numPoints].Coord)
		vec2next := NewVector(vertices[i].Coord, vertices[(i+2)%numPoints].Coord)
		if curvec.Cross(&vec2next) > 0 {
			positiveFlag = true
		} else {
			negativeFlag = true
		}
	}
	return positiveFlag != negativeFlag
}
