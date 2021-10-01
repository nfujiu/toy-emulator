package chip8

type Display struct {
		pixel [64][32]bool
}

func (d *Display) Render() [64][32]bool {
		return d.pixel
}

func (d *Display) Draw(x uint8, y uint8, n uint8, sprite [256]uint8)  (collision bool) {
		for iy := uint8(0); iy < n; iy++ {
				var and uint8 = 0b10000000
				for ix := uint8(0); ix < 8; ix++ {
						thisX := (x + ix) % 64
						thisY := (y + iy) % 32
						spriteOn := (sprite[iy] & and) == and
						displayOn := d.pixel[thisX][thisY]
						value := spriteOn != displayOn
						if !value && displayOn {
								collision = true
						}
						d.pixel[thisX][thisY] = value
						and = and >> 1
				}
		}
		return collision
}

func (d *Display) Clear() {
		for x := 0; x < 64; x++ {
				for y := 0; y < 32; y++ {
						d.pixel[x][y] = false
				}
		}
}