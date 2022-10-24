package readers

import (
	"fmt"
)

const (
	zone        = "/sys/class/powercap/intel-rapl/intel-rapl:%d/"
	zoneName    = "/sys/class/powercap/intel-rapl/intel-rapl:%d/name"
	zoneEnergy  = "/sys/class/powercap/intel-rapl/intel-rapl:%d/energy_uj"
	subZoneName   = "/sys/class/powercap/intel-rapl/intel-rapl:%d/intel-rapl:%d:%d/name"
	subZoneEnergy = "/sys/class/powercap/intel-rapl/intel-rapl:%d/intel-rapl:%d:%d/energy_uj"
)

type IntelRapl struct {
}

func (r *IntelRapl) Available() bool {
	return FileExists(fmt.Sprintf(zone, 0))
}

func (r *IntelRapl) Read() (map[string]uint64, error) {
	measurements := map[string]uint64{}
	_, packages, err := DetectPackages()
	if err != nil {
		return nil, err
	}

	for packageId, enabled := range packages {
		if enabled {
			name, err := ReadStringFromFile(fmt.Sprintf(zoneName, packageId))
			if err != nil {
				break
			}

			energy, err := ReadUintFromFile(fmt.Sprintf(zoneEnergy, packageId))
			if err != nil {
				break
			}

			measurements[name] = energy

			for raplDomainId := range raplDomains {
				name, err := ReadStringFromFile(fmt.Sprintf(subZoneName, packageId, packageId, raplDomainId))
				if err != nil {
					break
				}

				energy, err := ReadUintFromFile(fmt.Sprintf(subZoneEnergy, packageId, packageId, raplDomainId))
				if err != nil {
					break
				}

				measurements[name] = energy
			}
		}
	}

	return measurements, nil
}
