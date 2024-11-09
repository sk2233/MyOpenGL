/*
@author: sk
@date: 2024/10/27
*/
package main

import (
	"image"
	"math"
)

type VertexShader struct {
}

func NewVertexShader() *VertexShader {
	return &VertexShader{}
}

func (s *VertexShader) Transform(point *Point, mvp *Mat4, model *Mat4) {
	point.ClipPos = point.Pos.Mul(mvp)
	point.WorldPos = point.Pos.Mul(model)
	point.WorldNor = point.Nor.ToVec4().Mul(model).ToVec3()
}

type FragmentShader struct {
	LightPos   Vec3
	Ambient    float64 // 环境光
	Diffuse    float64 // 漫反射
	Specular   Vec3    // 高光颜色
	Shininess  float64 // 光滑度
	EyePos     Vec3    // 相机位置
	Texture    image.Image
	NorTexture image.Image
}

func NewFragmentShader() *FragmentShader {
	return &FragmentShader{
		LightPos:   Vec3{0, 1, 2},
		Ambient:    0.5,
		Diffuse:    0.1,
		Specular:   Vec3{1, 1, 1},
		Shininess:  32,
		Texture:    LoadImage("res/box.png"),
		NorTexture: LoadImage("res/box_high_light.png"),
	}
}

func (s *FragmentShader) Draw(colorBuff []byte, depthBuff []float64, triangle *Triangle) {
	// 求三角形的范围
	minX, minY, maxX, maxY := math.MaxInt, math.MaxInt, 0, 0
	for i := 0; i < 3; i++ {
		pos := triangle[i].FragPos
		minX = min(minX, int(pos.X))
		minY = min(minY, int(pos.Y))
		maxX = max(maxX, int(pos.X))
		maxY = max(maxY, int(pos.Y))
	} // 防止越界
	minX = max(minX, 0)
	minY = max(minY, 0)
	maxX = min(maxX, WinW-1)
	maxY = min(maxY, WinH-1)
	// 循环填充颜色
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			pos := Vec4{X: float64(x) + 0.5, Y: float64(y) + 0.5}
			weight := CalculateWeight(pos, triangle)
			if weight.X < 0 || weight.Y < 0 || weight.Z < 0 { // 不在三角形内部
				continue
			}
			// 这个权重是只有在 线性空间有效  后面有颜色纹理也是线性空间下的
			point := &Point{
				ClipPos:  WeightVec4(triangle[0].ClipPos, triangle[1].ClipPos, triangle[2].ClipPos, weight),
				Tex:      WeightVec2(triangle[0].Tex, triangle[1].Tex, triangle[2].Tex, weight),
				WorldNor: WeightVec3(triangle[0].WorldNor, triangle[1].WorldNor, triangle[2].WorldNor, weight),
			}
			CalculateNdc(point) // ndc 是非线性空间下的变量需要重新计算
			CalculateFrag(point, WinW, WinH)
			// 深度测试
			if !s.DepthTest(depthBuff, x, y, point.FragPos.Z) {
				continue
			}
			clr, discard := s.CalculateColor(point)
			if discard {
				continue
			}
			SetColor(colorBuff, x, y, clr)
		}
	}
}

func (s *FragmentShader) CalculateColor(point *Point) (Color, bool) {
	viewDir := s.EyePos.Sub(point.WorldPos.ToVec3()).Normal()
	lightDir := s.LightPos.Sub(point.WorldPos.ToVec3()).Normal()

	ambient := s.Ambient
	// WorldNor 本来就是正规化的
	diffuse := max(point.WorldNor.Dot(lightDir), 0) * s.Diffuse // 与 0取最大值是因为防止方向的光照
	halfDir := viewDir.Add(lightDir).Div(2)
	nor := GetColor(s.NorTexture, point.Tex) // 使用高光贴图
	specular := math.Pow(max(point.WorldNor.Dot(halfDir), 0), s.Shininess) * nor.X
	clr := GetColor(s.Texture, point.Tex)

	rate := Vec3{X: min((ambient+diffuse)*clr.X+specular*s.Specular.X, 1),
		Y: min((ambient+diffuse)*clr.Y+specular*s.Specular.Y, 1),
		Z: min((ambient+diffuse)*clr.Z+specular*s.Specular.Z, 1)}
	return [4]byte{byte(0xFF * rate.X), byte(0xFF * rate.Y), byte(0xFF * rate.Z), 0xFF}, false
}

func (s *FragmentShader) DepthTest(depthBuff []float64, x int, y int, z float64) bool {
	if x < 0 || x >= WinW || y < 0 || y >= WinH {
		return false
	}
	index := x + y*WinW
	if depthBuff[index] <= z {
		return false
	}
	depthBuff[index] = z
	return true
}

// 计算插值
func CalculateWeight(pos Vec4, triangle *Triangle) Vec3 {
	ab := triangle[1].FragPos.Sub(triangle[0].FragPos)
	ac := triangle[2].FragPos.Sub(triangle[0].FragPos)
	ap := pos.Sub(triangle[0].FragPos)
	factor := 1 / (ab.X*ac.Y - ab.Y*ac.X)
	s := (ac.Y*ap.X - ac.X*ap.Y) * factor // 这里计算的是齐次变化后的权重
	t := (ab.X*ap.Y - ab.Y*ap.X) * factor
	o := 1 - s - t

	w0 := triangle[0].FragPos.W * o // 需要转换到原始坐标下的
	w1 := triangle[1].FragPos.W * s
	w2 := triangle[2].FragPos.W * t
	normal := 1 / (w0 + w1 + w2) // 归一化
	return Vec3{
		X: w0 * normal,
		Y: w1 * normal,
		Z: w2 * normal,
	}
}
