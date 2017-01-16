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
	"image/color/palette"
	"log"

	"github.com/andreas-jonsson/go-wolf/engine"
	"github.com/andreas-jonsson/go-wolf/game"
	"github.com/andreas-jonsson/go-wolf/platform"
	"github.com/andreas-jonsson/go-wolf/world"
)

type renderTarget struct {
	bounds     image.Rectangle
	backBuffer *image.Paletted
	depth      []float64
}

func (rt *renderTarget) Bounds() image.Rectangle {
	return rt.bounds
}

func (rt *renderTarget) SetColorIndex(x, y int, c uint8) {
	if rt.backBuffer != nil {
		rt.backBuffer.SetColorIndex(x, y, c)
	}
}

func (rt *renderTarget) SetZ(x int, z float64) {
	if rt.backBuffer != nil {
		rt.depth[x] = z
	}
}

func (rt *renderTarget) setBackBuffer(img *image.Paletted) {
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
	w, err := world.NewWorld("level1", palette.Plan9)
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
	return nil
}

func (s *playState) Exit(to game.GameState) error {
	return nil
}

func (s *playState) Update(gctl game.GameControl) error {
	for event := gctl.PollEvent(); event != nil; event = gctl.PollEvent() {
		switch event.(type) {
		case *platform.KeyDownEvent, *platform.QuitEvent:
			gctl.Terminate()
		}
	}
	return nil
}

func (s *playState) Render(backBuffer *image.Paletted) error {
	if s.rt.backBuffer == nil {
		s.rt.setBackBuffer(backBuffer)
	}
	s.rc.Render()
	return nil
}
