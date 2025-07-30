package geo

import (
	"log"
	"math"
)

// Convex 凸多边形
type Convex struct {
	Index          int32       // convex序号，唯一标示
	Vertices       []Vertice   // 三角形包含三个顶点
	MergeTriangles []*Triangle // 合成的三角形
	EdgeIDs        []int32
	WtCoord        Coord
}

// NewConvex 转化三角形为凸多边形
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

// ToRect 获取四边形矩形边界
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

// MergeTriangle 合成新的凸多边形
func (c *Convex) MergeTriangle(p1, p2 Vertice, p3 []Vertice) bool {
	for index, v := range c.Vertices {
		if v.Index == p1.Index || v.Index == p2.Index {
			newVertices := make([]Vertice, len(c.Vertices))
			copy(newVertices, c.Vertices)
			// 将新的点插入到正确的位置
			if index != 0 || c.Vertices[index+1].Index == p2.Index || c.Vertices[index+1].Index == p1.Index {
				// 尝试插入两次,因为不能确定点的排列方式
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
					// 反序排列
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
					// 反序排列
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

// GetVectors 获取凸多边形向量组，逆时针排列,
// 凸多边形相邻边的差积大于0则是逆时针排序
func (c *Convex) GetVectors() []Vector {
	c.CounterClockWiseSort()
	vecs := make([]Vector, 0, len(c.Vertices))
	for _, v := range c.Vertices {
		vecs = append(vecs, NewVectorByCoord(v.Coord))
	}
	return vecs
}

// TriangleHasCoord 返回点所在的三角形
func (c *Convex) TriangleHasCoord(p Coord) int32 {
	for _, t := range c.MergeTriangles {
		if t.IsCoordInside(p) {
			return t.GetIndex()
		}
	}
	return -1
}

// IsCoordInside1 判断点在不在凸多边形内,射线法
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
		// i//j//aixs_x
		if vj.Coord.Z == vi.Coord.Z {
			if y == vi.Coord.Z && xmin <= x && x <= xmax {
				return true
			}
			continue
		}

		xt := (vj.Coord.X-vi.Coord.X)*(y-vi.Coord.Z)/(vj.Coord.Z-vi.Coord.Z) + vi.Coord.X
		if xt == x && ymin <= y && y <= ymax {
			// on edge [vj,vi]
			return true
		}
		if x < xt && ymin <= y && y < ymax {
			isIn = !isIn
		}

	}
	return isIn
}

// IsCoordInside2 Check if point is inside a convex
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
	// step 1: the point should between two vector point 0 to point 1 and point 0 to point n-1
	if (cp1 > 0) == (cp2 > 0) {
		return false
	}

	s := 1
	e := numOfVertice - 1
	// step 2: using binary search to determin point is between which two point
	for e != s+1 {
		m := (s + e) / 2
		if (CrossProduct(target, c.Vertices[m], c.Vertices[0]) > 0) == isCounterClockwise { // target在m顺时针方向
			e = m
		} else { // target在m逆时针方向
			s = m
		}
	}
	// step 3: the polygon will divided to a triangle finally,
	// check the point position of the finally vector,
	// the direction should be the same as the point position of point 0 to point 1
	vec3 := NewVector(c.Vertices[s].Coord, c.Vertices[e].Coord)
	vecS2Coord := NewVector(c.Vertices[s].Coord, p)
	return vec3.Cross(&vecS2Coord) > 0 == isCounterClockwise
}

// IsCoordInside Check if point is inside a convex
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

// CheckConvex 测试用查看这个凸多边形的合成正不正确
func (c *Convex) CheckConvex() bool {
	vers := make(map[int32]bool)
	for _, ver := range c.Vertices {
		vers[ver.Index] = false
	}
	if !IsConvex(c.Vertices) {
		log.Printf("[ERROR] convex error not a convex d\n")
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

// CounterClockWiseSort 逆时针排列点
func (c *Convex) CounterClockWiseSort() {
	if CrossProduct(c.Vertices[0], c.Vertices[1], c.Vertices[2]) < 0 {
		// 顺时针排序,反序排列
		for i, j := 0, len(c.Vertices)-1; i < j; i, j = i+1, j-1 {
			c.Vertices[i], c.Vertices[j] = c.Vertices[j], c.Vertices[i]
		}
	}
}

// GetIndex 获取序号
func (c *Convex) GetIndex() int32 {
	return c.Index
}

// GetNeighborPoints 获取两个凸多边形的邻接点
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
	// 返回多边形的点,前两个为临接边上的两点, 只有相邻点为2时是相邻凸多边形
	if len(pointsIndexs) == 2 {
		points := make([]Vertice, numOfVecs)
		first := pointsIndexs[0]
		second := pointsIndexs[1]
		if (first+1)%numOfVecs != second { // 如果first的下一个点不是second则交换, second的值不需要了直接给first赋值即可
			first = second
		}
		for i := 0; i < numOfVecs; i++ {
			points[i] = vertices[(first+i)%numOfVecs]
		}
		return points
	}
	return nil
}

// GetVertices 获取点列表
func (c *Convex) GetVertices() []Vertice {
	return c.Vertices
}

// GetCenterCoord 获取中心点
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

// GetCenterCoord1 获取重心
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

// GetEdgeIDs 返回多边形边的序号列表
func (c *Convex) GetEdgeIDs() []int32 {
	return c.EdgeIDs
}

// GetEdgeMidCoords 获取边中点
func (c *Convex) GetEdgeMidCoords() []Coord {
	nums := len(c.Vertices)
	coords := make([]Coord, 0, nums)
	for i, v := range c.Vertices {
		coords = append(coords, CalMidCoord(v.Coord, c.Vertices[(i+1)%nums].Coord))
	}
	return coords
}
