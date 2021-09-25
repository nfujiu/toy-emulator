package main

import (
		"chip8/internal/chip8"
		"github.com/faiface/pixel/pixelgl"
)

func main() {
		pixelgl.Run(chip8.Run)
}