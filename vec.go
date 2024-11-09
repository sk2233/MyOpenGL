/*
@author: sk
@date: 2024/10/27
*/
package main

import "math"

type Vec2 struct {
	X, Y float64
}

var (
	Vec3Up = Vec3{0, 1, 0}
)

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Dot(vec3 Vec3) float64 {
	return v.X*vec3.X + v.Y*vec3.Y + v.Z*vec3.Z
}

func (v Vec3) Cross(vec3 Vec3) Vec3 {
	return Vec3{
		X: v.Y*vec3.Z - v.Z*vec3.Y,
		Y: v.Z*vec3.X - v.X*vec3.Z,
		Z: v.X*vec3.Y - v.Y*vec3.X,
	}
}

func (v Vec3) Sub(vec3 Vec3) Vec3 {
	return Vec3{X: v.X - vec3.X, Y: v.Y - vec3.Y, Z: v.Z - vec3.Z}
}

func (v Vec3) Normal() Vec3 {
	l := v.Len()
	return v.Div(l)
}

func (v Vec3) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Div(val float64) Vec3 {
	return Vec3{X: v.X / val, Y: v.Y / val, Z: v.Z / val}
}

func (v Vec3) Mul(val float64) Vec3 {
	return Vec3{
		X: v.X * val,
		Y: v.Y * val,
		Z: v.Z * val,
	}
}

func (v Vec3) Add(val Vec3) Vec3 {
	return Vec3{
		X: v.X + val.X,
		Y: v.Y + val.Y,
		Z: v.Z + val.Z,
	}
}

func (v Vec3) ToVec4() Vec4 {
	return Vec4{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
		W: 1,
	}
}

func (v Vec3) Reflect(vn Vec3) Vec3 {
	// 计算 v 以 vn 为法线的反射光线 v vn 都是单位向量
	return v.Sub(vn.Scale(2 * v.Dot(vn)))
}

func (v Vec3) Scale(val float64) Vec3 {
	return Vec3{X: v.X * val, Y: v.Y * val, Z: v.Z * val}
}

type Vec4 struct {
	X, Y, Z, W float64
}

func (v Vec4) ToVec3() Vec3 {
	return Vec3{v.X, v.Y, v.Z}
}

func (v Vec4) Mul(mat4 *Mat4) Vec4 {
	return Vec4{
		X: v.X*mat4.Mat[0][0] + v.Y*mat4.Mat[1][0] + v.Z*mat4.Mat[2][0] + v.W*mat4.Mat[3][0],
		Y: v.X*mat4.Mat[0][1] + v.Y*mat4.Mat[1][1] + v.Z*mat4.Mat[2][1] + v.W*mat4.Mat[3][1],
		Z: v.X*mat4.Mat[0][2] + v.Y*mat4.Mat[1][2] + v.Z*mat4.Mat[2][2] + v.W*mat4.Mat[3][2],
		W: v.X*mat4.Mat[0][3] + v.Y*mat4.Mat[1][3] + v.Z*mat4.Mat[2][3] + v.W*mat4.Mat[3][3],
	}
}

func (v Vec4) Sub(vec4 Vec4) Vec4 {
	return Vec4{
		X: v.X - vec4.X,
		Y: v.Y - vec4.Y,
		Z: v.Z - vec4.Z,
		W: v.W - vec4.W,
	}
}

func (v Vec4) Scale(val float64) Vec4 {
	return Vec4{
		X: v.X * val,
		Y: v.Y * val,
		Z: v.Z * val,
		W: v.W * val,
	}
}
