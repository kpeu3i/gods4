package gods4

type Event string

const (
	// Cross
	EventCrossPress   Event = "cross.press"
	EventCrossRelease Event = "cross.release"

	// Circle
	EventCirclePress   Event = "circle.press"
	EventCircleRelease Event = "circle.release"

	// Square
	EventSquarePress   Event = "square.press"
	EventSquareRelease Event = "square.release"

	// Triangle
	EventTrianglePress   Event = "triangle.press"
	EventTriangleRelease Event = "triangle.release"

	// L1
	EventL1Press   Event = "l1.press"
	EventL1Release Event = "l1.release"

	// L2
	EventL2Press   Event = "l2.press"
	EventL2Release Event = "l2.release"

	// L3
	EventL3Press   Event = "l3.press"
	EventL3Release Event = "l3.release"

	// R1
	EventR1Press   Event = "r1.press"
	EventR1Release Event = "r1.release"

	// R2
	EventR2Press   Event = "r2.press"
	EventR2Release Event = "r2.release"

	// R3
	EventR3Press   Event = "r3.press"
	EventR3Release Event = "r3.release"

	// D-pad up
	EventDPadUpPress   Event = "dpad_up.press"
	EventDPadUpRelease Event = "dpad_up.release"

	// D-pad down
	EventDPadDownPress   Event = "dpad_down.press"
	EventDPadDownRelease Event = "dpad_down.release"

	// D-pad left
	EventDPadLeftPress   Event = "dpad_left.press"
	EventDPadLeftRelease Event = "dpad_left.release"

	// D-pad right
	EventDPadRightPress   Event = "dpad_right.press"
	EventDPadRightRelease Event = "dpad_right.release"

	// Share
	EventSharePress   Event = "share.press"
	EventShareRelease Event = "share.release"

	// Options
	EventOptionsPress   Event = "options.press"
	EventOptionsRelease Event = "options.release"

	// Touchpad
	EventTouchpadSwipe   Event = "touchpad.swipe"
	EventTouchpadPress   Event = "touchpad.press"
	EventTouchpadRelease Event = "touchpad.release"

	// PS
	EventPSPress   Event = "ps.press"
	EventPSRelease Event = "ps.release"

	// Left stick
	EventLeftStickMove Event = "left_stick.move"

	// Right stick
	EventRightStickMove Event = "right_stick.move"

	// Accelerometer
	EventAccelerometerUpdate Event = "accelerometer.update"

	// Gyroscope
	EventGyroscopeUpdate Event = "gyroscope.update"

	// Battery
	EventBatteryUpdate Event = "battery.update"
)
