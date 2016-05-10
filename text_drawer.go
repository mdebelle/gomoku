package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type TextDrawer struct {
	font		*ttf.Font
	color		sdl.Color
}

func NewTextDrawer() *TextDrawer {
	font, err := ttf.OpenFont("/Library/Fonts/Arial Black.ttf", 20)
	if err != nil {
		panic(err)
	}
	return &TextDrawer{font, sdl.Color{0, 0, 0, 255}}
}

// strconv.Itoa(values[i][j][0])
// sdl.Rect{int32(i) * 40 + 40 - surface.W / 2, int32(j) * 40 + 40 - surface.H / 2, surface.W, surface.H}

func (this *TextDrawer) Draw(renderer *sdl.Renderer, test string, w, h int) {
	surface, err := this.font.RenderUTF8_Blended(test, this.color)
	if err != nil { panic(err) }

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil { panic(err) }

	surface.Free()
	rect := sdl.Rect {int32(w) - surface.W / 2,
		int32(h) - surface.H / 2,
		surface.W,
		surface.H}
	_, _, rect.W, rect.H, _ = tex.Query()
	renderer.Copy(tex, nil, &rect)
	tex.Destroy()
}

func (this *TextDrawer) Dispose() {
	this.font.Close()
}
