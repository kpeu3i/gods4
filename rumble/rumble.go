package rumble

type Rumble struct {
	left  byte
	right byte
}

func (r *Rumble) Left() byte {
	return r.left
}

func (r *Rumble) Right() byte {
	return r.right
}

func New(left, right byte) *Rumble {
	return &Rumble{left: left, right: right}
}

func Left() *Rumble {
	return &Rumble{left: 255}
}

func Right() *Rumble {
	return &Rumble{right: 255}
}

func Both() *Rumble {
	return &Rumble{left: 255, right: 255}
}
