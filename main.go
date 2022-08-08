package main

import (
	"bytes"
	"embed"
	"flag"
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
	WinStartPos            Vector
	MouseStartPos          Vector
	Size                   int
	LastFed                time.Time
	NanosecondsUntilHungry time.Duration
	WalkChance             int
	StopChance             int
	X                      int
	Y                      int
	MaxX                   int
	MaxY                   int
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

func GlobalCursorPosition() Vector {
	cx, cy := ebiten.CursorPosition()
	wx, wy := ebiten.WindowPosition()
	return Vector{cx + wx, cy + wy}
}

func (g *Game) Update() error {
	if g.X < 0 {
		g.X = 0
	} else if g.X > g.MaxX {
		g.X = g.MaxX
	}
	if g.Y < 0 {
		g.Y = 0
	} else if g.Y > g.MaxY {
		g.Y = g.MaxY
	}
	ebiten.SetWindowPosition(g.X, g.Y)

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
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
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
		g.X -= g.Size
	case WalkRight:
		g.X += g.Size
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

func handleNonHungryInputs(g *Game) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if g.CurrentAnim == Idle {
			g.CurrentAnim = RightClick
			g.Ticks = 0
			g.CurrentFrame = 0
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.IsDragging = true
		g.CurrentAnim = Drag
		g.Ticks = 0
		g.CurrentFrame = 0
		g.PreviousMousePos = GlobalCursorPosition()
		g.WinStartPos = CreateVector(ebiten.WindowPosition())
		g.MouseStartPos = GlobalCursorPosition()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.IsDragging = false
		g.CurrentAnim = Idle
		g.Ticks = 0
		g.CurrentFrame = 0
	}

	mousePos := GlobalCursorPosition()
	if g.IsDragging && mousePos != g.PreviousMousePos {
		newWinPos := g.WinStartPos.Add(mousePos.Subtract(g.MouseStartPos))
		g.X = newWinPos.x
		g.Y = newWinPos.y
		ebiten.SetWindowPosition(g.X, g.Y)
	}

	g.PreviousMousePos = mousePos
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.CurrentAnim.Frames[g.CurrentFrame], nil)
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
	return SPRITE_X, SPRITE_Y
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
	game.X = xFlag
	game.Y = yFlag

	screenX, screenY := ebiten.ScreenSizeInFullscreen()
	game.MaxX = screenX - SPRITE_X*game.Size
	game.MaxY = screenY - SPRITE_Y*game.Size

	ebiten.SetWindowSize(SPRITE_X*sizeFlag, SPRITE_Y*sizeFlag)
	ebiten.SetWindowTitle("Shark!")
	ebiten.SetWindowDecorated(false)
	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowFloating(true)

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
