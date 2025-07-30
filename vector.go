package geo

import "math"

// Vector 位置向量
// 从坐标原点指向点所在位置的矢量称为位置矢量，简称位矢。
// https://zh.wikipedia.org/wiki/%E4%BD%8D%E7%BD%AE%E5%90%91%E9%87%8F
type Vector struct {
	X, Z int32
}

// NewVector 新建
func NewVector(start, end Coord) Vector {
	return Vector{
		X: end.X - start.X,
		Z: end.Z - start.Z,
	}
}

// NewVectorByCoord 新建
func NewVectorByCoord(p Coord) Vector {
	return Vector(p)
}

// Add 向量的加法
func (v *Vector) Add(vec *Vector) Vector {
	return Vector{
		X: v.X + vec.X,
		Z: v.Z + vec.Z,
	}
}

// Minus 向量的减法
func (v *Vector) Minus(vec *Vector) Vector {
	return Vector{
		X: v.X - vec.X,
		Z: v.Z - vec.Z,
	}
}

// Dot 向量的点积
func (v *Vector) Dot(vec *Vector) float64 {
	result := int64(v.X) * int64(vec.X)
	result += int64(v.Z) * int64(vec.Z)
	return float64(result)
}

// Cross 向量的叉乘
// result > 0 vec在v左侧; reslut <0 vec在v右; result.Y = 0 共线
func (v *Vector) Cross(vec *Vector) int64 {
	return int64(v.X)*int64(vec.Z) - int64(v.Z)*int64(vec.X)
}

// LengthSquared 距离平方
func (v *Vector) LengthSquared() float64 {
	return float64(int64(v.X)*int64(v.X) + int64(v.Z)*int64(v.Z))
}

// Length 向量的长度
func (v *Vector) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Trunc 向量截断
func (v *Vector) Trunc(ratio float64) Vector {
	return Vector{
		X: int32(ratio * float64(v.X)),
		Z: int32(ratio * float64(v.Z)),
	}
}

// TruncEdge 截断边为一个单位向量
func TruncEdge(start, end Coord) Coord {
	// 生成拐点数组
	vec0 := NewVector(start, end)
	// 单位向量
	vec0 = vec0.Trunc(1000 / vec0.Length())

	return vec0.ToCoord(start)
}

// ToCoord 转换成坐标
func (v Vector) ToCoord(start Coord) Coord {
	return Coord{
		X: start.X + v.X,
		Z: start.Z + v.Z,
	}
}

// Rotate 向量按照角度旋转获得新的向量，左手坐标系，旋转的时候默认向左边旋转
// 对于任意两个不同点A和B，A绕B旋转θ角度后的坐标为：
// (Δx*cosθ- Δy * sinθ+ xB, Δy*cosθ + Δx * sinθ+ yB )
// 注：xB、yB为B点坐标。
// https://blog.csdn.net/u013445530/article/details/44904017
func (v *Vector) Rotate(angle float64) Vector {
	x0 := float64(v.X)
	z0 := float64(v.Z)

	cos := math.Cos(angle)
	sin := math.Sin(angle)

	x := x0*cos - z0*sin
	z := x0*sin + z0*cos

	return Vector{
		X: int32(x),
		Z: int32(z),
	}
}

// CalCoordDst 获取点到向量的距离
func (v *Vector) CalCoordDst(start, target Coord) float64 {
	vec := NewVector(start, target)

	angle := v.GetAngle(&vec)
	return vec.Length() * math.Sin(angle)
}

// GetAngle 获取俩个向量的角度
func (v *Vector) GetAngle(vec *Vector) float64 {
	a := v.Dot(vec)
	b := v.Length() * vec.Length()
	t := a / b
	angle := math.Acos(t)
	if math.IsNaN(angle) {
		if t > 0 {
			return 0
		}
		return math.Pi
	}

	return angle
}

// cross 叉乘算法,
// 积=p3p1 X p3p2
func cross(p1, p2, p3 Coord) int64 {
	s := int64(p1.X-p3.X)*int64((p2.Z-p3.Z)) - int64(p2.X-p3.X)*int64(p1.Z-p3.Z)
	return s
}

// CalCoordByRatio Ratio表示newVec和Vec长度的比值
// 返回newVec终点坐标
func CalCoordByRatio(startCoord, endCoord Coord, ratio float64) Coord {
	vec := NewVector(startCoord, endCoord)
	return vec.Trunc(ratio).ToCoord(startCoord)
}
