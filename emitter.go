package gods4

import (
	"sync"
)

type emitter struct {
	mutex     sync.RWMutex
	callbacks map[Event]Callback
	checkers  []func(currState, prevState *state) error
}

func (e *emitter) emit(currState, prevState *state) error {
	for _, checker := range e.checkers {
		err := checker(currState, prevState)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *emitter) callback(event Event) (Callback, bool) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	if callback, ok := e.callbacks[event]; ok {
		return callback, true
	}

	return nil, false
}

func (e *emitter) setCallback(event Event, fn Callback) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.callbacks[event] = fn
}

func (e *emitter) unsetCallback(event Event) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.callbacks, event)
}

func (e *emitter) checkBattery(currState, prevState *state) error {
	if currState.battery.Capacity != prevState.battery.Capacity ||
		currState.battery.IsCharging != prevState.battery.IsCharging ||
		currState.battery.IsCableConnected != prevState.battery.IsCableConnected {
		event := EventBatteryUpdate
		if callback, ok := e.callback(event); ok {
			err := callback(currState.battery)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkTouchpad(currState, prevState *state) error {
	isSwipeChanged := false
	if len(currState.touchpad.Swipe) == len(prevState.touchpad.Swipe) {
		for i, currTouch := range currState.touchpad.Swipe {
			prevTouch := prevState.touchpad.Swipe[i]
			if currTouch.IsActive != prevTouch.IsActive ||
				currTouch.X != prevTouch.X ||
				currTouch.Y != prevTouch.Y {
				isSwipeChanged = true
				break
			}
		}
	} else {
		isSwipeChanged = true
	}

	if isSwipeChanged {
		event := EventTouchpadSwipe
		if callback, ok := e.callback(event); ok {
			err := callback(currState.touchpad)
			if err != nil {
				return err
			}
		}
	}

	if currState.touchpad.Press && !prevState.touchpad.Press {
		event := EventTouchpadPress
		if callback, ok := e.callback(event); ok {
			err := callback(currState.touchpad)
			if err != nil {
				return err
			}
		}
	}

	if !currState.touchpad.Press && prevState.touchpad.Press {
		event := EventTouchpadRelease
		if callback, ok := e.callback(event); ok {
			err := callback(currState.touchpad)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkAccelerometer(currState, prevState *state) error {
	if currState.accelerometer.X != prevState.accelerometer.X ||
		currState.accelerometer.Y != prevState.accelerometer.Y ||
		currState.accelerometer.Z != prevState.accelerometer.Z {
		event := EventAccelerometerUpdate
		if callback, ok := e.callback(event); ok {
			err := callback(currState.accelerometer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkGyroscope(currState, prevState *state) error {
	if currState.gyroscope.Roll != prevState.gyroscope.Roll ||
		currState.gyroscope.Yaw != prevState.gyroscope.Yaw ||
		currState.gyroscope.Pitch != prevState.gyroscope.Pitch {
		event := EventGyroscopeUpdate
		if callback, ok := e.callback(event); ok {
			err := callback(currState.gyroscope)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkRightStick(currState, prevState *state) error {
	if currState.rightStick.X != prevState.rightStick.X ||
		currState.rightStick.Y != prevState.rightStick.Y {
		event := EventRightStickMove
		if callback, ok := e.callback(event); ok {
			err := callback(currState.rightStick)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkLeftStick(currState, prevState *state) error {
	if currState.leftStick.X != prevState.leftStick.X ||
		currState.leftStick.Y != prevState.leftStick.Y {
		event := EventLeftStickMove
		if callback, ok := e.callback(event); ok {
			err := callback(currState.leftStick)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkPS(currState, prevState *state) error {
	if currState.ps && !prevState.ps {
		event := EventPSPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.ps && prevState.ps {
		event := EventPSRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkOptions(currState, prevState *state) error {
	if currState.options && !prevState.options {
		event := EventOptionsPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.options && prevState.options {
		event := EventOptionsRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkShare(currState, prevState *state) error {
	if currState.share && !prevState.share {
		event := EventSharePress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.share && prevState.share {
		event := EventShareRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkDPad(currState, prevState *state) error {
	// D-pad up
	if currState.dPadUp && !prevState.dPadUp {
		event := EventDPadUpPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.dPadUp && prevState.dPadUp {
		event := EventDPadUpRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	// D-pad down
	if currState.dPadDown && !prevState.dPadDown {
		event := EventDPadDownPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.dPadDown && prevState.dPadDown {
		event := EventDPadDownRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	// D-pad left
	if currState.dPadLeft && !prevState.dPadLeft {
		event := EventDPadLeftPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.dPadLeft && prevState.dPadLeft {
		event := EventDPadLeftRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	// D-pad right
	if currState.dPadRight && !prevState.dPadRight {
		event := EventDPadRightPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.dPadRight && prevState.dPadRight {
		event := EventDPadRightRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkR3(currState, prevState *state) error {
	if currState.r3 && !prevState.r3 {
		event := EventR3Press
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.r3 && prevState.r3 {
		event := EventR3Release
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkR2(currState, prevState *state) error {
	if currState.r2 != prevState.r2 {
		event := EventR2Press
		if callback, ok := e.callback(event); ok {
			err := callback(currState.r2)
			if err != nil {
				return err
			}
		}
	}

	if currState.r2 == 0 && prevState.r2 != 0 {
		event := EventR2Release
		if callback, ok := e.callback(event); ok {
			err := callback(currState.r2)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkR1(currState, prevState *state) error {
	if currState.r1 && !prevState.r1 {
		event := EventR1Press
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.r1 && prevState.r1 {
		event := EventR1Release
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkL3(currState, prevState *state) error {
	if currState.l3 && !prevState.l3 {
		event := EventL3Press
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.l3 && prevState.l3 {
		event := EventL3Release
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkL2(currState, prevState *state) error {
	if currState.l2 != prevState.l2 {
		event := EventL2Press
		if callback, ok := e.callback(event); ok {
			err := callback(currState.l2)
			if err != nil {
				return err
			}
		}
	}

	if currState.l2 == 0 && prevState.l2 != 0 {
		event := EventL2Release
		if callback, ok := e.callback(event); ok {
			err := callback(currState.l2)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkL1(currState, prevState *state) error {
	if currState.l1 && !prevState.l1 {
		event := EventL1Press
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.l1 && prevState.l1 {
		event := EventL1Release
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkTriangle(currState, prevState *state) error {
	if currState.triangle && !prevState.triangle {
		event := EventTrianglePress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.triangle && prevState.triangle {
		event := EventTriangleRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkSquare(currState, prevState *state) error {
	if currState.square && !prevState.square {
		event := EventSquarePress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.square && prevState.square {
		event := EventSquareRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkCircle(currState, prevState *state) error {
	if currState.circle && !prevState.circle {
		event := EventCirclePress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.circle && prevState.circle {
		event := EventCircleRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *emitter) checkCross(currState, prevState *state) error {
	if currState.cross && !prevState.cross {
		event := EventCrossPress
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	if !currState.cross && prevState.cross {
		event := EventCrossRelease
		if callback, ok := e.callback(event); ok {
			err := callback(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func newEmitter() *emitter {
	e := &emitter{callbacks: make(map[Event]Callback)}
	e.checkers = []func(currState, prevState *state) error{
		e.checkCross,
		e.checkCircle,
		e.checkSquare,
		e.checkTriangle,
		e.checkL1,
		e.checkL2,
		e.checkL3,
		e.checkR1,
		e.checkR2,
		e.checkR3,
		e.checkDPad,
		e.checkShare,
		e.checkOptions,
		e.checkPS,
		e.checkLeftStick,
		e.checkRightStick,
		e.checkTouchpad,
		e.checkAccelerometer,
		e.checkGyroscope,
		e.checkBattery,
	}

	return e
}
