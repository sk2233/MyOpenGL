/*
@author: sk
@date: 2024/10/27
*/
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	Pos         Vec3
	Up          Vec3
	Dir         Vec3
	Left        Vec3
	MoveSpeed   float64
	RotateSpeed float64
}

func (c *Camera) Update() {
	// 移动
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		c.Pos = c.Pos.Add(c.Up.Mul(c.MoveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		c.Pos = c.Pos.Sub(c.Up.Mul(c.MoveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Pos = c.Pos.Add(c.Dir.Mul(c.MoveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Pos = c.Pos.Sub(c.Dir.Mul(c.MoveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		c.Pos = c.Pos.Add(c.Left.Mul(c.MoveSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		c.Pos = c.Pos.Sub(c.Left.Mul(c.MoveSpeed))
	}
	// 旋转
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		mat4 := NewRotateYMat4(c.RotateSpeed)
		c.Dir = c.Dir.ToVec4().Mul(mat4).ToVec3()
		c.Left = c.Left.ToVec4().Mul(mat4).ToVec3()
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		mat4 := NewRotateYMat4(-c.RotateSpeed)
		c.Dir = c.Dir.ToVec4().Mul(mat4).ToVec3()
		c.Left = c.Left.ToVec4().Mul(mat4).ToVec3()
	}
}

func (c *Camera) GetView() *Mat4 {
	return NewLookAtMat4(c.Pos, c.Pos.Add(c.Dir), Vec3Up)
}

func NewCamera() *Camera {
	return &Camera{
		Pos:         Vec3{0, 0, 0},
		Up:          Vec3{0, 1, 0},
		Dir:         Vec3{0, 0, -1},
		Left:        Vec3{-1, 0, 0},
		MoveSpeed:   5.0 / 60,
		RotateSpeed: 1.0 / 60,
	}
}
