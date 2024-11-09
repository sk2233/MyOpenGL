/*
@author: sk
@date: 2024/10/27
*/
package main

import "math"

type Mat4 struct {
	Mat [4][4]float64
}

func NewIdentityMat4() *Mat4 {
	return &Mat4{
		Mat: [4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func NewScaleMat4(sx, sy, sz float64) *Mat4 {
	return &Mat4{
		Mat: [4][4]float64{
			{sx, 0, 0, 0},
			{0, sy, 0, 0},
			{0, 0, sz, 0},
			{0, 0, 0, 1},
		},
	}
}

func NewTranslateMat4(tx, ty, tz float64) *Mat4 {
	return &Mat4{
		Mat: [4][4]float64{
			{1, 0, 0, 0},
			{0, 1, 0, 0},
			{0, 0, 1, 0},
			{tx, ty, tz, 1},
		},
	}
}

func NewRotateXMat4(angle float64) *Mat4 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return &Mat4{
		Mat: [4][4]float64{
			{1, 0, 0, 0},
			{0, cos, sin, 0},
			{0, -sin, cos, 0},
			{0, 0, 0, 1},
		},
	}
}

func NewRotateYMat4(angle float64) *Mat4 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return &Mat4{
		Mat: [4][4]float64{
			{cos, 0, -sin, 0},
			{0, 1, 0, 0},
			{sin, 0, cos, 0},
			{0, 0, 0, 1},
		},
	}
}

func NewRotateZMat4(angle float64) *Mat4 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return &Mat4{
		Mat: [4][4]float64{
			{cos, sin, 0, 0},
			{-sin, cos, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 0, 1},
		},
	}
}

func NewLookAtMat4Raw(x, y, z Vec3, pos Vec3) *Mat4 {
	return &Mat4{ // 先平移使相机位于世界空间的原点 再旋转使世界空间方向与相机一致
		Mat: [4][4]float64{
			{x.X, y.X, z.X, 0},
			{x.Y, y.Y, z.Y, 0},
			{x.Z, y.Z, z.Z, 0},
			{-pos.Dot(x), -pos.Dot(y), -pos.Dot(z), 1}, // 这里相当于直接把平移矩阵与旋转矩阵结合在一起了
		},
	}
}

// 更常用的方法
func NewLookAtMat4(eye, target, up Vec3) *Mat4 {
	z := eye.Sub(target).Normal()
	x := up.Cross(z).Normal()
	y := z.Cross(x).Normal()
	return NewLookAtMat4Raw(x, y, z, eye)
}

/*
fovy 视野角度
aspect 宽高比
near 近平面
far 远平面     透视矩阵
*/
func NewPerspectiveMat4(fovy float64, aspect float64, near, far float64) *Mat4 {
	zRange := far - near
	return &Mat4{
		Mat: [4][4]float64{
			{1 / math.Tan(fovy/2) / aspect, 0, 0, 0},
			{0, 1 / math.Tan(fovy/2), 0, 0},
			{0, 0, -(near + far) / zRange, -1},
			{0, 0, -2 * near * far / zRange, 0},
		},
	}
}

func (m *Mat4) Mul(mat4 *Mat4) *Mat4 {
	res := &Mat4{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				res.Mat[i][j] += m.Mat[i][k] * mat4.Mat[k][j]
			}
		}
	}
	return res
}
