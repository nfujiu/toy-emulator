package chip8

type Cpu struct {
		v [16]uint8
		i uint16
		sp uint8
		pc uint16
		stack [16]uint16
}

func (c *Cpu) Increment() {
		c.pc += 2
}

func (c *Cpu) Tick(ram Ram) {
		println(c)
}