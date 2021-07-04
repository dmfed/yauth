package yauth

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

type device struct {
	deviceName string
	deviceID   string
}

// newDevice accepts g name as string and generates Device with name
// and UUID as DeviceID. If empty string is passed as argument then
// hostname is used as name. If we can't get hostname, name is set to
// "unknown".
func newDevice() device {
	var d device
	d.deviceID = generateUUID()
	hostname, err := os.Hostname()
	if err == nil {
		d.deviceName = hostname
	} else {
		d.deviceName = fmt.Sprintf("device with id: %v", d.deviceID)
	}
	return d
}

// generateUUID() returns UUIDs based on RFC 4122 and
// DCE 1.1: Authentication and Security Services.
func generateUUID() string {
	return uuid.NewString()
}
