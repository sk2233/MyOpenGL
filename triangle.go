/*
@author: sk
@date: 2024/10/27
*/
package main

type Point struct { // 包含一个点渲染过程中的各种数据
	// 入参
	Pos Vec4 // 顶点数据
	Tex Vec2 // 贴图
	Nor Vec3 // 法线
	// 中间变量
	ClipPos  Vec4
	NdcPos   Vec4
	FragPos  Vec4
	WorldPos Vec4
	WorldNor Vec3
}

type Triangle [3]*Point
