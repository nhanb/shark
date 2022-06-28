package main

import (
	"bytes"
	"embed"
	_ "image/png"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed sprites/idle/*
var IdleSprites embed.FS

//go:embed sprites/right-click/*
var RightClickSprites embed.FS

//go:embed sprites/drag/*
var DragSprites embed.FS

type Anim struct {
	Frames []*ebiten.Image
}

type Position struct{ x, y int }

type Game struct {
	CurrentAnim       *Anim
	CurrentFrame      int
	Ticks             int
	ShouldResetToIdle bool
	IsDragging        bool
	StartMouseX       int
	StartMouseY       int
}

func GlobalCursorPosition() (x, y int) {
	cx, cy := ebiten.CursorPosition()
	wx, wy := ebiten.WindowPosition()
	return cx + wx, cy + wy
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.CurrentAnim = RightClick
		g.Ticks = 0
		g.CurrentFrame = 0
		g.ShouldResetToIdle = true
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.IsDragging = true
		g.CurrentAnim = Drag
		g.Ticks = 0
		g.CurrentFrame = 0
		g.StartMouseX, g.StartMouseY = GlobalCursorPosition()
		return nil
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.IsDragging = false
		g.CurrentAnim = Idle
		g.Ticks = 0
		g.CurrentFrame = 0
		return nil
	}

	if g.IsDragging {
		currentX, currentY := GlobalCursorPosition()
		diffX := currentX - g.StartMouseX
		diffY := currentY - g.StartMouseY

		wx, wy := ebiten.WindowPosition()
		ebiten.SetWindowPosition(wx+diffX, wy+diffY)

		g.StartMouseX, g.StartMouseY = GlobalCursorPosition()
	}

	g.Ticks++
	if g.Ticks < 10 {
		return nil
	}

	g.Ticks = 0
	g.CurrentFrame++
	if g.CurrentFrame >= len(g.CurrentAnim.Frames) {
		g.CurrentFrame = 0
		if g.ShouldResetToIdle {
			g.CurrentAnim = Idle
			g.ShouldResetToIdle = false
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	/*
		x, y := ebiten.WindowPosition()
			ebitenutil.DebugPrint(
				screen,
				fmt.Sprintf(
					"Ticks: %d\nCurrentFrame: %d\nx: %v, y: %v",
					g.Ticks, g.CurrentFrame, x, y,
				),
			)
	*/
	screen.DrawImage(g.CurrentAnim.Frames[g.CurrentFrame], nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return outsideWidth / 2, outsideHeight / 2
}

func NewAnim(sprites embed.FS, subdir string) *Anim {
	files, err := sprites.ReadDir(filepath.Join("sprites", subdir))
	PanicIfErr(err)
	var frames []*ebiten.Image
	for _, direntry := range files {
		fname := direntry.Name()
		frame, err := sprites.ReadFile(filepath.Join("sprites", subdir, fname))
		PanicIfErr(err)
		img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(frame))
		PanicIfErr(err)
		frames = append(frames, img)
	}
	return &Anim{frames}
}

var Idle, RightClick, Drag *Anim

func init() {
	Idle = NewAnim(IdleSprites, "idle")
	Drag = NewAnim(DragSprites, "drag")
	RightClick = NewAnim(RightClickSprites, "right-click")
}

func main() {
	var game Game
	game.CurrentAnim = Idle

	ebiten.SetWindowSize(360, 360)
	ebiten.SetWindowTitle("Shark!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowPosition(9999, 9999)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func PanicIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
