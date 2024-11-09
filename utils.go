/*
@author: sk
@date: 2024/10/27
*/
package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func LerpVec4(pre, curr Vec4, rate float64) Vec4 {
	return Vec4{
		X: curr.X*rate + pre.X*(1-rate),
		Y: curr.Y*rate + pre.Y*(1-rate),
		Z: curr.Z*rate + pre.Z*(1-rate),
		W: curr.W*rate + pre.W*(1-rate),
	}
}

func LerpVec3(pre, curr Vec3, rate float64) Vec3 {
	return Vec3{
		X: curr.X*rate + pre.X*(1-rate),
		Y: curr.Y*rate + pre.Y*(1-rate),
		Z: curr.Z*rate + pre.Z*(1-rate),
	}
}

func LerpVec2(pre Vec2, curr Vec2, rate float64) Vec2 {
	return Vec2{
		X: curr.X*rate + pre.X*(1-rate),
		Y: curr.Y*rate + pre.Y*(1-rate),
	}
}

func FillArray[T any](arr []T, val T) {
	for i := 0; i < len(arr); i++ {
		arr[i] = val
	}
}

func SetColor(colorBuff []byte, x int, y int, clr Color) {
	if x < 0 || x >= WinW || y < 0 || y >= WinH { // 临时使用了 WinW
		return
	}
	index := (x + y*WinW) * 4
	colorBuff[index] = clr[0]
	colorBuff[index+1] = clr[1]
	colorBuff[index+2] = clr[2]
	colorBuff[index+3] = clr[3]
}

func WeightVec4(v1, v2, v3 Vec4, weight Vec3) Vec4 {
	return Vec4{
		X: v1.X*weight.X + v2.X*weight.Y + v3.X*weight.Z,
		Y: v1.Y*weight.X + v2.Y*weight.Y + v3.Y*weight.Z,
		Z: v1.Z*weight.X + v2.Z*weight.Y + v3.Z*weight.Z,
		W: v1.W*weight.X + v2.W*weight.Y + v3.W*weight.Z,
	}
}

func WeightVec3(v1 Vec3, v2 Vec3, v3 Vec3, weight Vec3) Vec3 {
	return Vec3{
		X: v1.X*weight.X + v2.X*weight.Y + v3.X*weight.Z,
		Y: v1.Y*weight.X + v2.Y*weight.Y + v3.Y*weight.Z,
		Z: v1.Z*weight.X + v2.Z*weight.Y + v3.Z*weight.Z,
	}
}

func WeightVec2(v1 Vec2, v2 Vec2, v3 Vec2, weight Vec3) Vec2 {
	return Vec2{
		X: v1.X*weight.X + v2.X*weight.Y + v3.X*weight.Z,
		Y: v1.Y*weight.X + v2.Y*weight.Y + v3.Y*weight.Z,
	}
}

func IsBack(a, b, c Vec4) bool {
	// 逆时针认为是可见的
	return a.X*b.Y-a.Y*b.X+b.X*c.Y-b.Y*c.X+c.X*a.Y-c.Y*a.X < 0
}

func LoadObjs(path string) []*Triangle {
	reader, err := os.Open(path)
	HandleErr(err)
	defer reader.Close()

	res := make([]*Triangle, 0)
	vs := make([]Vec3, 0)
	vts := make([]Vec2, 0)
	vns := make([]Vec3, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(strings.TrimSpace(line), " ")
		if len(items) < 3 { // 至少有 3 个
			continue
		}
		switch items[0] {
		case "v":
			vs = append(vs, Vec3{X: ParseFloat(items[1]), Y: ParseFloat(items[2]), Z: ParseFloat(items[3])})
		case "vt":
			vts = append(vts, Vec2{X: ParseFloat(items[1]), Y: ParseFloat(items[2])})
		case "vn":
			vns = append(vns, Vec3{X: ParseFloat(items[1]), Y: ParseFloat(items[2]), Z: ParseFloat(items[3])})
		case "f":
			res = append(res, &Triangle{ParsePoint(vs, vts, vns, items[1]), ParsePoint(vs, vts, vns, items[2]),
				ParsePoint(vs, vts, vns, items[3])})
		default:
			panic(fmt.Errorf("unknown obj type: %s", items[0]))
		}
	}
	return res
}

func ParsePoint(vs []Vec3, vts []Vec2, vns []Vec3, val string) *Point {
	items := strings.Split(strings.TrimSpace(val), "/")
	return &Point{
		Pos: vs[ParseInt(items[0])-1].ToVec4(),
		Tex: vts[ParseInt(items[1])-1],
		Nor: vns[ParseInt(items[2])-1],
	}
}

func ParseInt(val string) int {
	res, err := strconv.ParseInt(val, 10, 64)
	HandleErr(err)
	return int(res)
}

func ParseFloat(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	HandleErr(err)
	return res
}

func LoadImage(path string) image.Image {
	reader, err := os.Open(path)
	HandleErr(err)
	defer reader.Close()
	img, err := png.Decode(reader)
	HandleErr(err)
	return img
}

func GetColor(img image.Image, tex Vec2) Vec3 {
	bound := img.Bounds()
	clr := img.At(int(tex.X*float64(bound.Dx())), int(tex.Y*float64(bound.Dy()))).(color.NRGBA)
	return Vec3{float64(clr.R) / 0xFF, float64(clr.G) / 0xFF, float64(clr.B) / 0xFF}
}
