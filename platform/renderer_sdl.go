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

package platform

import (
	"image"
	"image/draw"
	"log"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const fullscreenFlag = sdl.WINDOW_FULLSCREEN //sdl.WINDOW_FULLSCREEN_DESKTOP

type Config func(*sdlRenderer) error

func ConfigWithSize(w, h int) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.windowSize = image.Point{w, h}
		return nil
	}
}

func ConfigWithTitle(title string) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.windowTitle = title
		return nil
	}
}

func ConfigWithDiv(n int) Config {
	return func(rnd *sdlRenderer) error {
		rnd.config.resolutionDiv = n
		return nil
	}
}

func ConfigWithFullscreen(rnd *sdlRenderer) error {
	rnd.config.fullscreen = true
	return nil
}

func ConfigWithDebug(rnd *sdlRenderer) error {
	rnd.config.debug = true
	return nil
}

func ConfigWithNoVSync(rnd *sdlRenderer) error {
	rnd.config.novsync = true
	return nil
}

type sdlRenderer struct {
	window           *sdl.Window
	backBuffer       *image.RGBA
	hwBuffer         *sdl.Texture
	internalRenderer *sdl.Renderer

	config struct {
		windowTitle   string
		windowSize    image.Point
		resolutionDiv int
		debug, novsync,
		fullscreen bool
	}
}

func NewRenderer(configs ...Config) (*sdlRenderer, error) {
	var (
		err error
		r   sdlRenderer
		dm  sdl.DisplayMode

		sdlFlags uint32 = sdl.WINDOW_SHOWN
	)

	for _, cfg := range configs {
		if err = cfg(&r); err != nil {
			return nil, err
		}
	}

	cfg := &r.config
	if cfg.fullscreen {
		sdlFlags |= fullscreenFlag
	}

	if err = sdl.GetDesktopDisplayMode(0, &dm); err != nil {
		return nil, err
	}

	if cfg.windowSize.X <= 0 {
		cfg.windowSize.X = int(dm.W)
	}
	if cfg.windowSize.Y <= 0 {
		cfg.windowSize.Y = int(dm.H)
	}

	if cfg.resolutionDiv > 0 {
		cfg.windowSize.X /= cfg.resolutionDiv
		cfg.windowSize.Y /= cfg.resolutionDiv
	}

	r.window, err = sdl.CreateWindow(cfg.windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, cfg.windowSize.X, cfg.windowSize.Y, sdlFlags)
	if err != nil {
		return nil, err
	}

	const width, height = 320, 200
	r.backBuffer = image.NewRGBA(image.Rect(0, 0, width, height))

	renderer, err := sdl.CreateRenderer(r.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	renderer.SetLogicalSize(width, height)
	r.internalRenderer = renderer

	r.hwBuffer, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		return nil, err
	}

	sdl.ShowCursor(0)
	return &r, nil
}

func (r *sdlRenderer) ToggleFullscreen() {
	isFullscreen := (r.window.GetFlags() & fullscreenFlag) != 0
	if isFullscreen {
		r.window.SetFullscreen(0)
	} else {
		r.window.SetFullscreen(fullscreenFlag)
	}
}

func (r *sdlRenderer) BackBuffer() draw.Image {
	return r.backBuffer
}

func (r *sdlRenderer) Clear() {
	pix := r.backBuffer.Pix
	for i := range pix {
		pix[i] = 0
	}
}

func (r *sdlRenderer) Present() {
	var (
		p     unsafe.Pointer
		pitch int
	)

	if err := r.hwBuffer.Lock(nil, &p, &pitch); err != nil {
		log.Panicln(err)
	}

	sz := len(r.backBuffer.Pix)
	h := &reflect.SliceHeader{Data: uintptr(p), Len: sz, Cap: sz}
	dest := *(*[]byte)(unsafe.Pointer(h))
	copy(dest, r.backBuffer.Pix)

	r.hwBuffer.Unlock()
	r.internalRenderer.Copy(r.hwBuffer, nil, nil)
	r.internalRenderer.Present()
}

func (r *sdlRenderer) Shutdown() {
	r.window.Destroy()
	r.hwBuffer.Destroy()
	r.internalRenderer.Destroy()
}

func (r *sdlRenderer) SetWindowTitle(title string) {
	r.window.SetTitle(title)
}
