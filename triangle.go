package geo

import (
	"math"
)

// Triangle 三角形
type Triangle struct {
	Index         int32     // 三角形序号，唯一标示
	Vertices      []Vertice // 三角形包含三个顶点
	EdgeIDs       []int32   // 三角形三边唯一序号，起服是生成
	EdgeKeyString []string  // 三角形三边序号，起服是生成
}

// IsCoordInside 判断点是否在三角形内部
// 三角形的向量是顺时针排列，所以当点p和三角形的向量之间求叉积，且全部小于0时，表示在同一侧，则点p在三角形内
// 需要注意的是：共线的情况也认为在同一侧，所以需要加上=0的判断
func (t *Triangle) IsCoordInside(p Coord) bool {
	pa := NewVector(p, t.Vertices[0].Coord)
	pb := NewVector(p, t.Vertices[1].Coord)
	pc := NewVector(p, t.Vertices[2].Coord)

	b1 := pa.Cross(&pb) >= 0
	b2 := pb.Cross(&pc) >= 0

	if b1 != b2 {
		return false
	}

	b3 := pc.Cross(&pa) >= 0
	return b2 == b3
}

// GetIndex 获取序号
func (t *Triangle) GetIndex() int32 {
	return t.Index
}

// GetLocationToBorder 判断三角形和边界的位置关系
func (t *Triangle) GetLocationToBorder(b *Border) LocationState {
	minX, minZ, maxX, maxZ := t.ToRect()
	return b.RectLocation(minX, minZ, maxX, maxZ)
}

// ToRect 获取三角形矩形边界
func (t *Triangle) ToRect() (minX, minZ, maxX, maxZ int32) {
	minX = int32(math.MaxInt32)
	minZ = int32(math.MaxInt32)
	for _, v := range t.Vertices {
		minX = min(v.Coord.X, minX)
		minZ = min(v.Coord.Z, minZ)
		maxX = max(v.Coord.X, maxX)
		maxZ = max(v.Coord.Z, maxZ)
	}
	return minX, minZ, maxX, maxZ
}

// GetEdgeIDs 返回三角形边的序号列表
func (t *Triangle) GetEdgeIDs() []int32 {
	return t.EdgeIDs
}

// GetEdgeMidCoords 返回三角形边的中点列表
func (t *Triangle) GetEdgeMidCoords() []Coord {
	coords := make([]Coord, 3)
	coords[0] = CalMidCoord(t.Vertices[0].Coord, t.Vertices[1].Coord)
	coords[1] = CalMidCoord(t.Vertices[1].Coord, t.Vertices[2].Coord)
	coords[2] = CalMidCoord(t.Vertices[2].Coord, t.Vertices[0].Coord)
	return coords
}

// GetVertices 获取点列表
func (t *Triangle) GetVertices() []Vertice {
	return t.Vertices
}

// GetVectors 获取三角形向量组，逆时针排列
func (t *Triangle) GetVectors() []Vector {
	vecs := make([]Vector, 3)

	vecs[0] = NewVectorByCoord(t.Vertices[0].Coord)
	vecs[1] = NewVectorByCoord(t.Vertices[2].Coord)
	vecs[2] = NewVectorByCoord(t.Vertices[1].Coord)

	return vecs
}

// GetNeighborEdgeNums 获取两个三角形邻接边的数量
func (t *Triangle) GetNeighborEdgeNums(t2 *Triangle) int {
	var cnt int
	for i := range t.Vertices {
		for j := range t2.Vertices {
			if t.Vertices[i].Index == t2.Vertices[j].Index {
				cnt++
			}
		}
	}
	return cnt
}
