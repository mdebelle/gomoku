package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

type TextDrawer struct {
	fonts		map[int]*ttf.Font
	color		sdl.Color
}

func NewTextDrawer() *TextDrawer {
	fontMap := make(map[int]*ttf.Font)
	return &TextDrawer{fontMap, sdl.Color{0, 0, 0, 255}}
}

func (this *TextDrawer) GetFont(size int) *ttf.Font {
	font := this.fonts[size]
	if font == nil {
		newFont, err := ttf.OpenFont("/Library/Fonts/Arial Black.ttf", size)
		if err != nil {
			panic(err)
		}
		this.fonts[size] = newFont
		font = newFont
	}
	return font
}

func (this *TextDrawer) Draw(renderer *sdl.Renderer, test string, w, h, size int) {
	font := this.GetFont(size)

	surface, err := font.RenderUTF8_Blended(test, this.color)
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
	for _, font := range this.fonts {
		font.Close()
	}
}
