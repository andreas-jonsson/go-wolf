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

package menu

import (
	"image/draw"

	"github.com/andreas-jonsson/go-wolf/game"
)

type menuState struct {
}

func NewMenuState() *menuState {
	return &menuState{}
}

func (s *menuState) Name() string {
	return "menu"
}

func (s *menuState) Enter(from game.GameState, args ...interface{}) error {
	return args[0].(game.GameControl).SwitchState("play", args[0])
}

func (s *menuState) Exit(to game.GameState) error {
	return nil
}

func (s *menuState) Update(gctl game.GameControl) error {
	gctl.PollAll()
	return nil
}

func (s *menuState) Render(backBuffer draw.Image) error {
	return nil
}
