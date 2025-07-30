package geo

import "github.com/busyster996/geo/util"

// Segment 线段
type Segment struct {
	A, B Coord
}

// NewSegment 新建
func NewSegment(a, b Coord) Segment {
	return Segment{
		A: a,
		B: b,
	}
}

// ToVector 获取向量
func (s *Segment) ToVector() Vector {
	return NewVector(s.A, s.B)
}

// CalCoordDst 计算点到线段的距离
func (s *Segment) CalCoordDst(coord Coord) float64 {
	// 点到线段端点的距离
	a := CalDstCoordToCoord(coord, s.A)
	b := CalDstCoordToCoord(coord, s.B)
	dst := min(a, b)

	// 判断点与线段的垂线是否与线段相交
	ab := NewVector(s.A, s.B)
	ap := NewVector(s.A, coord)
	if lab := ab.Length(); util.AC.Greater(ap.Dot(&ab)/lab, lab) {
		return dst
	}

	// 点到直线的距离
	vec := s.ToVector()
	c := vec.CalCoordDst(s.A, coord)

	return min(dst, c)
}

// Pan 线段平移，左手坐标系，按照叉积为正的方向平行移动一定距离
func (s *Segment) Pan(dst int32, positive bool) Segment {
	v := s.ToVector()
	// 法向量，叉积为正的方向
	normalV := NewVectorByCoord(Coord{X: -v.Z, Z: v.X})
	if !positive {
		normalV = NewVectorByCoord(Coord{X: v.Z, Z: -v.X})
	}
	ratio := float64(dst) / v.Length()
	newV := normalV.Trunc(ratio)

	return NewSegment(newV.ToCoord(s.A), newV.ToCoord(s.B))
}

// CrossCircle 判断线段是否和圆相交，相交则返回第一个交点
func (s *Segment) CrossCircle(circle Circle) (Coord, bool) {
	return GetLineCrossCircle(s.A, s.B, circle.Center, circle.Radius)
}

// IsRectCross 排斥判断
func IsRectCross(p0, p1, q0, q1 Coord) bool {
	ret := min(p0.X, p1.X) <= max(q0.X, q1.X) &&
		min(q0.X, q1.X) <= max(p0.X, p1.X) &&
		min(p0.Z, p1.Z) <= max(q0.Z, q1.Z) &&
		min(q0.Z, q1.Z) <= max(p0.Z, p1.Z)
	return ret
}

// IsLineSegmentCross 跨立判断
func IsLineSegmentCross(p0, p1, q0, q1 Coord) bool {
	// q0q1 X q0p0
	b1 := cross(q1, p0, q0)
	// q0q1 X q0p1
	b2 := cross(q1, p1, q0)

	// 叉积等于0，说明其中一点与另外的线段共线
	if b1 == 0 || b2 == 0 {
		return true
	}

	// p0p1 X p0q0
	a1 := cross(p1, q0, p0)
	// p0p1 X p0q1
	a2 := cross(p1, q1, p0)

	if a1 == 0 || a2 == 0 {
		return true
	}

	return ((b1 < 0) != (b2 < 0)) && ((a1 < 0) != (a2 < 0))
}

// GetCrossCoord 求线段P1P2与Q1Q2的交点。
// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect/565282#
func GetCrossCoord(p0, p1, q0, q1 Coord) (Coord, bool) {
	v1 := NewVector(p0, p1)
	v2 := NewVector(q0, q1)
	// 共线
	if v1.Cross(&v2) == 0 {
		return Coord{}, false
	}

	if IsRectCross(p0, p1, q0, q1) {
		if IsLineSegmentCross(p0, p1, q0, q1) {
			// 求交点
			s1X := float64(p1.X - p0.X)
			s1Z := float64(p1.Z - p0.Z)
			s2X := float64(q1.X - q0.X)
			s2Z := float64(q1.Z - q0.Z)

			t := (s2X*float64(p0.Z-q0.Z) - s2Z*float64(p0.X-q0.X)) / (-s2X*s1Z + s1X*s2Z)

			return Coord{X: p0.X + int32(t*s1X), Z: p0.Z + int32(t*s1Z)}, true
		}
	}
	return Coord{}, false
}
