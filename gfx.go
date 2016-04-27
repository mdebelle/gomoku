package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"strconv"
)

func drawGrid(renderer *sdl.Renderer) {
	_ = renderer.SetDrawColor(236, 0, 0, 0)
	for i := 1; i < 20; i++ {
		_ = renderer.DrawLine(40, 40 * i, 40 * 19, 40 * i)
		_ = renderer.DrawLine(40 *i, 40, 40 * i, 40 * 19)
	}	
}

func drawClic(renderer *sdl.Renderer, values *Board, capture *[3]int, freeThrees *Board) {

	dr := func (x, y, lenx, leny int, vertical bool) {
		_ = renderer.SetDrawColor(0, 0, 0, 0)
		if (!vertical) {
			_ = renderer.DrawLine(x, y - 2, x + lenx, y + leny - 2)
			_ = renderer.DrawLine(x, y - 1, x + lenx, y + leny - 1)
			_ = renderer.DrawLine(x, y, x + lenx, y + leny)
			_ = renderer.DrawLine(x, y + 1, x + lenx, y + leny + 1)
			_ = renderer.DrawLine(x, y + 2, x + lenx, y + leny + 2)
		} else {
			_ = renderer.DrawLine(x - 2, y, x - 2, y + leny)
			_ = renderer.DrawLine(x - 1, y, x - 1, y + leny)
			_ = renderer.DrawLine(x, y, x, y + leny)
			_ = renderer.DrawLine(x + 1, y, x + 1, y + leny)
			_ = renderer.DrawLine(x + 2, y, x + 2, y + leny)
		}
		return
	}

	bitValueAtPosition := func (number, pos int) bool {
		if bit := ((number >> uint(pos - 1)) & 1); bit == 1 {
			return true
		}
		return false
	}

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if values[j][i] == player_one {
				_ = renderer.SetDrawColor(0, 236, 0, 0)
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			} else if values[j][i] == player_two {
				_ = renderer.SetDrawColor(0, 0, 236, 0)
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}
			}
					
			if freeThrees[j][i] != 0 {
				if bitValueAtPosition(freeThrees[j][i], 1) == true {
						dr((i+1)*40, ((j+1)*40)-15, 0, 30, true)
				}
				if bitValueAtPosition(freeThrees[j][i], 2) == true {
						dr((i+1)*40-15, ((j+1)*40), 30, 0, false)
				}
				if bitValueAtPosition(freeThrees[j][i], 3) == true {
						dr((i+1)*40-15, ((j+1)*40)+15, 30, -30, false)
				}
				if bitValueAtPosition(freeThrees[j][i], 4) == true {
						dr((i+1)*40-15, ((j+1)*40)-15, 30, 30, false)
				}
			}

		}
	}
	_ = renderer.SetDrawColor(0, 236, 0, 0)
	for i := 0; i < capture[0]; i++ {
		for k := 0; k < 20; k++ {
			_ = renderer.DrawLine(((i+1)*40)-10, 800+(k-10), ((i+1)*40)+10, 800+(k-10))
		}
	}
	_ = renderer.SetDrawColor(0, 0, 236, 0)
	for i := 0; i < capture[2]; i++ {
		for k := 0; k < 20; k++ {
			_ = renderer.DrawLine(((i+1)*40)-10, 840+(k-10), ((i+1)*40)+10, 840+(k-10))
		}
	}


}

func draweval(renderer *sdl.Renderer, values *BoardData) {

	// dr := func (x, y, lenx, leny int, vertical bool) {
	// 	_ = renderer.SetDrawColor(0, 0, 0, 0)
	// 	if (!vertical) {
	// 		_ = renderer.DrawLine(x, y - 2, x + lenx, y + leny - 2)
	// 		_ = renderer.DrawLine(x, y - 1, x + lenx, y + leny - 1)
	// 		_ = renderer.DrawLine(x, y, x + lenx, y + leny)
	// 		_ = renderer.DrawLine(x, y + 1, x + lenx, y + leny + 1)
	// 		_ = renderer.DrawLine(x, y + 2, x + lenx, y + leny + 2)
	// 	} else {
	// 		_ = renderer.DrawLine(x - 2, y, x - 2, y + leny)
	// 		_ = renderer.DrawLine(x - 1, y, x - 1, y + leny)
	// 		_ = renderer.DrawLine(x, y, x, y + leny)
	// 		_ = renderer.DrawLine(x + 1, y, x + 1, y + leny)
	// 		_ = renderer.DrawLine(x + 2, y, x + 2, y + leny)
	// 	}
	// 	return
	// }

	// bitValueAtPosition := func (number, pos int) bool {
	// 	if bit := ((number >> uint(pos - 1)) & 1); bit == 1 {
	// 		return true
	// 	}
	// 	return false
	// }

	font, err := ttf.OpenFont("/Library/Fonts/Arial.ttf", 9)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if values[j][i][0] != 0 {
				switch {
					case values[j][i][0] > 4:
						_ = renderer.SetDrawColor(0, 0, 250, 0)
					case values[j][i][0] > 3:
						_ = renderer.SetDrawColor(50, 50, 250, 0)
					case values[j][i][0] > 2:
						_ = renderer.SetDrawColor(100, 100, 250, 0)
					case values[j][i][0] > 1:
						_ = renderer.SetDrawColor(150, 150, 250, 0)
					case values[j][i][0] > 0:
						_ = renderer.SetDrawColor(200, 200, 250, 0)
				}
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40), ((j+1)*40)+(k-10))
				}

				surface, err := font.RenderUTF8_Blended(strconv.Itoa(values[i][j][0]), sdl.Color {0, 0, 0, 1})
				if err != nil { panic(err) }

				tex, err := renderer.CreateTextureFromSurface(surface)
				if err != nil { panic(err) }

				surface.Free()
				rect := sdl.Rect {int32(i) * 40 + 40 - surface.W / 2, int32(j) * 40 + 40 - surface.H / 2, surface.W, surface.H}
				_, _, rect.W, rect.H, _ = tex.Query()
				renderer.Copy(tex, nil, &rect)
			}
			if values[j][i][1] != 0 {
				switch {
					case values[j][i][1] > 4:
						_ = renderer.SetDrawColor(0, 250, 0, 0)
					case values[j][i][1] > 3:
						_ = renderer.SetDrawColor(50, 250, 50, 0)
					case values[j][i][1] > 2:
						_ = renderer.SetDrawColor(100, 250, 100, 0)
					case values[j][i][1] > 1:
						_ = renderer.SetDrawColor(150, 250, 150, 0)
					case values[j][i][1] > 0:
						_ = renderer.SetDrawColor(200, 250, 200, 0)
				}
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40), ((j+1)*40)+(k-10), ((i+1)*40)+10, ((j+1)*40)+(k-10))
				}

			}
			if values[j][i][2] != 0 {
				_ = renderer.SetDrawColor(250, 0, 0, 0)
				for k := 0; k < 10; k++ {
					_ = renderer.DrawLine(((i+1)*40)-5, ((j+1)*40)+(k-5), ((i+1)*40)+5, ((j+1)*40)+(k-5))
				}
			}

			// if freeThrees[j][i] != 0 {
			// 	if bitValueAtPosition(freeThrees[j][i], 1) == true {
			// 			dr((i+1)*40, ((j+1)*40)-15, 0, 30, true)
			// 	}
			// 	if bitValueAtPosition(freeThrees[j][i], 2) == true {
			// 			dr((i+1)*40-15, ((j+1)*40), 30, 0, false)
			// 	}
			// 	if bitValueAtPosition(freeThrees[j][i], 3) == true {
			// 			dr((i+1)*40-15, ((j+1)*40)+15, 30, -30, false)
			// 	}
			// 	if bitValueAtPosition(freeThrees[j][i], 4) == true {
			// 			dr((i+1)*40-15, ((j+1)*40)-15, 30, 30, false)
			// 	}
			// }

		}
	}
}
