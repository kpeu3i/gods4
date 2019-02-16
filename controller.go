package gods4

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"sync"

	"github.com/pkg/errors"

	"github.com/kpeu3i/gods4/hid"
	"github.com/kpeu3i/gods4/led"
	"github.com/kpeu3i/gods4/rumble"
)

var (
	ErrInvalidConnectionType    = errors.New("ds4: can't detect connection type")
	ErrControllerIsConnected    = errors.New("ds4: controller is already connected")
	ErrControllerIsNotConnected = errors.New("ds4: controller is not connected")
	ErrControllerIsListening    = errors.New("ds4: controller is already listening for events")
)

const getFeatureReportCode0x04 = 0x04

type Device interface {
	VendorID() uint16
	ProductID() uint16
	Path() string
	Release() uint16
	Serial() string
	Manufacturer() string
	Product() string
	Open() error
	Close() error
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	GetFeatureReport(code byte) ([]byte, error)
}

type Controller struct {
	mutex          sync.RWMutex
	device         Device
	connectionType ConnectionType
	emitter        *emitter
	inputOffset    uint
	inputCurrState *state
	inputPrevState *state
	outputOffset   uint
	outputState    []byte
	isListening    bool
	errors         chan error
	quit           chan struct{}
}

type Callback func(data interface{}) error

func NewController(device Device) *Controller {
	return &Controller{
		device:         device,
		connectionType: ConnectionTypeNone,
		emitter:        newEmitter(),
		errors:         make(chan error),
		quit:           make(chan struct{}),
	}
}

func (c *Controller) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.errorIfConnected()
	if err != nil {
		return err
	}

	err = c.device.Open()
	if err != nil {
		return err
	}

	connectionType, err := detectConnectionType(c.device)
	if err != nil {
		return err
	}

	c.connectionType = connectionType

	switch c.connectionType {
	case ConnectionTypeBluetooth:
		_, err = c.device.GetFeatureReport(getFeatureReportCode0x04)
		if err != nil {
			return err
		}

		c.inputOffset = 2
		c.outputOffset = 3

		c.outputState = make([]byte, 79)
		c.outputState[0] = 0xA2
		c.outputState[1] = 0x11
		c.outputState[2] = 0x80
		c.outputState[4] = 0x0F
	case ConnectionTypeUSB:
		c.inputOffset = 0
		c.outputOffset = 0

		c.outputState = make([]byte, 79)
		c.outputState[0] = 0x05
		c.outputState[1] = 0x07
	}

	return nil
}

func (c *Controller) Disconnect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.errorIfNotConnected()
	if err != nil {
		return err
	}

	c.quit <- struct{}{}

	err = c.device.Close()
	if err != nil {
		return err
	}

	c.connectionType = ConnectionTypeNone

	c.errors <- nil

	return nil
}

func (c *Controller) ConnectionType() ConnectionType {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.connectionType
}

func (c *Controller) Listen() error {
	c.mutex.Lock()

	err := c.errorIfNotConnected()
	if err != nil {
		return err
	}

	err = c.errorIfListening()
	if err != nil {
		return err
	}

	c.isListening = true
	defer func() {
		c.isListening = false
	}()

	go c.handle()

	c.mutex.Unlock()

	return <-c.errors
}

func (c *Controller) On(event Event, fn Callback) {
	c.emitter.setCallback(event, fn)
}

func (c *Controller) Off(event Event) {
	c.emitter.unsetCallback(event)
}

func (c *Controller) Rumble(rumble *rumble.Rumble) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.errorIfNotConnected()
	if err != nil {
		return err
	}

	patch := make(map[uint]byte, 2)
	patch[4+c.outputOffset] = rumble.Left()
	patch[5+c.outputOffset] = rumble.Right()

	return c.set(patch)
}

func (c *Controller) Led(led *led.Led) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.errorIfNotConnected()
	if err != nil {
		return err
	}

	patch := make(map[uint]byte, 3)
	patch[6+c.outputOffset] = led.Red()
	patch[7+c.outputOffset] = led.Green()
	patch[8+c.outputOffset] = led.Blue()
	patch[9+c.outputOffset] = led.FlashOn()
	patch[10+c.outputOffset] = led.FlashOff()

	return c.set(patch)
}

func (c *Controller) handle() {
	bytes := make([]byte, 64)
	bytes[0+c.inputOffset] = 1
	bytes[1+c.inputOffset] = 128
	bytes[2+c.inputOffset] = 128
	bytes[3+c.inputOffset] = 128
	bytes[4+c.inputOffset] = 128
	bytes[5+c.inputOffset] = 8

	c.inputPrevState = newState(bytes, c.inputOffset, nil)

	for {
		select {
		case <-c.quit:
			return
		default:
			_, err := c.device.Read(bytes)
			if err != nil {
				c.errors <- err

				return
			}

			c.inputCurrState = newState(bytes, c.inputOffset, c.inputPrevState)

			err = c.emitter.emit(c.inputCurrState, c.inputPrevState)
			if err != nil {
				c.errors <- err

				return
			}

			c.inputPrevState = c.inputCurrState
		}
	}
}

func (c *Controller) VendorID() uint16 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.device.VendorID()
}

func (c *Controller) ProductID() uint16 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.device.ProductID()
}

func (c *Controller) Name() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.device.Product()
}

func (c *Controller) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return fmt.Sprintf("%s (vendor: %v, product: %v)", c.device.Product(), c.device.VendorID(), c.device.ProductID())
}

func (c *Controller) set(patch map[uint]byte) error {
	for i, b := range patch {
		c.outputState[i] = b
	}

	switch c.connectionType {
	case ConnectionTypeBluetooth:
		crc := crc32.ChecksumIEEE(c.outputState[0:75])
		binary.LittleEndian.PutUint32(c.outputState[75:], crc)
		_, err := c.device.Write(c.outputState[1:])
		if err != nil {
			return err
		}
	case ConnectionTypeUSB:
		_, err := c.device.Write(c.outputState)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) errorIfConnected() error {
	if c.connectionType != ConnectionTypeNone {
		return ErrControllerIsConnected
	}

	return nil
}

func (c *Controller) errorIfNotConnected() error {
	if c.connectionType == ConnectionTypeNone {
		return ErrControllerIsNotConnected
	}

	return nil
}

func (c *Controller) errorIfListening() error {
	if c.isListening {
		return ErrControllerIsListening
	}

	return nil
}

func Find() []*Controller {
	devices := hid.Find()
	controllers := make([]*Controller, 0, len(devices))
	for _, device := range devices {
		controllers = append(controllers, NewController(device))
	}

	return controllers
}
