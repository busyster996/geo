package geo

import "math/rand"

// Rectangle 矩形对象
type Rectangle struct {
	Coord  // 左下角点
	Width  int32
	Height int32
}

// NewRectangle 创建矩形
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

// RandCoord 随机一个坐标
func (rec *Rectangle) RandCoord() Coord {
	return Coord{
		X: rec.X + rand.Int31n(rec.Width),
		Z: rec.Z + rand.Int31n(rec.Height),
	}
}

// GetVerticeCoords 获取4个顶点坐标,边界点按照逆时针排列
func (rec *Rectangle) GetVerticeCoords() [4]Coord {
	var p [4]Coord
	p[0] = Coord{X: rec.Coord.X, Z: rec.Coord.Z}
	p[1] = Coord{X: rec.Coord.X + rec.Width, Z: rec.Coord.Z}
	p[2] = Coord{X: rec.Coord.X + rec.Width, Z: rec.Coord.Z + rec.Height}
	p[3] = Coord{X: rec.Coord.X, Z: rec.Coord.Z + rec.Height}
	return p
}

// GetVectors 获取4条线段，按照逆时针排列
func (rec *Rectangle) GetVectors() [4]Vector {
	coords := rec.GetVerticeCoords()
	return [4]Vector{
		NewVectorByCoord(coords[0]),
		NewVectorByCoord(coords[1]),
		NewVectorByCoord(coords[2]),
		NewVectorByCoord(coords[3]),
	}
}

// GetLocationToBorder 获得矩形和给定边界的位置关系
func (rec *Rectangle) GetLocationToBorder(b *Border) LocationState {
	minX := rec.X
	maxX := rec.X + rec.Width
	minZ := rec.Z
	maxZ := rec.Z + rec.Height
	return b.RectLocation(minX, minZ, maxX, maxZ)
}

// IsCoordInside 点是否在矩形内
// 矩形的向量是逆时针排列，所以当点p和矩形的向量之间求叉积，且全部大于0时，表示在同一侧，则点p在矩形内
// 需要注意的是：共线的情况也认为在同一侧，所以需要加上=0的判断
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
