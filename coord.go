package geo

import "math"

// Coord 坐标
type Coord struct {
	X int32 `bson:"x"`
	Z int32 `bson:"z"`
}

// NewCoord 新建
func NewCoord(x, z int32) Coord {
	return Coord{
		X: x,
		Z: z,
	}
}

// IsEqual 判断是否相等
func (c *Coord) IsEqual(target Coord) bool {
	return *c == target
}

// GetLocationToBorder 获取和边界的位置关系
func (c *Coord) GetLocationToBorder(b *Border) LocationState {
	return b.CoordLocation(*c)
}

// CalDstCoordToCoord 计算点和点之间的距离
func CalDstCoordToCoord(coord1, coord2 Coord) float64 {
	dx := int64(coord1.X - coord2.X)
	dz := int64(coord1.Z - coord2.Z)

	squared := dx*dx + dz*dz
	return math.Sqrt(float64(squared))
}

// CalDstCoordToCoordWithoutSqrt 计算点和点之间的距离,不开根号
func CalDstCoordToCoordWithoutSqrt(coord1, coord2 Coord) float64 {
	dx := int64(coord1.X - coord2.X)
	dz := int64(coord1.Z - coord2.Z)

	squared := dx*dx + dz*dz
	return float64(squared)
}
