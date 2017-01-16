// +build !mobile

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

package world

import (
	"encoding/json"
	"image"
	"image/color"
	"os"
	"path"

	"github.com/andreas-jonsson/go-wolf/engine"
)

type palettedUniform struct {
	rect image.Rectangle
	c    int
}

func (u palettedUniform) Bounds() image.Rectangle {
	return u.rect
}

func (u palettedUniform) ColorIndexAt(x, y int) uint8 {
	return uint8(u.c)
}

type World struct {
	mapData  [][]int
	textures []palettedUniform
}

func NewWorld(name string, pal color.Palette) (*World, error) {
	bounds := image.Rect(0, 0, 1, 1)
	textures := []palettedUniform{
		{bounds, pal.Index(color.RGBA{255, 0, 0, 255})},
		{bounds, pal.Index(color.RGBA{255, 255, 0, 255})},
		{bounds, pal.Index(color.RGBA{255, 0, 255, 255})},
		{bounds, pal.Index(color.RGBA{0, 255, 0, 255})},
		{bounds, pal.Index(color.RGBA{0, 0, 255, 255})},
		{bounds, pal.Index(color.RGBA{127, 255, 0, 255})},
		{bounds, pal.Index(color.RGBA{0, 127, 127, 255})},
		{bounds, pal.Index(color.RGBA{0, 0, 127, 255})},
	}

	w := &World{textures: textures}

	fp, err := os.Open(path.Join("data", "maps", name+".json"))
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&w.mapData); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *World) GetTexture(index, shade int) engine.Texture {
	// Ignore shades for now...
	return &w.textures[index]
}

func (w *World) GetTile(x, y int) int {
	return w.mapData[x][y]
}
