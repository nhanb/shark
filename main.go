package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const SPRITE_X = 100
const SPRITE_Y = 135

//go:embed sprites/idle/*
var IdleSprites embed.FS

//go:embed sprites/right-click/*
var RightClickSprites embed.FS

//go:embed sprites/drag/*
var DragSprites embed.FS

//go:embed sprites/hungry/*
var HungrySprites embed.FS

//go:embed sprites/feeding/*
var FeedingSprites embed.FS

//go:embed sprites/walk-left/*
var WalkLeftSprites embed.FS

//go:embed sprites/walk-right/*
var WalkRightSprites embed.FS

//go:embed icon.png
var IconFile []byte

type Anim struct {
	Frames []*ebiten.Image
}

type Position struct{ x, y int }

type Game struct {
	CurrentAnim            *Anim
	CurrentFrame           int
	Ticks                  int
	IsDragging             bool
	PreviousMousePos       Vector
	SpriteStartPos         Vector
	MouseStartPos          Vector
	Size                   int
	LastFed                time.Time
	NanosecondsUntilHungry time.Duration
	WalkChance             int
	StopChance             int
	ScreenSize             Vector
	SpritePos              Vector
	op                     *ebiten.DrawImageOptions
}

type Vector struct{ x, y int }

func CreateVector(x, y int) Vector {
	return Vector{x, y}
}

func (this Vector) Add(that Vector) Vector {
	return Vector{this.x + that.x, this.y + that.y}
}
func (this Vector) Subtract(that Vector) Vector {
	return Vector{this.x - that.x, this.y - that.y}
}

func (g *Game) Update() error {
	isHungry := false

	if time.Now().Sub(g.LastFed) >= g.NanosecondsUntilHungry {
		// The only allowed interaction when hungry is right-click to feed.
		isHungry = true
		g.IsDragging = false
		if g.CurrentAnim != Hungry {
			g.CurrentAnim = Hungry
			g.Ticks = 0
			g.CurrentFrame = 0
			return nil
		} else if IsSpriteJustPressed(g, ebiten.MouseButtonRight) {
			g.CurrentAnim = Feeding
			g.Ticks = 0
			g.CurrentFrame = 0
			g.LastFed = time.Now()
			return nil
		}
	}

	if !isHungry && g.CurrentAnim != Feeding {
		handleNonHungryInputs(g)
	}

	switch g.CurrentAnim {
	case WalkLeft:
		x, y := ebiten.WindowPosition()
		ebiten.SetWindowPosition(x-g.Size, y)
	case WalkRight:
		x, y := ebiten.WindowPosition()
		ebiten.SetWindowPosition(x+g.Size, y)
	}

	g.Ticks++
	if g.Ticks < 10 {
		return nil
	}
	g.Ticks = 0
	g.CurrentFrame++

	if g.CurrentFrame >= len(g.CurrentAnim.Frames) {
		g.CurrentFrame = 0
		if g.CurrentAnim == RightClick || g.CurrentAnim == Feeding {
			g.CurrentAnim = Idle
		}

		if g.CurrentAnim == Idle {
			if randBool(g.WalkChance) {
				if randBool(50) {
					g.CurrentAnim = WalkLeft
				} else {
					g.CurrentAnim = WalkRight
				}
			}
		} else if g.CurrentAnim == WalkLeft || g.CurrentAnim == WalkRight {
			if randBool(g.StopChance) {
				g.CurrentAnim = Idle
			}
		}
	}
	return nil
}

func randBool(chance int) bool {
	return rand.Intn(100) < chance
}

func IsSpriteJustPressed(g *Game, btn ebiten.MouseButton) bool {
	if !inpututil.IsMouseButtonJustPressed(btn) {
		return false
	}
	minX := g.SpritePos.x
	minY := g.SpritePos.y
	maxX := g.SpritePos.x + SPRITE_X*g.Size
	maxY := g.SpritePos.y + SPRITE_Y*g.Size
	x, y := ebiten.CursorPosition()
	return x >= minX && x <= maxX && y >= minY && y <= maxY
}

func handleNonHungryInputs(g *Game) {
	if IsSpriteJustPressed(g, ebiten.MouseButtonRight) {
		if g.CurrentAnim == Idle {
			g.CurrentAnim = RightClick
			g.Ticks = 0
			g.CurrentFrame = 0
		}
	}

	if IsSpriteJustPressed(g, ebiten.MouseButtonLeft) {
		g.IsDragging = true
		g.CurrentAnim = Drag
		g.Ticks = 0
		g.CurrentFrame = 0
		g.PreviousMousePos = CreateVector(ebiten.CursorPosition())
		g.SpriteStartPos = g.SpritePos
		g.MouseStartPos = CreateVector(ebiten.CursorPosition())
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.IsDragging = false
		g.CurrentAnim = Idle
		g.Ticks = 0
		g.CurrentFrame = 0
	}

	mousePos := CreateVector(ebiten.CursorPosition())
	if g.IsDragging && mousePos != g.PreviousMousePos {
		g.SpritePos = g.SpriteStartPos.Add(mousePos.Subtract(g.MouseStartPos))
	}

	g.PreviousMousePos = mousePos
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(g.SpritePos.x), float64(g.SpritePos.y))
	screen.DrawImage(g.CurrentAnim.Frames[g.CurrentFrame], g.op)
	/*
		debugStr := ""
		debugStr += fmt.Sprintf("%v\n", g.Ticks)
		debugStr += fmt.Sprintf("%v\n", g.LastFed)
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			debugStr += "Dragging\n"
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			debugStr += "Right click\n"
		}
		ebitenutil.DebugPrint(screen, debugStr)
	*/
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return outsideWidth, outsideHeight
}

func NewAnim(sprites embed.FS, subdir string) *Anim {
	files, err := sprites.ReadDir("sprites/" + subdir)
	PanicIfErr(err)
	var frames []*ebiten.Image
	for _, direntry := range files {
		fname := direntry.Name()
		frame, err := sprites.ReadFile("sprites/" + subdir + "/" + fname)
		PanicIfErr(err)
		img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(frame))
		PanicIfErr(err)
		frames = append(frames, img)
	}
	return &Anim{frames}
}

var Idle, RightClick, Drag, Hungry, Feeding, WalkLeft, WalkRight *Anim

func init() {
	Idle = NewAnim(IdleSprites, "idle")
	Drag = NewAnim(DragSprites, "drag")
	RightClick = NewAnim(RightClickSprites, "right-click")
	Hungry = NewAnim(HungrySprites, "hungry")
	Feeding = NewAnim(FeedingSprites, "feeding")
	WalkLeft = NewAnim(WalkLeftSprites, "walk-left")
	WalkRight = NewAnim(WalkRightSprites, "walk-right")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var sizeFlag, xFlag, yFlag, walkChanceFlag, stopChanceFlag int
	var secondsUntilHungryFlag int64
	flag.IntVar(
		&sizeFlag, "size", 1, "Size multiplier: make Gura as big as you want",
	)
	flag.Int64Var(
		&secondsUntilHungryFlag,
		"hungry",
		3600,
		"The number of seconds it takes for Gura to go hungry",
	)
	flag.IntVar(&xFlag, "x", 9999, "X position on screen")
	flag.IntVar(&yFlag, "y", 9999, "Y position on screen")
	flag.IntVar(&walkChanceFlag, "walk", 5, "chance to start walking, in %")
	flag.IntVar(&stopChanceFlag, "stop", 40, "chance to stop walking, in %")
	flag.Parse()

	var game Game
	game.CurrentAnim = Idle
	game.LastFed = time.Now()
	game.NanosecondsUntilHungry = time.Duration(secondsUntilHungryFlag) * 1_000_000_000
	game.Size = sizeFlag
	game.WalkChance = walkChanceFlag
	game.StopChance = stopChanceFlag
	game.ScreenSize = CreateVector(ebiten.ScreenSizeInFullscreen())

	if xFlag > game.ScreenSize.x-SPRITE_X {
		xFlag = game.ScreenSize.x - SPRITE_X
	}
	if yFlag > game.ScreenSize.x-SPRITE_Y {
		yFlag = game.ScreenSize.y - SPRITE_Y
	}
	game.SpritePos = Vector{xFlag, yFlag}

	game.op = &ebiten.DrawImageOptions{}
	game.op.GeoM.Translate(float64(game.SpritePos.x), float64(game.SpritePos.y))
	fmt.Println(game.SpritePos)

	ebiten.SetWindowTitle("Shark!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowSize(game.ScreenSize.x, game.ScreenSize.y)
	ebiten.SetWindowPosition(0, 0)

	AppIcon, _, iconerr := image.Decode(bytes.NewReader(IconFile))
	PanicIfErr(iconerr)
	ebiten.SetWindowIcon([]image.Image{AppIcon})

	err := ebiten.RunGame(&game)
	PanicIfErr(err)
}

func PanicIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
