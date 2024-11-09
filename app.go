/*
@author: sk
@date: 2024/10/27
*/
package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type App struct {
	ColorBuff      []byte          // 颜色缓冲 w*h*4  R G B A
	DepthBuff      []float64       // 深度缓冲 w*h
	VertexShader   *VertexShader   // 顶点着色器
	FragmentShader *FragmentShader // 片源着色器
	MVP            *Mat4
	Triangles      []*Triangle
	Camera         *Camera
}

func NewApp() *App {
	obj := LoadObjs("res/box.obj")
	return &App{
		ColorBuff:      make([]byte, WinW*WinH*4),
		DepthBuff:      make([]float64, WinW*WinH),
		VertexShader:   NewVertexShader(),
		FragmentShader: NewFragmentShader(),
		Triangles:      obj,
		Camera:         NewCamera(),
	}
}

func (a *App) Update() error {
	a.Camera.Update()
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	model := NewIdentityMat4()
	view := a.Camera.GetView()
	proj := NewPerspectiveMat4(math.Pi/3, 16.0/9.0, 0.1, 100)
	a.MVP = model.Mul(view).Mul(proj)
	// 顶点着色器
	for _, triangle := range a.Triangles {
		for i := 0; i < 3; i++ {
			a.VertexShader.Transform(triangle[i], a.MVP, model)
		}
	}
	// 三角形裁剪
	triangles := make([]*Triangle, 0)
	for _, triangle := range a.Triangles {
		triangles = append(triangles, Clip(triangle)...)
	}
	// 齐次变换 (变化后就有了近大远小的特征了)
	for _, triangle := range triangles {
		for i := 0; i < 3; i++ {
			CalculateNdc(triangle[i])
		}
	}
	// 计算最终屏幕空间的位置
	for _, triangle := range triangles {
		for i := 0; i < 3; i++ {
			CalculateFrag(triangle[i], WinW, WinH)
		}
	}
	// 片源着色器
	FillArray(a.ColorBuff, 0) // 使用前先清空
	FillArray(a.DepthBuff, math.MaxFloat64)
	a.FragmentShader.EyePos = a.Camera.Pos // 准备参数

	for _, triangle := range triangles {
		// 这里是通过 2 维平面判断顺时针还是逆时针
		if IsBack(triangle[0].NdcPos, triangle[1].NdcPos, triangle[2].NdcPos) {
			continue
		}
		a.FragmentShader.Draw(a.ColorBuff, a.DepthBuff, triangle)
	}
	screen.WritePixels(a.ColorBuff) // 直接覆盖提高性能 必须每次重新覆盖
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (a *App) Layout(w, h int) (int, int) {
	return w, h
}
