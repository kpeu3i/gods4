package gods4

import (
	"encoding/binary"
	"math"
)

const analogSticksSmoothing = 4

type state struct {
	cross         bool
	circle        bool
	square        bool
	triangle      bool
	l1            bool
	l2            byte
	l3            bool
	r1            bool
	r2            byte
	r3            bool
	dPadUp        bool
	dPadDown      bool
	dPadLeft      bool
	dPadRight     bool
	share         bool
	options       bool
	ps            bool
	leftStick     Stick
	rightStick    Stick
	touchpad      Touchpad
	accelerometer Accelerometer
	gyroscope     Gyroscope
	battery       Battery
}

type Stick struct {
	X byte
	Y byte
}

type Touchpad struct {
	Press bool
	Swipe []Touch
}

type Touch struct {
	IsActive bool
	X        byte
	Y        byte
}

type Accelerometer struct {
	X int16
	Y int16
	Z int16
}

type Gyroscope struct {
	Roll  int16
	Yaw   int16
	Pitch int16
}

type Battery struct {
	Capacity         byte
	IsCharging       bool
	IsCableConnected bool
}

func newState(bytes []byte, offset uint, prevState *state) *state {
	s := &state{
		cross:         buttonCrossState(bytes, offset),
		circle:        buttonCircleState(bytes, offset),
		square:        buttonSquareState(bytes, offset),
		triangle:      buttonTriangleState(bytes, offset),
		l1:            buttonL1State(bytes, offset),
		l2:            buttonL2State(bytes, offset),
		l3:            buttonL3State(bytes, offset),
		r1:            buttonR1State(bytes, offset),
		r2:            buttonR2State(bytes, offset),
		r3:            buttonR3State(bytes, offset),
		dPadUp:        buttonDPadUpState(bytes, offset),
		dPadDown:      buttonDPadDownState(bytes, offset),
		dPadLeft:      buttonDPadLeftState(bytes, offset),
		dPadRight:     buttonDPadRightState(bytes, offset),
		share:         buttonShareState(bytes, offset),
		options:       buttonOptionsState(bytes, offset),
		ps:            buttonPSState(bytes, offset),
		leftStick:     buttonLeftStickState(bytes, offset, prevState),
		rightStick:    buttonRightStickState(bytes, offset, prevState),
		touchpad:      touchpadState(bytes, offset),
		accelerometer: accelerometerState(bytes, offset),
		gyroscope:     gyroscopeState(bytes, offset),
		battery:       batteryState(bytes, offset),
	}

	return s
}

func buttonCrossState(bytes []byte, offset uint) bool {
	return bytes[5+offset]&32 != 0
}

func buttonCircleState(bytes []byte, offset uint) bool {
	return bytes[5+offset]&64 != 0
}

func buttonSquareState(bytes []byte, offset uint) bool {
	return bytes[5+offset]&16 != 0
}

func buttonTriangleState(bytes []byte, offset uint) bool {
	return bytes[5+offset]&128 != 0
}

func buttonL1State(bytes []byte, offset uint) bool {
	return bytes[6+offset]&1 != 0
}

func buttonL2State(bytes []byte, offset uint) byte {
	if bytes[6+offset]&4 != 0 {
		return bytes[8+offset]
	}

	return 0
}

func buttonL3State(bytes []byte, offset uint) bool {
	return bytes[6+offset]&64 != 0
}

func buttonR1State(bytes []byte, offset uint) bool {
	return bytes[6+offset]&2 != 0
}

func buttonR2State(bytes []byte, offset uint) byte {
	if bytes[6+offset]&8 != 0 {
		return bytes[9+offset]
	}

	return 0
}

func buttonR3State(bytes []byte, offset uint) bool {
	return bytes[6+offset]&128 != 0
}

func buttonDPadUpState(bytes []byte, offset uint) bool {
	v := bytes[5+offset] & 15

	return v == 0 || v == 1 || v == 7
}

func buttonDPadDownState(bytes []byte, offset uint) bool {
	v := bytes[5+offset] & 15

	return v == 3 || v == 4 || v == 5
}

func buttonDPadLeftState(bytes []byte, offset uint) bool {
	v := bytes[5+offset] & 15

	return v == 5 || v == 6 || v == 7
}

func buttonDPadRightState(bytes []byte, offset uint) bool {
	v := bytes[5+offset] & 15

	return v == 1 || v == 2 || v == 3
}

func buttonShareState(bytes []byte, offset uint) bool {
	return bytes[6+offset]&16 != 0
}

func buttonOptionsState(bytes []byte, offset uint) bool {
	return bytes[6+offset]&32 != 0
}

//func buttonTouchpadState(bytes []byte, offset uint) bool {
//	return bytes[7+offset]&2 != 0
//}

func buttonPSState(bytes []byte, offset uint) bool {
	return bytes[7+offset]&1 != 0
}

func buttonLeftStickState(bytes []byte, offset uint, prevState *state) Stick {
	var prevX, prevY byte

	if prevState == nil {
		prevX, prevY = bytes[1+offset], bytes[2+offset]
	} else {
		prevX, prevY = prevState.leftStick.X, prevState.leftStick.Y
	}

	if math.Abs(float64(bytes[1+offset])-float64(prevX)) >= float64(analogSticksSmoothing) ||
		math.Abs(float64(bytes[2+offset])-float64(prevY)) >= float64(analogSticksSmoothing) {
		return Stick{X: bytes[1+offset], Y: bytes[2+offset]}
	}

	return Stick{X: prevX, Y: prevY}
}

func buttonRightStickState(bytes []byte, offset uint, prevState *state) Stick {
	var prevX, prevY byte

	if prevState == nil {
		prevX, prevY = bytes[3+offset], bytes[4+offset]
	} else {
		prevX, prevY = prevState.rightStick.X, prevState.rightStick.Y
	}

	if math.Abs(float64(bytes[3+offset])-float64(prevX)) >= float64(analogSticksSmoothing) ||
		math.Abs(float64(bytes[4+offset])-float64(prevY)) >= float64(analogSticksSmoothing) {
		return Stick{X: bytes[3+offset], Y: bytes[4+offset]}
	}

	return Stick{X: prevX, Y: prevY}
}

func touchpadState(bytes []byte, offset uint) Touchpad {
	var (
		touches     []Touch
		touchOffset uint
	)

	for i := 1; i <= 2; i++ {
		touch := Touch{
			IsActive: (bytes[35+touchOffset+offset] >> 7) == 0,
			X:        ((bytes[37+touchOffset+offset] & 0x0F) << 8) | bytes[36+touchOffset+offset],
			Y:        bytes[38+touchOffset+offset]<<4 | ((bytes[37+touchOffset+offset] & 0xF0) >> 4),
		}

		touches = append(touches, touch)
		touchOffset += 4
	}

	t := Touchpad{
		Press: bytes[7+offset]&2 != 0,
		Swipe: touches,
	}

	return t
}

func accelerometerState(bytes []byte, offset uint) Accelerometer {
	a := Accelerometer{
		X: int16(binary.LittleEndian.Uint16(bytes[13+offset:])),
		Y: -int16(binary.LittleEndian.Uint16(bytes[15+offset:])),
		Z: -int16(binary.LittleEndian.Uint16(bytes[17+offset:])),
	}

	return a
}

func gyroscopeState(bytes []byte, offset uint) Gyroscope {
	g := Gyroscope{
		Roll:  -int16(binary.LittleEndian.Uint16(bytes[19+offset:])),
		Yaw:   int16(binary.LittleEndian.Uint16(bytes[21+offset:])),
		Pitch: int16(binary.LittleEndian.Uint16(bytes[23+offset:])),
	}

	return g
}

func batteryState(bytes []byte, offset uint) Battery {
	var (
		isCharging  bool
		maxCapacity byte
	)

	capacity := bytes[30+offset] & 0x0F
	isCableConnected := ((bytes[30+offset] >> 4) & 0x01) == 1

	if !isCableConnected || capacity > 10 {
		isCharging = false
	} else {
		isCharging = true
	}

	if isCableConnected {
		maxCapacity = 10
	} else {
		maxCapacity = 9
	}

	if capacity > maxCapacity {
		capacity = maxCapacity
	}

	capacity = byte(float64(float64(capacity) / float64(maxCapacity) * 100))

	battery := Battery{
		Capacity:         capacity,
		IsCharging:       isCharging,
		IsCableConnected: isCableConnected,
	}

	return battery
}
