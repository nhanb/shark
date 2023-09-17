package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StateMachine struct {
	state          State
	anim           *Anim
	animFrameCount int // number of frames in current anim
	frameIdx       int // frame index in current animation
	ticks          int // ticks since last animation frame change
}

func NewStateMachine() *StateMachine {
	sm := StateMachine{}
	sm.SetState(&StateIdle{})
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
	sm.state.Update(sm)

	// Advance or reset animation frame
	sm.ticks += 1
	if sm.ticks < 10 {
		return nil
	}
	sm.ticks = 0
	if sm.frameIdx >= sm.animFrameCount-1 {
		sm.frameIdx = 0
	} else {
		sm.frameIdx += 1
	}
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
}

type StateRClick struct{}
type StateHungry struct{}
type StateFeed struct{}
type StateWalkL struct{}
type StateWalkR struct{}

type StateIdle struct{}

func (s *StateIdle) Enter(sm *StateMachine) { sm.SetAnim(Idle) }
func (s *StateIdle) Update(sm *StateMachine) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		sm.SetState(&StateDrag{})
		return
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
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		sm.SetState(&StateIdle{})
		return
	}
	mousePos := GlobalCursorPosition()
	if mousePos != s.PreviousMousePos {
		winPos := s.WinStartPos.Add(mousePos.Subtract(s.MouseStartPos))
		ebiten.SetWindowPosition(winPos.x, winPos.y)
	}
	s.PreviousMousePos = mousePos
}
