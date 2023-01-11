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

// PowerCap is collecting RAPL results on Linux by reading the files under /sys/class/powercap/intel-rapl/intel-rapl:0 using the powercap interface. This requires no special permissions, and was introduced in Linux 3.13
type PowerCap struct {
}

//Available checks if this RAPL reading strategy is available on this machine
func (r *PowerCap) Available() bool {
	return FileExists(fmt.Sprintf(zone, 0))
}

//Read a measurement using this reader strategy
func (r *PowerCap) Read() (map[string]uint64, error) {
	measurements := map[string]uint64{}
	var readError error

	//for packageId, enabled := range Packages {
	//	//if enabled {
	//	//	name, err := ReadStringFromFile(fmt.Sprintf(zoneName, packageId))
	//	//	if err != nil && !errors.Is(err, os.ErrNotExist) {
	//	//		klog.Errorln(err)
	//	//		readError = err
	//	//	} else if errors.Is(err, os.ErrNotExist) {
	//	//		break
	//	//	}
	//	//
	//	//	energy, err := ReadUintFromFile(fmt.Sprintf(zoneEnergy, packageId))
	//	//	if err != nil && !errors.Is(err, os.ErrNotExist) {
	//	//		klog.Errorln(err)
	//	//		readError = err
	//	//	} else if errors.Is(err, os.ErrNotExist) {
	//	//		break
	//	//	}
	//	//
	//	//	measurements[name] = energy
	//	//
	//	//	for raplDomainId := range raplDomains {
	//	//		name, err := ReadStringFromFile(fmt.Sprintf(subZoneName, packageId, packageId, raplDomainId))
	//	//		if err != nil && !errors.Is(err, os.ErrNotExist) {
	//	//			klog.Errorln(err)
	//	//			readError = err
	//	//		} else if errors.Is(err, os.ErrNotExist) {
	//	//			break
	//	//		}
	//	//
	//	//		energy, err := ReadUintFromFile(fmt.Sprintf(subZoneEnergy, packageId, packageId, raplDomainId))
	//	//		if err != nil && !errors.Is(err, os.ErrNotExist) {
	//	//			klog.Errorln(err)
	//	//			readError = err
	//	//		} else if errors.Is(err, os.ErrNotExist) {
	//	//			break
	//	//		}
	//	//
	//	//		measurements[name] = energy
	//	//	}
	//	//}
	//}

	return measurements, readError
}
