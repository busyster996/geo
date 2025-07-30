package geo

// Vertice represents a vertex with unique index and coordinate
type Vertice struct {
	Index int32 // Unique vertex identifier
	Coord Coord // Vertex coordinate
}

// Edge represents an edge connecting two vertices
type Edge struct {
	WtCoord            Coord       // Weight coordinate (center of mass)
	Vertices           [2]Vertice  // Vertex array containing two vertices (0 and 1)
	AdjacenctTriangles []*Triangle // Two adjacent triangles
	Inflects           [2]Vertice  // Inflection point array
	IsAdjacency        bool        // Whether this is an adjacent edge
}

// CalMidCoord calculates the midpoint of the edge
func (e *Edge) CalMidCoord() Coord {
	return CalMidCoord(e.Vertices[0].Coord, e.Vertices[1].Coord)
}

// GenKey generates a unique key for the edge
func (e *Edge) GenKey() int32 {
	return GenEdgeKey(e.Vertices[0].Index, e.Vertices[1].Index)
}

// GenEdgeKey generates a unique key for an edge given two vertex indices
func GenEdgeKey(i, j int32) int32 {
	if i < j {
		return 10000*i + j
	}
	return 10000*j + i
}
