package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func drawGrid(renderer *sdl.Renderer) {
	_ = renderer.SetDrawColor(236, 0, 0, 0)
	for i := 1; i < 20; i++ {
		_ = renderer.DrawLine(40, 40 * i, 40 * 19, 40 * i)
		_ = renderer.DrawLine(40 *i, 40, 40 * i, 40 * 19)
	}	
}

func drawClic(renderer *sdl.Renderer, values *Board, capture *[3]int) {
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

func draweval(renderer *sdl.Renderer, values *[19][19][5]int, freeThrees *Board) {

	dr := func (x, y, lenx, leny int, vertical bool) {
		_ = renderer.SetDrawColor(0, 0, 0, 0)
		if (!vertical) {
			_ = renderer.DrawLine(x, y - 2, x + lenx, y + leny - 2)
			_ = renderer.DrawLine(x, y - 1, x + lenx, y + leny - 1)
			_ = renderer.DrawLine(x, y, x + lenx, y)
			_ = renderer.DrawLine(x, y - 2, x + lenx, y + leny + 1)
			_ = renderer.DrawLine(x, y - 2, x + lenx, y + leny + 2)
		} else {
			_ = renderer.DrawLine(x, y - 2, x - 2, y + leny)
			_ = renderer.DrawLine(x, y - 1, x - 1, y + leny)
			_ = renderer.DrawLine(x, y, x, y + leny)
			_ = renderer.DrawLine(x, y - 2, x + 1, y + leny)
			_ = renderer.DrawLine(x, y - 2, x + 2, y + leny)
		}
		return
	}

	bitValueAtPosition := func (number, pos int) bool {
		number = number << uint(8 - pos)
		number = number >> uint(7)
		if number == 1 {
			return true
		}
		return false
	}

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

			if freeThrees[j][i] != 0 {
				if bitValueAtPosition(freeThrees[j][i], 1) == true {
						dr((i+1)*40, ((j+1)*40)-20, 0, 40, true)
				}
				if bitValueAtPosition(freeThrees[j][i], 2) == true {
						dr((i+1)*40-20, ((j+1)*40), 40, 0, false)
				}
				if bitValueAtPosition(freeThrees[j][i], 4) == true {
						dr((i+1)*40-10, ((j+1)*40)-10, 20, 20, false)
				}
				if bitValueAtPosition(freeThrees[j][i], 8) == true {
						dr((i+1)*40-10, ((j+1)*40)+10, 20, 20, false)
				}
			}

		}
	}
}
