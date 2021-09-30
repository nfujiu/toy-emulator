package chip8

type Cpu struct {
		v [16]uint8
		i uint16
		sp uint8
		pc uint8
}