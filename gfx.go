package main

import (
	"github.com/veandco/go-sdl2/sdl"
//	"github.com/veandco/go-sdl2/sdl_ttf"
	"strconv"
)

func drawGrid(renderer *sdl.Renderer) {

	_ = renderer.SetDrawColor(44, 62, 80, 255)
	_ = renderer.FillRect(&sdl.Rect{0, 0, 800, 880})

	_ = renderer.SetDrawColor(149, 165, 166, 255)
	for i := 1; i < 20; i++ {
		_ = renderer.DrawLine(40, 40 * i, 40 * 19, 40 * i)
		_ = renderer.DrawLine(40 *i, 40, 40 * i, 40 * 19)
	}	
}

func drawPanel(renderer *sdl.Renderer) {
	
	renderer.Clear()
	_ = renderer.SetDrawColor(44, 62, 80, 255)
	_ = renderer.FillRect(&sdl.Rect{0, 0, 800, 880})
	_ = renderer.SetDrawColor(149, 165, 166, 255)
	_ = renderer.FillRect(&sdl.Rect{10, 400, 780, 80})
	textDrawer.Draw(renderer, "Payer vs Computer Press 1", 400, 420)
	textDrawer.Draw(renderer, "Payer vs Player Press 2", 400, 450)
	renderer.Present()

}

func drawRestartPanel(renderer *sdl.Renderer) {
	
	_ = renderer.SetDrawColor(149, 165, 166, 255)
	_ = renderer.FillRect(&sdl.Rect{10, 400, 780, 80})
	textDrawer.Draw(renderer, "Start Again Press a", 400, 420)
	textDrawer.Draw(renderer, "Quit Press q", 400, 450)


}


func drawClic(renderer *sdl.Renderer, values *Board, capture *[3]int, lstCaptures []Position, freeThrees *[2]Board, searchTime int) {

	dr := func (x, y, lenx, leny int, vertical bool) {
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

	drawOctogone := func (i, j int) {
		_ = renderer.DrawLine((i+1)*40 - 5, (j+1)*40 - 10, (i+1)*40 + 5, (j+1)*40 - 10)
		_ = renderer.DrawLine((i+1)*40 - 6, (j+1)*40 - 9, (i+1)*40 + 6, (j+1)*40 - 9)
		_ = renderer.DrawLine((i+1)*40 - 7, (j+1)*40 - 8, (i+1)*40 + 7, (j+1)*40 - 8)
		_ = renderer.DrawLine((i+1)*40 - 8, (j+1)*40 - 7, (i+1)*40 + 8, (j+1)*40 - 7)
		_ = renderer.DrawLine((i+1)*40 - 9, (j+1)*40 - 6, (i+1)*40 + 9, (j+1)*40 - 6)
		_ = renderer.FillRect(&sdl.Rect{int32((i+1)*40 - 10), int32((j+1)*40 - 5), 20, 11})
		_ = renderer.DrawLine((i+1)*40 - 9, (j+1)*40 + 6, (i+1)*40 + 9, (j+1)*40 + 6)
		_ = renderer.DrawLine((i+1)*40 - 8, (j+1)*40 + 7, (i+1)*40 + 8, (j+1)*40 + 7)
		_ = renderer.DrawLine((i+1)*40 - 7, (j+1)*40 + 8, (i+1)*40 + 7, (j+1)*40 + 8)
		_ = renderer.DrawLine((i+1)*40 - 6, (j+1)*40 + 9, (i+1)*40 + 6, (j+1)*40 + 9)
		_ = renderer.DrawLine((i+1)*40 - 5, (j+1)*40 + 10, (i+1)*40 + 5, (j+1)*40 + 10)
	}

	drawDoubleFree := func (i, j, player int) {
		freeThrees := &freeThrees[(player + 1) / 2]
		if doesDoubleFreeThreePlayer(freeThrees, i, j) {
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

	drawGrid(renderer)
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if values[j][i] == player_one {
				_ = renderer.SetDrawColor(231, 76, 60, 255)
				drawOctogone(i, j)

			} else if values[j][i] == player_two {
				_ = renderer.SetDrawColor(52, 152, 219, 255)
				drawOctogone(i, j)
			}

			_ = renderer.SetDrawColor(220, 20, 60, 255)
			drawDoubleFree(i, j, player_one)
			_ = renderer.SetDrawColor(21, 96, 189, 255)
			drawDoubleFree(i, j, player_two)
		}
	}


	textDrawer.Draw(renderer, strconv.Itoa(capture[0]), 484, 790)
	_ = renderer.SetDrawColor(231, 76, 60, 255)
	for i := 0; i < capture[0]; i++ {
		_ = renderer.FillRect(&sdl.Rect{int32((i+1)*40 - 10), int32(800 - 10), 20, 20})
	}
	textDrawer.Draw(renderer, strconv.Itoa(capture[2]), 484, 830)
	_ = renderer.SetDrawColor(52, 152, 219, 255)
	for i := 0; i < capture[2]; i++ {
		_ = renderer.FillRect(&sdl.Rect{int32((i+1)*40 - 10), int32(840 - 10), 20, 20})
	}

	// TODO: Display forced captures
	/*
	if (victory.Todo == true) {
		
		_ = renderer.FillRect(&sdl.Rect{int32((victory.X+1)*40 - 10), int32((victory.Y+1)*40 - 10), 20, 20})
	}
	*/
	_ = renderer.SetDrawColor(220, 32, 220, 255)
	for i := 0; i < len(lstCaptures); i++ {
		_ = renderer.FillRect(&sdl.Rect{int32((lstCaptures[i].x+1)*40 - 10), int32((lstCaptures[i].y+1)*40 - 10), 20, 20})
	}

	t := searchTime / 1000000
	textDrawer.Draw(renderer, strconv.Itoa(t), 400, 20)

}

func draweval(renderer *sdl.Renderer, values *BoardData) {

	var alpha uint8

	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			if values[j][i][0] != 0 {
				switch {
					case values[j][i][0] > 4:
						alpha = 240
					case values[j][i][0] > 3:
						alpha = 180
					case values[j][i][0] > 2:
						alpha = 120
					case values[j][i][0] > 1:
						alpha = 60
					case values[j][i][0] > 0:
						alpha = 0
				}
				_ = renderer.SetDrawColor(52, 152, 219, alpha)
				_ = renderer.FillRect(&sdl.Rect{int32((i+1)*40 - 10), int32((j+1)*40 - 10), 20, 20})
				for k := 0; k < 20; k++ {
					_ = renderer.DrawLine(((i+1)*40)-10, ((j+1)*40)+(k-10), ((i+1)*40), ((j+1)*40)+(k-10))
				}
			}
			if values[j][i][1] != 0 {
				switch {
					case values[j][i][1] > 4:
						alpha = 240
					case values[j][i][1] > 3:
						alpha = 180
					case values[j][i][1] > 2:
						alpha = 120
					case values[j][i][1] > 1:
						alpha = 60
					case values[j][i][1] > 0:
						alpha = 0
				}
				_ = renderer.SetDrawColor(231, 76, 60, alpha)
				_ = renderer.FillRect(&sdl.Rect{int32((i+1)*40 - 10), int32((j+1)*40 - 10), 20, 20})

			}
			if values[j][i][2] != 0 {
				_ = renderer.SetDrawColor(46, 204, 113, 255)
				_ = renderer.FillRect(&sdl.Rect{int32((i+1)*40 - 5), int32((j+1)*40 - 5), 10, 10})
			}

			if values[j][i][6] != 0 {
				if (values[j][i][5] > 2e9) {
					textDrawer.Draw(renderer, "∞", (i + 1) * 40, (j + 1) * 40)
				} else if (values[j][i][5] < -2e9) {
					textDrawer.Draw(renderer, "-∞", (i + 1) * 40, (j + 1) * 40)
				} else {
					textDrawer.Draw(renderer, strconv.Itoa(values[j][i][5]), (i + 1) * 40, (j + 1) * 40)
				}
				
				if (values[j][i][7] > 2e9) {
					textDrawer.Draw(renderer, "∞", (i + 1) * 40 + 10, (j + 1) * 40 + 10)
				} else if (values[j][i][7] < -2e9) {
					textDrawer.Draw(renderer, "-∞", (i + 1) * 40 + 10, (j + 1) * 40 + 10)
				} else {
					textDrawer.Draw(renderer, strconv.Itoa(values[j][i][7]), (i + 1) * 40 + 10, (j + 1) * 40 + 10)
				}
			}
		}
	}
}
