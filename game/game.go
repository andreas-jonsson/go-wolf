// +-------------------=D=r=i=v=e=-=E=n=g=i=n=e=---------------------+
// | Copyright (C) 2016-2017 Andreas T Jonsson. All rights reserved. |
// | Contact <mail@andreasjonsson.se>                                |
// +-----------------------------------------------------------------+

package game

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/andreas-jonsson/go-wolf/platform"
)

type (
	GameState interface {
		Name() string
		Enter(from GameState, args ...interface{}) error
		Exit(to GameState) error
		Update(gctl GameControl) error
		Render(backBuffer *image.Paletted) error
	}

	GameControl interface {
		SwitchState(to string, args ...interface{}) error
		CurrentStateName() string
		Timing() (time.Duration, time.Duration, int)
		PollAll()
		PollEvent() platform.Event
		Terminate()
	}
)

var startupTime = time.Now()

type Game struct {
	currentState GameState
	states       map[string]GameState

	t, ft     time.Time
	fps       int
	dt, tick  time.Duration
	numFrames int
	running   bool
}

func NewGame(states map[string]GameState) (*Game, error) {
	return &Game{running: true, states: states, t: time.Now()}, nil
}

func (g *Game) PollAll() {
	for g.PollEvent() != nil {
	}
}

func (g *Game) PollEvent() platform.Event {
	for {
		event := platform.PollEvent()
		if event == nil {
			return nil
		}

		switch t := event.(type) {
		case *platform.QuitEvent:
			g.running = false
		case *platform.KeyDownEvent:
			switch t.Key {
			case platform.KeyEsc:
				g.running = false
				continue
			}
			return event
		default:
			return event
		}
	}
}

func (g *Game) CurrentStateName() string {
	return g.currentState.Name()
}

func (g *Game) SwitchState(to string, args ...interface{}) error {
	newState, ok := g.states[to]
	if !ok {
		return fmt.Errorf("invalid state: %s", to)
	}

	currentState := g.currentState

	if currentState != nil {
		log.Printf("Exiting state: %v", currentState.Name())
		if err := currentState.Exit(newState); err != nil {
			return err
		}
	}

	g.currentState = newState

	log.Printf("Enter state: %v", to)
	if err := newState.Enter(currentState, args...); err != nil {
		return err
	}

	return nil
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Timing() (time.Duration, time.Duration, int) {
	return g.dt, g.tick, g.fps
}

func (g *Game) Terminate() {
	g.running = false
}

func (g *Game) Update() error {
	g.dt = time.Since(g.t)
	g.tick = time.Since(startupTime)
	g.t = time.Now()

	if err := g.currentState.Update(g); err != nil {
		return err
	}

	g.numFrames++
	if time.Since(g.ft) >= time.Second {
		g.fps = g.numFrames
		g.ft = time.Now()
		g.numFrames = 0
	}

	return nil
}

func (g *Game) Render(backBuffer *image.Paletted) error {
	if err := g.currentState.Render(backBuffer); err != nil {
		return err
	}
	return nil
}

func (g *Game) Shutdown() {
}
