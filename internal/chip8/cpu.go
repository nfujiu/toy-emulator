package chip8

type Cpu struct {
		v [16]unit8
		i uint16
		sp unit8
		pc unit8
}