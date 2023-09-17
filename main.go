package main

import (
	"bytes"
	"embed"
	"flag"
	"image"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"go.imnhan.com/shark/must"
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

//go:embed sprites/feed/*
var FeedSprites embed.FS

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

var DurationTillHungry time.Duration
var WalkChance, StopChance int
var WindowWidth, WindowHeight int

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

func GlobalCursorPosition() Vector {
	cx, cy := ebiten.CursorPosition()
	wx, wy := ebiten.WindowPosition()
	return Vector{cx + wx, cy + wy}
}

func randBool(chance int) bool {
	return rand.Intn(100) < chance
}

func NewAnim(sprites embed.FS, subdir string) *Anim {
	files := must.One(sprites.ReadDir("sprites/" + subdir))
	var frames []*ebiten.Image
	for _, direntry := range files {
		fname := direntry.Name()
		frame := must.One(sprites.ReadFile("sprites/" + subdir + "/" + fname))
		img, _ := must.Two(ebitenutil.NewImageFromReader(bytes.NewReader(frame)))
		frames = append(frames, img)
	}
	return &Anim{frames}
}

var Idle, RightClick, Drag, Hungry, Feed, WalkLeft, WalkRight *Anim

func init() {
	Idle = NewAnim(IdleSprites, "idle")
	Drag = NewAnim(DragSprites, "drag")
	RightClick = NewAnim(RightClickSprites, "right-click")
	Hungry = NewAnim(HungrySprites, "hungry")
	Feed = NewAnim(FeedSprites, "feed")
	WalkLeft = NewAnim(WalkLeftSprites, "walk-left")
	WalkRight = NewAnim(WalkRightSprites, "walk-right")
}

func main() {
	var sizeFlag, xFlag, yFlag int
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
	flag.IntVar(&WalkChance, "walk", 5, "Chance to start walking, in %")
	flag.IntVar(&StopChance, "stop", 40, "Chance to stop walking, in %")
	flag.Parse()

	DurationTillHungry = time.Duration(secondsUntilHungryFlag) * 1_000_000_000

	ebiten.SetWindowPosition(xFlag, yFlag)
	WindowWidth = SPRITE_X * sizeFlag
	WindowHeight = SPRITE_Y * sizeFlag
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Shark!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)

	AppIcon, _ := must.Two(image.Decode(bytes.NewReader(IconFile)))
	ebiten.SetWindowIcon([]image.Image{AppIcon})

	must.Zero(ebiten.RunGameWithOptions(
		NewStateMachine(),
		&ebiten.RunGameOptions{ScreenTransparent: true},
	))
}
