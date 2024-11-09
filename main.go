/*
@author: sk
@date: 2024/10/27
*/
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(WinW, WinH)
	ebiten.SetScreenClearedEveryFrame(false) // 自己控制刷新
	err := ebiten.RunGame(NewApp())
	HandleErr(err)
}
