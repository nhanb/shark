package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type StateMachine struct {
	state          State
	anim           *Anim
	animFrameCount int // number of frames in current anim
	frameIdx       int // frame index in current animation
	ticks          int // ticks since last animation frame change
}

func NewStateMachine() *StateMachine {
	sm := StateMachine{state: StateIdle{}}
	sm.state.Enter(&sm)
	return &sm
}

func (sm *StateMachine) setAnim(anim *Anim) {
	sm.anim = anim
	sm.animFrameCount = len(anim.Frames)
}
func (sm *StateMachine) Frame() *ebiten.Image {
	return sm.anim.Frames[sm.frameIdx]
}

func (sm *StateMachine) HandleInput() {
	sm.state.HandleInput(sm)
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
	HandleInput(sm *StateMachine)
	Update(sm *StateMachine)
}

type StateIdle struct{}
type StateDrag struct{}
type StateRClick struct{}
type StateHungry struct{}
type StateFeed struct{}
type StateWalkL struct{}
type StateWalkR struct{}

func (s StateIdle) Enter(sm *StateMachine)       { sm.setAnim(Idle) }
func (s StateIdle) HandleInput(sm *StateMachine) {}
func (s StateIdle) Update(sm *StateMachine)      {}
