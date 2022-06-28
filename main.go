package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	CurrentFrame int
	Ticks        int
	Sprites      []*ebiten.Image
}

func (g *Game) Update() error {
	g.Ticks++
	if g.Ticks >= 10 {
		g.Ticks = 0
		g.CurrentFrame++
		if g.CurrentFrame > 3 {
			g.CurrentFrame = 0
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.WindowPosition()
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"Ticks: %d\nCurrentFrame: %d\nx: %v, y: %v",
			g.Ticks, g.CurrentFrame, x, y,
		),
	)
	screen.DrawImage(g.Sprites[g.CurrentFrame], nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return outsideWidth / 2, outsideHeight / 2
}

func PanicIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var game Game

	// Should probably use go:embed somehow here
	for i := 0; i <= 3; i++ {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("sprites/idle-%d.png", i))
		PanicIfErr(err)
		game.Sprites = append(game.Sprites, img)
	}

	ebiten.SetWindowSize(360, 360)
	ebiten.SetWindowTitle("Shark!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowPosition(9999, 9999)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
