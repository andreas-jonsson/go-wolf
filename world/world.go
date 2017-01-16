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
	"image/draw"
	"image/png"
	"os"
	"path"

	"github.com/andreas-jonsson/go-wolf/engine"
)

type World struct {
	mapData  [][]int
	textures []*image.Paletted
}

func NewWorld(name string, pal color.Palette) (*World, error) {
	w := new(World)
	if err := w.loadTextures(pal); err != nil {
		return nil, err
	}

	if err := w.loadMap(name); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *World) loadTextures(pal color.Palette) error {
	fp, err := os.Open(path.Join("data", "textures", "textures.json"))
	if err != nil {
		return err
	}
	defer fp.Close()

	var textureList []string

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&textureList); err != nil {
		return err
	}

	for _, s := range textureList {
		fp, err := os.Open(path.Join("data", "textures", s))
		if err != nil {
			return err
		}

		img, err := png.Decode(fp)
		fp.Close()

		if err != nil {
			return err
		}

		bounds := img.Bounds()
		convertedImage := image.NewPaletted(bounds, pal)
		draw.Draw(convertedImage, bounds, img, image.ZP, draw.Src)

		w.textures = append(w.textures, convertedImage)
	}

	return nil
}

func (w *World) loadMap(name string) error {
	fp, err := os.Open(path.Join("data", "maps", name+".json"))
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&w.mapData); err != nil {
		return err
	}

	return nil
}

func (w *World) GetTexture(index, shade int) engine.Texture {
	// Ignore shades for now...
	return w.textures[index]
}

func (w *World) GetTile(x, y int) int {
	return w.mapData[x][y]
}
