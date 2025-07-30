package geo

import (
	"math"

	"github.com/busyster996/geo/util"
)

// GetCoordsAround 获得圆圈周围的点
// startCoord: 圆上一点
// centerCoord: 圆心
// endCoord: 圆外点
func GetCoordsAround(startCoord, endCoord, centerCoord Coord) []Coord {
	centerVector := NewVector(centerCoord, startCoord)
	endVector := NewVector(centerCoord, endCoord)

	// 向量夹角
	angle := centerVector.GetAngle(&endVector)
	radius := CalDstCoordToCoord(startCoord, centerCoord)
	// 切线角度
	cutAngle := GetCutOffCoordAngle(endCoord, centerCoord, radius)
	angle -= cutAngle
	// 叉积确定方向
	if centerVector.Cross(&endVector) > 0 {
		angle = -angle
	}
	coords := GetArcCoords(startCoord, centerCoord, angle)

	return coords
}

// GetCoordsAround2 获取圆周围均匀分布的点
// circleCoord: 圆周上的点
// centerCoord: 圆心
// n: 均匀分配点的个数
func GetCoordsAround2(circleCoord, centerCoord Coord, n int) []Coord {
	return getCoordsAround(circleCoord, centerCoord, n, 2*math.Pi)
}

// getCoordsAround 获取圆上的startCoord出发, angle角度内的点, 包括startCoord
// n: 采样数
func getCoordsAround(startCoord, centerCoord Coord, n int, angle float64) []Coord {
	if n < 2 {
		return nil
	}

	ret := make([]Coord, n)
	ret[0] = startCoord

	vec := NewVector(centerCoord, startCoord)
	ave := angle / float64(n-1)
	angle = 0
	for i := 1; i < n; i++ {
		angle += ave
		v := vec.Rotate(angle)
		ret[i] = v.ToCoord(centerCoord)
	}

	return ret
}

// GetArcCoords 获取圆弧上的点, 圆弧从startCoord点出发, angle弧度
// angle: 角度, >0 为顺时针方向, <0 为逆时针方向
func GetArcCoords(startCoord, centerCoord Coord, angle float64) []Coord {
	// pi 34个点
	n := max(int32(util.Abs(angle)*10), 2)
	return getCoordsAround(startCoord, centerCoord, int(n), -angle)
}

// GetSpiralCoords 按照螺线获取路径点
// startCoord: 起始坐标
// centerCoord: 中心点
// angle: 角度, >0 为顺时针方向, <0 为逆时针方向
// delta: 最后一个点与起点 相对于中心的距离偏差
func GetSpiralCoords(startCoord, centerCoord Coord, angle, delta float64) []Coord {
	n := int32(util.Abs(angle) * 10)
	if n < 2 {
		dst := CalDstCoordToCoord(startCoord, centerCoord)
		vec := NewVector(centerCoord, startCoord)
		vec = vec.Trunc((delta + dst) / dst)
		return []Coord{startCoord, vec.ToCoord(centerCoord)}
	}

	ret := make([]Coord, n)
	ret[0] = startCoord

	aveAngle := angle / float64(n)
	aveDelta := (delta) / float64(n)
	angle = 0
	delta = 0

	vec := NewVector(centerCoord, startCoord)
	radius := vec.Length()
	for i := int32(1); i < n; i++ {
		// 将原始向量旋转
		angle += aveAngle
		v := vec.Rotate(angle)
		// 增加将向量长度
		delta += aveDelta
		v = v.Trunc((radius + delta) / radius)

		ret[i] = v.ToCoord(centerCoord)
	}

	return ret
}

// GetIntersectCoord 求圆外一点到中心的线段和圆的交点
func GetIntersectCoord(centerP, endP Coord, radius int32) Coord {
	c := NewCirCle(centerP, radius)
	return c.GetIntersectCoord(endP)
}

// GetLineCrossCircle 线段与圆的交点
// startP: 起点坐标
// endP: 终点坐标
// centerP: 终点坐标
// radius: 圆半径
// 如果有交点，返回第一个交点，如果没有，返回nil
func GetLineCrossCircle(startP, endP, centerP Coord, radius int32) (Coord, bool) {
	c := NewCirCle(centerP, radius)
	seg := NewSegment(startP, endP)
	return c.GetLineCross(&seg)
}

// GetCutOffCoordAngle 求圆外一点（endCoord）的到圆(center)的切点的角度
func GetCutOffCoordAngle(endCoord, centerCoord Coord, radius float64) float64 {
	// 令终点为e, 圆心为c, 切点为 p
	// ep ⊥ cp
	// cos∠ecp = ep / cp
	dst := CalDstCoordToCoord(endCoord, centerCoord)
	a := math.Acos(radius / dst)

	if math.IsNaN(a) {
		return 0.
	}

	return a
}

// GetBresenhamCoord https://zh.wikipedia.org/wiki/布雷森漢姆直線演算法
func GetBresenhamCoord(p1, p2 Coord) []Coord {
	xstart := p1.X
	zstart := p1.Z
	xend := p2.X
	zend := p2.Z

	var steep = false
	var swapped = false

	if util.Abs(zend-zstart) > util.Abs(xend-xstart) {
		steep = true
	}

	if steep {
		xstart, zstart = zstart, xstart
		xend, zend = zend, xend
	}
	if xstart > xend {
		xstart, zstart = zstart, xstart
		xend, zend = zend, xend
		swapped = true
	}
	var deltax = xend - xstart
	var deltay = util.Abs(zend - zstart)
	var err = deltax / 2
	var zstep int32
	z := zstart
	if zstart < zend {
		zstep = 1
	} else {
		zstep = -1
	}
	tmpCoordList := make([]Coord, 0, xend-xstart+1)
	for x := xstart; x <= xend; x++ {
		var c Coord
		if steep {
			c = Coord{
				X: z, Z: x,
			}
		} else {
			c = Coord{
				X: x, Z: z,
			}
		}
		tmpCoordList = append(tmpCoordList, c)
		err -= deltay
		if err < 0 {
			z += zstep
			err += deltax
		}
	}

	if swapped {
		reverse(tmpCoordList)
	}
	return tmpCoordList
}

// reverse 将slice反转
func reverse(slice []Coord) {
	for i, j := 0, len(slice)-1; j > i; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// CalMidCoord 计算中心点
func CalMidCoord(p1, p2 Coord) Coord {
	return Coord{
		X: (p1.X + p2.X) / 2,
		Z: (p1.Z + p2.Z) / 2,
	}
}

// GetCrossRect 计算线段经过的范围的矩形坐标算法，范围被划分为n个小矩形
// p0, p1: 线段的起始点坐标
// rWidth, rHeight: 小矩形的宽
// width, height: 范围的宽和高
func GetCrossRect(p0, p1 Coord, rWidth, rHeight, width, height int32) map[Coord]bool {
	coordSet := map[Coord]bool{}
	// 将起点和终点所在的矩形加入
	startCoords := []Coord{p0, p1}
	for i := 0; i < len(startCoords); i++ {
		c := Coord{X: (startCoords[i].X / rHeight) * rHeight, Z: (startCoords[i].Z / rHeight) * rHeight}
		_, exist := coordSet[c]
		if !exist {
			coordSet[c] = true
		}
	}
	// 由小到大
	zStart := p0
	zEnd := p1
	if p1.Z < p0.Z {
		zStart = p1
		zEnd = p0
	}
	startZ := zStart.Z / rHeight
	endZ := zEnd.Z / rHeight
	// 遍历横轴
	for i := startZ; i <= endZ; i++ {
		q0 := &Coord{X: 0, Z: i * rHeight}
		q1 := &Coord{X: width, Z: i * rHeight}
		// 计算交点
		coord, ok := GetCrossCoord(p0, p1, *q0, *q1)
		if ok {
			// 校正误差
			coord.Z = i * rHeight
			upCoord := Coord{X: (coord.X / rHeight) * rHeight, Z: (coord.Z / rHeight) * rHeight}
			downCoord := Coord{X: upCoord.X, Z: (coord.Z/rHeight - 1) * rHeight}
			// 上下矩形必然为经过的矩形，如果交点恰好在矩形十字交点上，那么会多算一些矩形
			_, exist := coordSet[upCoord]
			if !exist {
				coordSet[upCoord] = true
			}
			_, exist = coordSet[downCoord]
			if !exist {
				coordSet[downCoord] = true
			}
		}
	}

	xStart := p0
	xEnd := p1
	if p1.X < p0.X {
		xStart = p1
		xEnd = p0
	}
	// 遍历纵轴
	startX := xStart.X / rWidth
	endX := xEnd.X / rWidth
	for i := startX; i <= endX; i++ {
		q0 := &Coord{X: i * rWidth, Z: 0}
		q1 := &Coord{X: i * rWidth, Z: height}
		// 求交点
		coord, ok := GetCrossCoord(p0, p1, *q0, *q1)
		if ok {
			coord.X = i * rWidth
			rightCoord := Coord{X: (coord.X / rWidth) * rWidth, Z: (coord.Z / rWidth) * rWidth}
			leftCoord := Coord{X: (coord.X/rWidth - 1) * rWidth, Z: rightCoord.Z}
			// 左右矩形必然为经过的矩形
			_, exist := coordSet[rightCoord]
			if !exist {
				coordSet[rightCoord] = true
			}
			_, exist = coordSet[leftCoord]
			if !exist {
				coordSet[leftCoord] = true
			}
		}
	}

	return coordSet
}

// GetIntersectRect 求两个矩形相交区域
func GetIntersectRect(r0 Rectangle, r1 Rectangle) (Rectangle, bool) {
	rect := Rectangle{}
	rect.X = max(r0.X, r1.X)
	rect.Z = max(r0.Z, r1.Z)
	rect.Width = min(r0.X+r0.Width, r1.X+r1.Width) - rect.X
	rect.Height = min(r0.Z+r0.Height, r1.Z+r1.Height) - rect.Z

	if rect.X >= rect.X+rect.Width || rect.Z >= rect.Z+rect.Height {
		return rect, false
	}
	return rect, true
}
