package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var screenW, screenH = ebiten.ScreenSizeInFullscreen()

type StateMachine struct {
	state          State
	anim           *Anim
	animFrameCount int // number of frames in current anim
	frameIdx       int // frame index in current animation
	ticks          int // ticks since last animation frame change
	lastFed        time.Time
	x, y           int // logical window position

	isXYInit bool
}

func NewStateMachine() *StateMachine {
	sm := StateMachine{}
	sm.SetState(&StateIdle{})
	sm.lastFed = time.Now()
	return &sm
}

func (sm *StateMachine) SetAnim(anim *Anim) {
	sm.anim = anim
	sm.animFrameCount = len(anim.Frames)
	sm.frameIdx = 0
	sm.ticks = 0
}
func (sm *StateMachine) Frame() *ebiten.Image {
	return sm.anim.Frames[sm.frameIdx]
}
func (sm *StateMachine) SetState(s State) {
	sm.state = s
	sm.state.Enter(sm)
}
func (sm *StateMachine) Update() error {
	if !sm.isXYInit {
		sm.x, sm.y = ebiten.WindowPosition()
		sm.isXYInit = true
	}

	sm.state.Update(sm)

	sm.x = min(max(sm.x, 0), screenW-WindowWidth)
	sm.y = min(max(sm.y, 0), screenH-WindowHeight)
	ebiten.SetWindowPosition(sm.x, sm.y)

	// Advance to next animation frame
	sm.ticks += 1
	if sm.ticks < 10 {
		return nil
	}
	sm.ticks = 0
	if sm.frameIdx < sm.animFrameCount-1 {
		sm.frameIdx += 1
		return nil
	}

	// At end of current anim, restart the animation,
	// execute state-specific hook if any.
	sm.frameIdx = 0
	sm.state.EndAnimHook(sm)
	return nil
}

func (sm *StateMachine) Draw(screen *ebiten.Image) {
	screen.DrawImage(sm.Frame(), nil)
}
func (sm *StateMachine) Layout(ow, oh int) (sw, sh int) {
	return SPRITE_X, SPRITE_Y
}

type State interface {
	Enter(sm *StateMachine)
	Update(sm *StateMachine)
	EndAnimHook(sm *StateMachine)
}

type StateIdle struct{}

func (s *StateIdle) Enter(sm *StateMachine) { sm.SetAnim(Idle) }
func (s *StateIdle) Update(sm *StateMachine) {
	if checkHunger(sm) {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		sm.SetState(&StateDrag{})
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		sm.SetState(&StateRClick{})
		return
	}
}
func (s *StateIdle) EndAnimHook(sm *StateMachine) {
	if randBool(WalkChance) {
		sm.SetState(&StateWalk{isLeft: randBool(50)})
	}
}

type StateDrag struct {
	PreviousMousePos Vector
	WinStartPos      Vector
	MouseStartPos    Vector
}

func (s *StateDrag) Enter(sm *StateMachine) {
	sm.SetAnim(Drag)
	s.PreviousMousePos = GlobalCursorPosition()
	s.WinStartPos = CreateVector(ebiten.WindowPosition())
	s.MouseStartPos = GlobalCursorPosition()
}
func (s *StateDrag) Update(sm *StateMachine) {
	if checkHunger(sm) {
		return
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		sm.SetState(&StateIdle{})
		return
	}
	mousePos := GlobalCursorPosition()
	if mousePos != s.PreviousMousePos {
		winPos := s.WinStartPos.Add(mousePos.Subtract(s.MouseStartPos))
		sm.x, sm.y = winPos.x, winPos.y
	}
	s.PreviousMousePos = mousePos
}
func (s *StateDrag) EndAnimHook(sm *StateMachine) {
	syncWinPos(sm)
}

type StateRClick struct{}

func (s *StateRClick) Enter(sm *StateMachine) {
	sm.SetAnim(RightClick)
}
func (s *StateRClick) Update(sm *StateMachine) {
	if checkHunger(sm) {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		sm.SetState(&StateDrag{})
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		sm.SetState(&StateRClick{})
		return
	}
}
func (s *StateRClick) EndAnimHook(sm *StateMachine) {
	sm.SetState(&StateIdle{})
}

type StateHungry struct{}

func (s *StateHungry) Enter(sm *StateMachine) { sm.SetAnim(Hungry) }
func (s *StateHungry) Update(sm *StateMachine) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		sm.SetState(&StateFeed{})
		return
	}
}
func (s *StateHungry) EndAnimHook(sm *StateMachine) {}

func checkHunger(sm *StateMachine) (isHungry bool) {
	now := time.Now()
	if now.Sub(sm.lastFed) >= DurationTillHungry {
		sm.SetState(&StateHungry{})
		return true
	}
	return false
}

type StateFeed struct{}

func (s *StateFeed) Enter(sm *StateMachine)  { sm.SetAnim(Feed) }
func (s *StateFeed) Update(sm *StateMachine) {}
func (s *StateFeed) EndAnimHook(sm *StateMachine) {
	sm.SetState(&StateIdle{})
	sm.lastFed = time.Now()
}

type StateWalk struct {
	isLeft bool
}

func (s *StateWalk) Enter(sm *StateMachine) {
	if s.isLeft {
		sm.SetAnim(WalkLeft)
	} else {
		sm.SetAnim(WalkRight)
	}
}
func (s *StateWalk) Update(sm *StateMachine) {
	if checkHunger(sm) {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		sm.SetState(&StateDrag{})
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		sm.SetState(&StateRClick{})
		return
	}

	if s.isLeft {
		sm.x -= 1
	} else {
		sm.x += 1
	}
}
func (s *StateWalk) EndAnimHook(sm *StateMachine) {
	// On KDE, our window isn't allowed to walk over a vertical taskbar
	// ("panel" in KDE parlance), so when gura walks agains a taskbar, the
	// in-game position (sm.x) will get out of sync with the actual window
	// position. Therefore, we must run a sync:
	syncWinPos(sm)

	if randBool(StopChance) {
		sm.SetState(&StateIdle{})
	}
}

// Sometimes the sm.x/y doesn't match the window's actual position on screen.
// Run this at the end of States that might end up in that situation.
func syncWinPos(sm *StateMachine) {
	actualX, actualY := ebiten.WindowPosition()
	if actualX != sm.x || actualY != sm.y {
		//fmt.Printf("ingame: %d,%d - actual: %d,%d\n", sm.x, sm.y, actualX, actualY)
		sm.x = actualX
		sm.y = actualY
	}
}
