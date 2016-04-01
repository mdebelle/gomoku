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

func drawClic(renderer *sdl.Renderer, values *[19][19]int, capture *[3]int) {
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

