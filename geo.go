package geo

import (
	"math"

	"github.com/busyster996/geo/util"
)

// GetCoordsAround gets points around a circle
// startCoord: a point on the circle
// centerCoord: center of the circle
// endCoord: external point outside the circle
func GetCoordsAround(startCoord, endCoord, centerCoord Coord) []Coord {
	centerVector := NewVector(centerCoord, startCoord)
	endVector := NewVector(centerCoord, endCoord)

	// Vector angle
	angle := centerVector.GetAngle(&endVector)
	radius := CalDstCoordToCoord(startCoord, centerCoord)
	// Tangent angle
	cutAngle := GetCutOffCoordAngle(endCoord, centerCoord, radius)
	angle -= cutAngle
	// Use cross product to determine direction
	if centerVector.Cross(&endVector) > 0 {
		angle = -angle
	}
	coords := GetArcCoords(startCoord, centerCoord, angle)

	return coords
}

// GetCoordsAround2 gets evenly distributed points around a circle
// circleCoord: a point on the circle
// centerCoord: center of the circle
// n: number of evenly distributed points
func GetCoordsAround2(circleCoord, centerCoord Coord, n int) []Coord {
	return getCoordsAround(circleCoord, centerCoord, n, 2*math.Pi)
}

// getCoordsAround gets points on circle starting from startCoord within angle range, including startCoord
// n: number of sample points
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

// GetArcCoords gets points on arc starting from startCoord with given angle in radians
// angle: angle in radians, >0 for clockwise, <0 for counter-clockwise
func GetArcCoords(startCoord, centerCoord Coord, angle float64) []Coord {
	// pi radians with 34 points
	n := max(int32(util.Abs(angle)*10), 2)
	return getCoordsAround(startCoord, centerCoord, int(n), -angle)
}

// GetSpiralCoords generates spiral path points
// startCoord: starting coordinate
// centerCoord: center point
// angle: angle in radians, >0 for clockwise, <0 for counter-clockwise
// delta: distance difference between last point and start point relative to center
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
		// Rotate the original vector
		angle += aveAngle
		v := vec.Rotate(angle)
		// Increase vector length
		delta += aveDelta
		v = v.Trunc((radius + delta) / radius)

		ret[i] = v.ToCoord(centerCoord)
	}

	return ret
}

// GetIntersectCoord calculates intersection point between line from external point to center and circle
func GetIntersectCoord(centerP, endP Coord, radius int32) Coord {
	c := NewCirCle(centerP, radius)
	return c.GetIntersectCoord(endP)
}

// GetLineCrossCircle calculates intersection point between line segment and circle
// startP: start point coordinate
// endP: end point coordinate
// centerP: circle center coordinate
// radius: circle radius
// Returns first intersection point if exists, otherwise returns false
func GetLineCrossCircle(startP, endP, centerP Coord, radius int32) (Coord, bool) {
	c := NewCirCle(centerP, radius)
	seg := NewSegment(startP, endP)
	return c.GetLineCross(&seg)
}

// GetCutOffCoordAngle calculates the angle from external point (endCoord) to tangent point of circle (center)
func GetCutOffCoordAngle(endCoord, centerCoord Coord, radius float64) float64 {
	// Let endpoint be e, center be c, tangent point be p
	// ep ⊥ cp
	// cos∠ecp = radius / distance_ec
	dst := CalDstCoordToCoord(endCoord, centerCoord)
	a := math.Acos(radius / dst)

	if math.IsNaN(a) {
		return 0.
	}

	return a
}

// GetBresenhamCoord implements Bresenham's line algorithm
// Reference: https://zh.wikipedia.org/wiki/布雷森漢姆直線演算法
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

// reverse reverses a slice of coordinates
func reverse(slice []Coord) {
	for i, j := 0, len(slice)-1; j > i; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// CalMidCoord calculates the midpoint between two coordinates
func CalMidCoord(p1, p2 Coord) Coord {
	return Coord{
		X: (p1.X + p2.X) / 2,
		Z: (p1.Z + p2.Z) / 2,
	}
}

// GetCrossRect calculates rectangles that a line segment passes through
// The area is divided into n small rectangles
// p0, p1: line segment start and end points
// rWidth, rHeight: small rectangle dimensions
// width, height: total area dimensions
func GetCrossRect(p0, p1 Coord, rWidth, rHeight, width, height int32) map[Coord]bool {
	coordSet := map[Coord]bool{}
	// Add rectangles containing start and end points
	startCoords := []Coord{p0, p1}
	for i := 0; i < len(startCoords); i++ {
		c := Coord{X: (startCoords[i].X / rHeight) * rHeight, Z: (startCoords[i].Z / rHeight) * rHeight}
		_, exist := coordSet[c]
		if !exist {
			coordSet[c] = true
		}
	}
	// Sort from small to large
	zStart := p0
	zEnd := p1
	if p1.Z < p0.Z {
		zStart = p1
		zEnd = p0
	}
	startZ := zStart.Z / rHeight
	endZ := zEnd.Z / rHeight
	// Traverse horizontal axis
	for i := startZ; i <= endZ; i++ {
		q0 := &Coord{X: 0, Z: i * rHeight}
		q1 := &Coord{X: width, Z: i * rHeight}
		// Calculate intersection point
		coord, ok := GetCrossCoord(p0, p1, *q0, *q1)
		if ok {
			// Correct for errors
			coord.Z = i * rHeight
			upCoord := Coord{X: (coord.X / rHeight) * rHeight, Z: (coord.Z / rHeight) * rHeight}
			downCoord := Coord{X: upCoord.X, Z: (coord.Z/rHeight - 1) * rHeight}
			// Upper and lower rectangles are definitely crossed by the line
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
	// Traverse vertical axis
	startX := xStart.X / rWidth
	endX := xEnd.X / rWidth
	for i := startX; i <= endX; i++ {
		q0 := &Coord{X: i * rWidth, Z: 0}
		q1 := &Coord{X: i * rWidth, Z: height}
		// Calculate intersection point
		coord, ok := GetCrossCoord(p0, p1, *q0, *q1)
		if ok {
			coord.X = i * rWidth
			rightCoord := Coord{X: (coord.X / rWidth) * rWidth, Z: (coord.Z / rWidth) * rWidth}
			leftCoord := Coord{X: (coord.X/rWidth - 1) * rWidth, Z: rightCoord.Z}
			// Left and right rectangles are definitely crossed by the line
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

// GetIntersectRect calculates the intersection area between two rectangles
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
