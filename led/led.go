package led

type Led struct {
	red      byte
	green    byte
	blue     byte
	flashOn  byte
	flashOff byte
}

func (l *Led) Red() byte {
	return l.red
}

func (l *Led) Green() byte {
	return l.green
}

func (l *Led) Blue() byte {
	return l.blue
}

func (l *Led) FlashOn() byte {
	return l.flashOn
}

func (l *Led) FlashOff() byte {
	return l.flashOff
}

func (l *Led) Flash(on, off byte) *Led {
	l.flashOn = on
	l.flashOff = off

	return l
}

func RGB(red, green, blue byte) *Led {
	return &Led{red: red, green: green, blue: blue}
}

func None() *Led {
	return &Led{}
}

func White() *Led {
	return &Led{red: 255, green: 255, blue: 255}
}

func Red() *Led {
	return &Led{red: 255}
}

func Green() *Led {
	return &Led{green: 128}
}

func Blue() *Led {
	return &Led{blue: 255}
}

func Lime() *Led {
	return &Led{green: 255}
}

func Yellow() *Led {
	return &Led{red: 255, green: 255}
}

func Cyan() *Led {
	return &Led{green: 255, blue: 255}
}

func Magenta() *Led {
	return &Led{red: 255, blue: 255}
}

func Silver() *Led {
	return &Led{red: 192, green: 192, blue: 192}
}

func Gray() *Led {
	return &Led{red: 128, green: 128, blue: 128}
}

func Maroon() *Led {
	return &Led{red: 128}
}

func Olive() *Led {
	return &Led{red: 128, green: 128}
}

func Purple() *Led {
	return &Led{red: 128, blue: 128}
}

func Teal() *Led {
	return &Led{green: 128, blue: 128}
}

func Navy() *Led {
	return &Led{blue: 128}
}
