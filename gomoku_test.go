package main

import (
	"testing"
)

func test(board, sample *Board, x, y, color int, t *testing.T) {
	freeThrees := [2]Board {}
	checkDoubleThree(board, &freeThrees[(color + 1) / 2], x, y, color)
	if (freeThrees[(color + 1) / 2] != *sample) {
		t.Error("Expected ", sample, ", got ", freeThrees[1])
	}
}

func TestFreeThrees0(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}

func TestFreeThrees1(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 3, 0, 1, t)
}

func TestFreeThrees2(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}

func TestFreeThrees3(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}

func TestFreeThrees4(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}

func TestFreeThrees5(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}

func TestFreeThrees6(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<2}}, 1, 0, 1, t)
}

func TestFreeThrees7(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<3}}, 1, 0, 1, t)
}

func TestFreeThrees8(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}

func TestFreeThrees9(t *testing.T) {
	test(&Board{{0, 1, 0, 1}}, &Board{{0, 0, 1<<1, 0, 1<<1}}, 1, 0, 1, t)
}
