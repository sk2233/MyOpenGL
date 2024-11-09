/*
@author: sk
@date: 2024/10/27
*/
package main

import (
	"fmt"
	"math"
)

const (
	PlanPosX = 1
	PlanNegX = 2
	PlanPosY = 3
	PlanNegY = 4
	PlanPosZ = 5
	PlanNegZ = 6
)

// 对一个三角形裁剪 可能产生 0~n 个三角形
func Clip(triangle *Triangle) []*Triangle {
	// 都是可见的不用裁剪
	if PosInClip(triangle[0].ClipPos) && PosInClip(triangle[1].ClipPos) && PosInClip(triangle[2].ClipPos) {
		return []*Triangle{triangle}
	}
	points := []*Point{triangle[0], triangle[1], triangle[2]}
	// 先对 z 平面进行切割 保证 w 为负数(物体在相机后面) 的点先全部移除  否者 w 为负数 判断 x y 平面裁剪会失败
	points = CLipPlan(points, PlanPosZ)
	if len(points) == 0 { // 全裁剪了
		return make([]*Triangle, 0)
	}
	points = CLipPlan(points, PlanNegZ)
	if len(points) == 0 { // 全裁剪了
		return make([]*Triangle, 0)
	}
	points = CLipPlan(points, PlanPosX)
	if len(points) == 0 { // 全裁剪了
		return make([]*Triangle, 0)
	}
	points = CLipPlan(points, PlanNegX)
	if len(points) == 0 { // 全裁剪了
		return make([]*Triangle, 0)
	}
	points = CLipPlan(points, PlanPosY)
	if len(points) == 0 { // 全裁剪了
		return make([]*Triangle, 0)
	}
	points = CLipPlan(points, PlanNegY)
	if len(points) == 0 { // 全裁剪了
		return make([]*Triangle, 0)
	}
	triangles := make([]*Triangle, 0)
	for i := 2; i < len(points); i++ { // 重新组成 三角形
		triangles = append(triangles, &Triangle{points[0], points[i-1], points[i]})
	}
	return triangles
}

func CLipPlan(points []*Point, planType int) []*Point {
	res := make([]*Point, 0)
	cnt := len(points)
	for i := 0; i < cnt; i++ {
		currPoint := points[i]
		prePoint := points[(i-1+cnt)%cnt]
		currInPlan := InPlan(currPoint.ClipPos, planType)
		preInPlan := InPlan(prePoint.ClipPos, planType)
		if currInPlan != preInPlan { // 有穿墙计算交点
			rate := GetRate(prePoint.ClipPos, currPoint.ClipPos, planType)
			res = append(res, &Point{ // TODO 后续若是有其他变量也要跟着插值
				Pos:      LerpVec4(prePoint.Pos, currPoint.Pos, rate),
				ClipPos:  LerpVec4(prePoint.ClipPos, currPoint.ClipPos, rate),
				Tex:      LerpVec2(prePoint.Tex, currPoint.Tex, rate),
				WorldPos: LerpVec4(prePoint.WorldPos, currPoint.WorldPos, rate),
				WorldNor: LerpVec3(prePoint.WorldNor, currPoint.WorldNor, rate),
			})
		}
		if currInPlan {
			res = append(res, currPoint)
		}
	}
	return res
}

// 不能使用齐次除法后的数据，因为齐次除法会导致非线性变化，其他插值也是类似必须放在齐次除法前
func GetRate(pre Vec4, curr Vec4, planType int) float64 {
	switch planType {
	case PlanPosX: // w 可以就是对应点与轴做垂线与平面的交点 因此能用于判断是否在平面内
		return (pre.W - pre.X) / ((pre.W - pre.X) + (curr.X - curr.W))
	case PlanNegX:
		return (pre.W + pre.X) / ((pre.W + pre.X) + (-curr.X - curr.W))
	case PlanPosY:
		return (pre.W - pre.Y) / ((pre.W - pre.Y) + (curr.Y - curr.W))
	case PlanNegY:
		return (pre.W + pre.Y) / ((pre.W + pre.Y) + (-curr.Y - curr.W))
	case PlanPosZ:
		return (pre.W - pre.Z) / ((pre.W - pre.Z) + (curr.Z - curr.W))
	case PlanNegZ:
		return (pre.W + pre.Z) / ((pre.W + pre.Z) + (-curr.Z - curr.W))
	default:
		panic(fmt.Sprintf("invalid planType %d", planType))
	}
}

func InPlan(clipPos Vec4, planType int) bool {
	switch planType {
	case PlanPosX:
		return clipPos.X <= clipPos.W
	case PlanNegX:
		return clipPos.X >= -clipPos.W
	case PlanPosY:
		return clipPos.Y <= clipPos.W
	case PlanNegY:
		return clipPos.Y >= -clipPos.W
	case PlanPosZ:
		return clipPos.Z <= clipPos.W
	case PlanNegZ:
		return clipPos.Z >= -clipPos.W
	default:
		panic(fmt.Sprintf("invalid planType %d", planType))
	}
}

// 一个裁剪空间下的点看其 x y z 绝对值 是不是小于 w
func PosInClip(pos Vec4) bool {
	return math.Abs(pos.X) <= pos.W && math.Abs(pos.Y) <= pos.W && math.Abs(pos.Z) <= pos.W
}

func CalculateNdc(point *Point) {
	w := point.ClipPos.W
	point.NdcPos.X = point.ClipPos.X / w
	point.NdcPos.Y = point.ClipPos.Y / w
	point.NdcPos.Z = point.ClipPos.Z / w
	point.NdcPos.W = 1 / w // 额外记录倒数值后面有用
}

func CalculateFrag(point *Point, winW, winH float64) {
	// x y 从 -1 ~ 1 扩展到  0 ～ winSize
	point.FragPos.X = (point.NdcPos.X + 1) / 2 * winW
	point.FragPos.Y = (-point.NdcPos.Y + 1) / 2 * winH // y轴是相反的
	// z 只需要保留遮挡关系即可  从 -1 ~ 1 转换到 0 ~ 1
	point.FragPos.Z = (point.NdcPos.X + 1) / 2
	point.FragPos.W = point.NdcPos.W // w 直接保留
}
