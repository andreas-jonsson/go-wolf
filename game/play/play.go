/*
Copyright (C) 2017 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package play

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/andreas-jonsson/go-wolf/engine"
	"github.com/andreas-jonsson/go-wolf/game"
	"github.com/andreas-jonsson/go-wolf/platform"
	"github.com/andreas-jonsson/go-wolf/world"
	"github.com/ungerik/go3d/float64/vec2"
)

type renderTarget struct {
	bounds     image.Rectangle
	backBuffer draw.Image
	depth      []float64
}

func (rt *renderTarget) Bounds() image.Rectangle {
	return rt.bounds
}

func (rt *renderTarget) Set(x, y int, c color.Color) {
	if rt.backBuffer != nil {
		rt.backBuffer.Set(x, y, c)
	}
}

func (rt *renderTarget) SetZ(x int, z float64) {
	if rt.backBuffer != nil {
		rt.depth[x] = z
	}
}

func (rt *renderTarget) setBackBuffer(img draw.Image) {
	rt.bounds = img.Bounds()
	rt.backBuffer = img
	rt.depth = make([]float64, rt.Bounds().Size().X)
}

type playState struct {
	rt *renderTarget
	rc *engine.Raycaster
}

func NewPlayState() *playState {
	rt := new(renderTarget)
	w, err := world.NewWorld("level1")
	if err != nil {
		log.Panicln(err)
	}

	return &playState{
		rt: rt,
		rc: engine.NewRaycaster(rt, w),
	}
}

func (s *playState) Name() string {
	return "play"
}

func (s *playState) Enter(from game.GameState, args ...interface{}) error {
	s.rc.Move(vec2.T{22.0, 11.5})
	return nil
}

func (s *playState) Exit(to game.GameState) error {
	return nil
}

func (s *playState) Update(gctl game.GameControl) error {
	rc := s.rc
	dt, _, _ := gctl.Timing()
	dtf := dt.Seconds()

	const (
		moveSpeed = 10.0
		rotSpeed  = 7.5
	)

	for event := gctl.PollEvent(); event != nil; event = gctl.PollEvent() {
		switch t := event.(type) {
		case *platform.QuitEvent:
			gctl.Terminate()
		case *platform.KeyDownEvent:
			switch t.Key {
			case platform.KeyUp:
				v := rc.Dir()
				v.Scale(moveSpeed * dtf)
				rc.Move(v)
			case platform.KeyDown:
				v := rc.Dir()
				v.Scale(-moveSpeed * dtf)
				rc.Move(v)
			case platform.KeyLeft:
				rc.Rotate(rotSpeed * dtf)
			case platform.KeyRight:
				rc.Rotate(-rotSpeed * dtf)
			}
		}
	}
	return nil
}

func (s *playState) Render(backBuffer draw.Image) error {
	if s.rt.backBuffer == nil {
		s.rt.setBackBuffer(backBuffer)
	}

	size := backBuffer.Bounds().Size()
	roofColor := color.RGBA{75, 75, 75, 255}
	floorColor := color.RGBA{100, 100, 100, 255}

	for y := 0; y < size.Y; y++ {
		c := roofColor
		if y > size.Y/2 {
			c = floorColor
		}
		for x := 0; x < size.X; x++ {
			backBuffer.Set(x, y, c)
		}
	}

	s.rc.Render()
	return nil
}
