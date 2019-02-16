package gods4

type ConnectionType uint

const (
	ConnectionTypeNone ConnectionType = iota
	ConnectionTypeUSB
	ConnectionTypeBluetooth
)

func (t ConnectionType) String() string {
	switch t {
	case ConnectionTypeNone:
		return "NONE"
	case ConnectionTypeUSB:
		return "USB"
	case ConnectionTypeBluetooth:
		return "BT"
	default:
		return ""
	}
}

func detectConnectionType(device Device) (ConnectionType, error) {
	_, _ = device.GetFeatureReport(getFeatureReportCode0x04)

	bytes := make([]byte, 2)
	for i := 1; i <= 100; i++ {
		_, err := device.Read(bytes)
		if err != nil {
			return 0, ErrInvalidConnectionType
		}
	}

	if bytes[0] == 0x11 && bytes[1] == 0xC0 {
		return ConnectionTypeBluetooth, nil
	}

	return ConnectionTypeUSB, nil
}
