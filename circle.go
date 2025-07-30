package geo

import (
	"math"

	"github.com/busyster996/geo/util"
)

// Circle 圆
type Circle struct {
	Center Coord
	Radius int32
}

// NewCirCle 新建
func NewCirCle(center Coord, radius int32) Circle {
	return Circle{
		Center: center,
		Radius: radius,
	}
}

// GetLocationToBorder 获取和Border的相对位置
func (c *Circle) GetLocationToBorder(b *Border) LocationState {
	minX, minZ, maxX, maxZ := c.ToRect()
	return b.RectLocation(minX, minZ, maxX, maxZ)
}

// ToRect 泛化为矩形
func (c *Circle) ToRect() (minX, minZ, maxX, maxZ int32) {
	minX = c.Center.X - c.Radius
	minZ = c.Center.Z - c.Radius
	maxX = c.Center.X + c.Radius
	maxZ = c.Center.Z + c.Radius
	return
}

// GetIntersectCoord 求圆外一点到中心的线段和圆的交点
func (c *Circle) GetIntersectCoord(p Coord) Coord {
	// 中心点和p重合时，直接返回中心点
	if c.Center == p {
		return c.Center
	}
	vec := NewVector(c.Center, p)
	length := vec.Length()
	ratio := float64(c.Radius) / length
	vec = vec.Trunc(ratio)
	return vec.ToCoord(c.Center)
}

// IsIntersect 线段与圆是否相交
func (c *Circle) IsIntersect(s *Segment) bool {
	_, ok := c.GetLineCross(s)
	return ok
}

// GetLineCross 线段与圆的交点
// 如果有交点，返回第一个交点，如果没有，返回nil
func (c *Circle) GetLineCross(s *Segment) (Coord, bool) {
	var coord1 *Coord
	var coord2 *Coord
	fDis := CalDstCoordToCoord(s.A, s.B)

	dx := float64(s.B.X-s.A.X) / fDis
	dz := float64(s.B.Z-s.A.Z) / fDis

	ex := float64(c.Center.X - s.A.X)
	ez := float64(c.Center.Z - s.A.Z)

	a := ex*dx + ez*dz
	a2 := a * a
	e2 := ex*ex + ez*ez
	r2 := float64(c.Radius * c.Radius)
	if util.AC.Smaller(r2-e2+a2, 0) {
		return Coord{}, false
	}
	f := math.Sqrt(r2 - e2 + a2)
	t := a - f
	if t > -util.Epsilon && (t-fDis) < util.Epsilon {
		coord1 = &Coord{
			X: s.A.X + int32(t*dx),
			Z: s.A.Z + int32(t*dz),
		}
	}
	t = a + f
	if t > -util.Epsilon && (t-fDis) < util.Epsilon {
		coord2 = &Coord{
			X: s.A.X + int32(t*dx),
			Z: s.A.Z + int32(t*dz),
		}
	}

	if coord1 == nil {
		coord1 = coord2
	}
	if coord1 == nil {
		return Coord{}, false
	}
	return *coord1, true
}

// IsInterPolygon 判断圆和多边形是否相交，相交返回true，不相交返回false
// 圆在多边形内部，相交；多边形在圆内部，相交；圆和多边形只有部分相交，相交。圆和多边形完全分离，不相交
// 此处的相交判断可以类比于碰撞判断
// circle: 圆
// vectors: 多边形向量数组，逆时针
// 参考链接：https://bitlush.com/blog/circle-vs-polygon-collision-detection-in-c-sharp
func (c *Circle) IsInterPolygon(vectors []Vector) bool {
	radiusSquared := float64(c.Radius * c.Radius)

	vertex := vectors[len(vectors)-1]
	center := NewVectorByCoord(c.Center)

	nearestDistance := math.MaxFloat64
	nearestIsInside := false
	nearestVertex := -1
	lastIsInside := false

	for i := 0; i < len(vectors); i++ {
		nextVertex := vectors[i]
		axis := center.Minus(&vertex)
		distance := axis.LengthSquared() - radiusSquared

		if util.AC.SmallerOrEqual(distance, 0) {
			return true
		}

		isInside := false
		edge := nextVertex.Minus(&vertex)
		edgeLengthSquared := edge.LengthSquared()
		if !util.AC.Equal(edgeLengthSquared, 0) {
			dot := edge.Dot(&axis)
			if util.AC.GreaterOrEqual(dot, 0) && util.AC.SmallerOrEqual(dot, edgeLengthSquared) {
				projection := edge.Trunc(dot / edgeLengthSquared)
				projection = vertex.Add(&projection)

				axis = projection.Minus(&center)
				if util.AC.SmallerOrEqual(axis.LengthSquared(), radiusSquared) {
					return true
				}

				if !isInsideEdge(&edge, &axis) {
					return false
				}

				isInside = true
			}
		}

		if util.AC.Smaller(distance, nearestDistance) {
			nearestDistance = distance
			nearestIsInside = isInside || lastIsInside
			nearestVertex = i
		}

		vertex = nextVertex
		lastIsInside = isInside
	}

	if nearestVertex == 0 {
		return nearestIsInside || lastIsInside
	}
	return nearestIsInside
}

func isInsideEdge(edge, axis *Vector) bool {
	switch {
	case edge.X > 0 && axis.Z > 0:
		return false
	case edge.X < 0 && axis.Z < 0:
		return false
	case edge.X == 0 && edge.Z > 0 && axis.X < 0:
		return false
	case edge.X == 0 && edge.Z <= 0 && axis.X > 0:
		return false
	}
	return true
}
