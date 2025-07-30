package geo

import (
	"log"
	"math"
)

// Convex represents a convex polygon
type Convex struct {
	Index          int32       // Unique identifier for the convex polygon
	Vertices       []Vertice   // Vertices of the polygon (triangles contain three vertices)
	MergeTriangles []*Triangle // Triangles that compose this convex polygon
	EdgeIDs        []int32     // Edge identifiers
	WtCoord        Coord       // Weight coordinate (center of mass)
}

// NewConvex converts a triangle to a convex polygon
func NewConvex(t *Triangle, id int32) *Convex {
	return &Convex{
		Index: id,
		Vertices: []Vertice{
			t.Vertices[0],
			t.Vertices[1],
			t.Vertices[2],
		},
		MergeTriangles: []*Triangle{t},
		EdgeIDs:        t.EdgeIDs,
	}
}

// ToRect returns the bounding rectangle of the convex polygon
// Returns minimum and maximum X,Z coordinates
func (c *Convex) ToRect() (minX, minZ, maxX, maxZ int32) {
	minX = int32(math.MaxInt32)
	minZ = int32(math.MaxInt32)
	for _, v := range c.Vertices {
		minX = min(v.Coord.X, minX)
		minZ = min(v.Coord.Z, minZ)
		maxX = max(v.Coord.X, maxX)
		maxZ = max(v.Coord.Z, maxZ)
	}
	return minX, minZ, maxX, maxZ
}

// MergeTriangle merges new triangle to form a new convex polygon
// p1, p2: shared vertices between triangles
// p3: additional vertices from the new triangle
// Returns true if merge is successful and results in a valid convex polygon
func (c *Convex) MergeTriangle(p1, p2 Vertice, p3 []Vertice) bool {
	for index, v := range c.Vertices {
		if v.Index == p1.Index || v.Index == p2.Index {
			newVertices := make([]Vertice, len(c.Vertices))
			copy(newVertices, c.Vertices)
			// Insert new points at correct position
			if index != 0 || c.Vertices[index+1].Index == p2.Index || c.Vertices[index+1].Index == p1.Index {
				// Try inserting twice since we can't determine the point arrangement
				for i := 0; i < 2; i++ {
					newVertices = make([]Vertice, len(c.Vertices)+len(p3))
					insertIndex := 0
					for insertIndex <= index {
						newVertices[insertIndex] = c.Vertices[insertIndex]
						insertIndex++
					}
					for _, p := range p3 {
						newVertices[insertIndex] = p
						insertIndex++
					}
					for _, p := range c.Vertices[index+1:] {
						newVertices[insertIndex] = p
						insertIndex++
					}
					if IsConvex(newVertices) {
						c.Vertices = newVertices
						return true
					}
					// Reverse order
					for i, j := 0, len(p3)-1; i < j; i, j = i+1, j-1 {
						p3[i], p3[j] = p3[j], p3[i]
					}
				}
			} else {
				for i := 0; i < 2; i++ {
					newVertices = make([]Vertice, len(c.Vertices))
					copy(newVertices, c.Vertices)
					newVertices = append(newVertices, p3...)
					if IsConvex(newVertices) {
						c.Vertices = newVertices
						return true
					}
					// Reverse order
					for i, j := 0, len(p3)-1; i < j; i, j = i+1, j-1 {
						p3[i], p3[j] = p3[j], p3[i]
					}
				}
			}
			return false
		}
	}
	return false
}

// GetVectors returns vector array of convex polygon in counter-clockwise order
// Adjacent edges of convex polygon have positive cross product in counter-clockwise order
func (c *Convex) GetVectors() []Vector {
	c.CounterClockWiseSort()
	vecs := make([]Vector, 0, len(c.Vertices))
	for _, v := range c.Vertices {
		vecs = append(vecs, NewVectorByCoord(v.Coord))
	}
	return vecs
}

// TriangleHasCoord returns the triangle index that contains the given point
func (c *Convex) TriangleHasCoord(p Coord) int32 {
	for _, t := range c.MergeTriangles {
		if t.IsCoordInside(p) {
			return t.GetIndex()
		}
	}
	return -1
}

// IsCoordInside1 checks if point is inside convex polygon using ray casting algorithm
func (c *Convex) IsCoordInside1(p Coord) bool {
	if len(c.MergeTriangles) == 1 {
		return c.MergeTriangles[0].IsCoordInside(p)
	}
	c.CounterClockWiseSort()
	x := p.X
	y := p.Z
	sz := len(c.Vertices)
	isIn := false

	for i := 0; i < sz; i++ {
		j := i - 1
		if i == 0 {
			j = sz - 1
		}
		vi := c.Vertices[i]
		vj := c.Vertices[j]

		xmin := vi.Coord.X
		xmax := vj.Coord.X
		if xmin > xmax {
			t := xmin
			xmin = xmax
			xmax = t
		}
		ymin := vi.Coord.Z
		ymax := vj.Coord.Z
		if ymin > ymax {
			t := ymin
			ymin = ymax
			ymax = t
		}
		// Check if point is on horizontal edge
		if vj.Coord.Z == vi.Coord.Z {
			if y == vi.Coord.Z && xmin <= x && x <= xmax {
				return true
			}
			continue
		}

		xt := (vj.Coord.X-vi.Coord.X)*(y-vi.Coord.Z)/(vj.Coord.Z-vi.Coord.Z) + vi.Coord.X
		if xt == x && ymin <= y && y <= ymax {
			// Point is on edge [vj,vi]
			return true
		}
		if x < xt && ymin <= y && y < ymax {
			isIn = !isIn
		}

	}
	return isIn
}

// IsCoordInside2 checks if point is inside convex polygon using binary search
func (c *Convex) IsCoordInside2(p Coord) bool {
	numOfVertice := len(c.Vertices)
	target := Vertice{Coord: p}
	vec1 := NewVector(c.Vertices[0].Coord, c.Vertices[1].Coord)
	vec2 := NewVector(c.Vertices[0].Coord, c.Vertices[numOfVertice-1].Coord)
	vec2Coord := NewVector(c.Vertices[0].Coord, p)
	vec2CoordLen := vec2Coord.Length()
	cp1 := vec1.Cross(&vec2Coord)
	cp2 := vec2.Cross(&vec2Coord)
	isCounterClockwise := cp1 > 0
	if cp1 == 0 && vec2CoordLen <= vec1.Length() || cp2 == 0 && vec2CoordLen <= vec2.Length() {
		return true
	}
	// Step 1: the point should be between two vectors from point 0 to point 1 and point 0 to point n-1
	if (cp1 > 0) == (cp2 > 0) {
		return false
	}

	s := 1
	e := numOfVertice - 1
	// Step 2: use binary search to determine which two points the target point is between
	for e != s+1 {
		m := (s + e) / 2
		if (CrossProduct(target, c.Vertices[m], c.Vertices[0]) > 0) == isCounterClockwise { // target is clockwise to m
			e = m
		} else { // target is counter-clockwise to m
			s = m
		}
	}
	// Step 3: the polygon is divided into a triangle finally,
	// check the point position of the final vector,
	// the direction should be the same as the point position from point 0 to point 1
	vec3 := NewVector(c.Vertices[s].Coord, c.Vertices[e].Coord)
	vecS2Coord := NewVector(c.Vertices[s].Coord, p)
	return vec3.Cross(&vecS2Coord) > 0 == isCounterClockwise
}

// IsCoordInside checks if point is inside convex polygon
func (c *Convex) IsCoordInside(p Coord) bool {
	numOfVertice := len(c.Vertices)
	vec1 := NewVector(c.Vertices[0].Coord, c.Vertices[1].Coord)
	vec2 := NewVector(c.Vertices[0].Coord, p)
	isCounterClockwise := vec1.Cross(&vec2) > 0
	for i := 1; i < numOfVertice; i++ {
		vec1 := NewVector(c.Vertices[i].Coord, c.Vertices[(i+1)%numOfVertice].Coord)
		vec2 := NewVector(c.Vertices[i].Coord, p)
		if (vec1.Cross(&vec2) > 0) != isCounterClockwise {
			return false
		}
	}
	return true
}

// CheckConvex validates if this convex polygon is correctly formed (for testing)
func (c *Convex) CheckConvex() bool {
	vers := make(map[int32]bool)
	for _, ver := range c.Vertices {
		vers[ver.Index] = false
	}
	if !IsConvex(c.Vertices) {
		log.Printf("[ERROR] convex error not a convex polygon\n")
		return false
	}
	count1 := len(vers)
	for _, triangle := range c.MergeTriangles {
		for _, ver := range triangle.GetVertices() {
			vers[ver.Index] = true
		}
	}
	count2 := len(vers)
	for key, value := range vers {
		if !value {
			log.Printf("[ERROR] convex error has not exist value %d\n", key)
			return false
		}
	}

	if count1 != count2 {
		log.Printf("[ERROR] convex error has not merged\nvertice convex vertices count: %d merge vertices count: %d\n", count1, count2)
		return false
	}
	return true
}

// CounterClockWiseSort sorts vertices in counter-clockwise order
func (c *Convex) CounterClockWiseSort() {
	if CrossProduct(c.Vertices[0], c.Vertices[1], c.Vertices[2]) < 0 {
		// Clockwise order, reverse the array
		for i, j := 0, len(c.Vertices)-1; i < j; i, j = i+1, j-1 {
			c.Vertices[i], c.Vertices[j] = c.Vertices[j], c.Vertices[i]
		}
	}
}

// GetIndex returns the polygon index
func (c *Convex) GetIndex() int32 {
	return c.Index
}

// GetNeighborPoints returns neighboring points between two convex polygons
func (c *Convex) GetNeighborPoints(t2 Polygon) []Vertice {
	vertices := t2.GetVertices()
	numOfVecs := len(vertices)
	pointsIndexs := make([]int, 0, numOfVecs)

	for _, v := range c.Vertices {
		for j := range vertices {
			if v.Index == vertices[j].Index {
				pointsIndexs = append(pointsIndexs, j)
			}
		}
	}
	// Return polygon points, first two are points on adjacent edge
	// Only when adjacent points equal 2, they are neighboring convex polygons
	if len(pointsIndexs) == 2 {
		points := make([]Vertice, numOfVecs)
		first := pointsIndexs[0]
		second := pointsIndexs[1]
		if (first+1)%numOfVecs != second { // If next point of first is not second, swap them
			first = second
		}
		for i := 0; i < numOfVecs; i++ {
			points[i] = vertices[(first+i)%numOfVecs]
		}
		return points
	}
	return nil
}

// GetVertices returns the list of vertices
func (c *Convex) GetVertices() []Vertice {
	return c.Vertices
}

// GetCenterCoord returns the center point (arithmetic mean of vertices)
func (c *Convex) GetCenterCoord() Coord {
	coordx := int32(0)
	coordz := int32(0)
	length := int32(len(c.Vertices))
	for _, v := range c.Vertices {
		coordx += v.Coord.X
		coordz += v.Coord.Z
	}
	return Coord{coordx / length, coordz / length}
}

// GetCenterCoord1 returns the centroid (center of mass) of the polygon
func (c *Convex) GetCenterCoord1() Coord {
	S := int32(0)
	for i := 0; i < len(c.Vertices)-1; i++ {
		S += c.Vertices[i].Coord.X*c.Vertices[i+1].Coord.Z - c.Vertices[i+1].Coord.X*c.Vertices[i].Coord.Z
	}
	S /= 2
	centerX := int32(0)
	centerY := int32(0)
	for i := 0; i < len(c.Vertices)-1; i++ {
		xi := c.Vertices[i].Coord.X
		yi := c.Vertices[i].Coord.Z
		xni := c.Vertices[i+1].Coord.X
		yni := c.Vertices[i+1].Coord.Z
		centerX += (xi + xni) * (xi*yni - xni*yi)
		centerY += (yi + yni) * (xi*yni - xni*yi)
	}
	centerX /= 6 * S
	centerY /= 6 * S
	return Coord{
		X: centerX,
		Z: centerY,
	}

}

// GetEdgeIDs returns the list of edge indices
func (c *Convex) GetEdgeIDs() []int32 {
	return c.EdgeIDs
}

// GetEdgeMidCoords returns the midpoints of all edges
func (c *Convex) GetEdgeMidCoords() []Coord {
	nums := len(c.Vertices)
	coords := make([]Coord, 0, nums)
	for i, v := range c.Vertices {
		coords = append(coords, CalMidCoord(v.Coord, c.Vertices[(i+1)%nums].Coord))
	}
	return coords
}
