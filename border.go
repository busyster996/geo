package geo

// LocationState 位置状态
type LocationState int

// LocationState
const (
	LeftTop = 1 << iota
	RightTop
	LeftButtom
	RightButtom
)

// Border 边界
type Border struct {
	Rectangle
}

// NewBorder 创建新的边界
func NewBorder(x, z, width, height int32) Border {
	return Border{
		Rectangle: NewRectangle(x, z, width, height),
	}
}

// RectLocation 判断边界位置
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
			location |= LeftButtom
		}
	}
	if maxX > centerZ {
		if maxZ > centerZ {
			location |= RightTop
		}
		if minZ < centerZ {
			location |= RightButtom
		}
	}
	return location
}

// CoordLocation 点在Border中的位置
func (b *Border) CoordLocation(p Coord) LocationState {
	// 边界外
	if !b.IsCoordInside(p) {
		return 0
	}

	centerX := b.X + b.Width/2
	centerZ := b.Z + b.Height/2
	if p.X <= centerX {
		if p.Z >= centerZ {
			return LeftTop
		}
		return LeftButtom
	}
	if p.Z >= centerZ {
		return RightTop
	}
	return RightButtom
}
