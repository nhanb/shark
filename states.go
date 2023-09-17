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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		sm.SetState(&StateRClick{})
		return
	}
}
func (s *StateIdle) EndAnimHook(sm *StateMachine) {}

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
func (s *StateDrag) EndAnimHook(sm *StateMachine) {}

type StateRClick struct{}

func (s *StateRClick) Enter(sm *StateMachine) {
	sm.SetAnim(RightClick)
}
func (s *StateRClick) Update(sm *StateMachine) {}
func (s *StateRClick) EndAnimHook(sm *StateMachine) {
	sm.SetState(&StateIdle{})
}
