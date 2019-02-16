package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kpeu3i/gods4"
	"github.com/kpeu3i/gods4/led"
	"github.com/kpeu3i/gods4/rumble"
)

func main() {
	// Find all controllers connected to your machine via USB or Bluetooth
	controllers := gods4.Find()
	if len(controllers) == 0 {
		panic("No connected DS4 controllers found")
	}

	// Select first controller from the list
	controller := controllers[0]

	// Connect to the controller
	err := controller.Connect()
	if err != nil {
		panic(err)
	}

	log.Printf("* Controller #1 | %-10s | name: %s, connection: %s\n", "Connect", controller, controller.ConnectionType())

	// Disconnect controller when a program is terminated
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		err := controller.Disconnect()
		if err != nil {
			panic(err)
		}
		log.Printf("* Controller #1 | %-10s | bye!\n", "Disconnect")
	}()

	// Register callback for "BatteryUpdate" event
	controller.On(gods4.EventBatteryUpdate, func(data interface{}) error {
		battery := data.(gods4.Battery)
		log.Printf("* Controller #1 | %-10s | capacity: %v%%, charging: %v, cable: %v\n",
			"Battery",
			battery.Capacity,
			battery.IsCharging,
			battery.IsCableConnected,
		)

		return nil
	})

	// Register callback for "CrossPress" event
	controller.On(gods4.EventCrossPress, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: press\n", "Cross")

		return nil
	})

	// Register callback for "CrossRelease" event
	controller.On(gods4.EventCrossRelease, func(data interface{}) error {
		log.Printf("* Controller #1 | %-10s | state: release\n", "Cross")

		return nil
	})

	// Register callback for "RightStickMove" event
	controller.On(gods4.EventRightStickMove, func(data interface{}) error {
		stick := data.(gods4.Stick)
		log.Printf("* Controller #1 | %-10s | x: %v, y: %v\n", "RightStick", stick.X, stick.Y)

		return nil
	})

	// Enable left and right rumble motors
	err = controller.Rumble(rumble.Both())
	if err != nil {
		panic(err)
	}

	// Enable LED (yellow) with flash
	err = controller.Led(led.Yellow().Flash(50, 50))
	if err != nil {
		panic(err)
	}

	// Start listening for controller events
	err = controller.Listen()
	if err != nil {
		panic(err)
	}

	// Output:
	// 2019/02/16 17:00:23 * Controller #1 | Connect    | name: Wireless Controller (vendor: 1356, product: 2508), connection: BT
	// 2019/02/16 17:00:23 * Controller #1 | Battery    | capacity: 77%, charging: false, cable: false
	// 2019/02/16 17:00:34 * Controller #1 | Cross      | state: press
	// 2019/02/16 17:00:34 * Controller #1 | Cross      | state: release
	// 2019/02/16 17:00:39 * Controller #1 | RightStick | x: 187, y: 98
	// 2019/02/16 17:00:39 * Controller #1 | RightStick | x: 191, y: 94
	// 2019/02/16 17:00:39 * Controller #1 | RightStick | x: 196, y: 93
	// 2019/02/16 17:00:39 * Controller #1 | RightStick | x: 212, y: 88
	// 2019/02/16 17:00:39 * Controller #1 | RightStick | x: 228, y: 79
	// 2019/02/16 17:02:52 * Controller #1 | Disconnect | bye!
}
