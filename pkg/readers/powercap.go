package readers

import (
	"fmt"
)

const (
	zone          = "/sys/class/powercap/intel-rapl/intel-rapl:%d/"
	zoneName      = "/sys/class/powercap/intel-rapl/intel-rapl:%d/name"
	zoneEnergy    = "/sys/class/powercap/intel-rapl/intel-rapl:%d/energy_uj"
	subZoneName   = "/sys/class/powercap/intel-rapl/intel-rapl:%d/intel-rapl:%d:%d/name"
	subZoneEnergy = "/sys/class/powercap/intel-rapl/intel-rapl:%d/intel-rapl:%d:%d/energy_uj"
)

var (
	raplDomains [5]string = [5]string{"energy-cores", "energy-gpu", "energy-pkg", "energy-ram", "energy-psys"}
)

// PowerCap is collecting RAPL results on Linux by reading the files under /sys/class/powercap/intel-rapl/intel-rapl:0 using the powercap interface. This requires no special permissions, and was introduced in Linux 3.13
type PowerCap struct {
}

//Available checks if this RAPL reading strategy is available on this machine
func (r *PowerCap) Available() bool {
	return FileExists(fmt.Sprintf(zone, 0))
}

//Read a measurement using this reader strategy
func (r *PowerCap) Read() (Measurement, error) {
	return nil, nil
}
