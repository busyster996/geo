package geo

// LocationState represents the position state
type LocationState int

// LocationState constants
const (
	LeftTop = 1 << iota
	RightTop
	LeftBottom  // Fixed typo: Buttom -> Bottom
	RightBottom // Fixed typo: Buttom -> Bottom
)

// Border represents a boundary
type Border struct {
	Rectangle
}

// NewBorder creates a new border
func NewBorder(x, z, width, height int32) Border {
	return Border{
		Rectangle: NewRectangle(x, z, width, height),
	}
}

// RectLocation determines the boundary position of a rectangle
func (b *Border) RectLocation(minX, minZ, maxX, maxZ int32) LocationState {
	if minX > b.X+b.Width ||
		minZ > b.Z+b.Height ||
		maxX < b.X ||
		maxZ < b.Z {
		return 0
	}

	var location LocationState
	centerX := b.X + b.Width/2
	centerZ := b.Z + b.Height/2
	if minX <= centerX {
		if maxZ >= centerZ {
			location |= LeftTop
		}
		if minZ <= centerZ {
			location |= LeftBottom
		}
	}
	if maxX > centerZ {
		if maxZ > centerZ {
			location |= RightTop
		}
		if minZ < centerZ {
			location |= RightBottom
		}
	}
	return location
}

// CoordLocation returns the position of a coordinate point within the Border
func (b *Border) CoordLocation(p Coord) LocationState {
	// Outside boundary
	if !b.IsCoordInside(p) {
		return 0
	}

	centerX := b.X + b.Width/2
	centerZ := b.Z + b.Height/2
	if p.X <= centerX {
		if p.Z >= centerZ {
			return LeftTop
		}
		return LeftBottom
	}
	if p.Z >= centerZ {
		return RightTop
	}
	return RightBottom
}
