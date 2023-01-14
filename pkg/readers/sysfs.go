package readers

import (
	"errors"
	"fmt"
	"os"
	"time"
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

// Sysfs is collecting RAPL results on Linux by reading the files under /sys/class/powercap/intel-rapl/intel-rapl:0 using the sysfs interface. This requires no special permissions, and was introduced in Linux 3.13
type Sysfs struct {
}

//Available checks if this RAPL reading strategy is available on this machine
func (r *Sysfs) Available() bool {
	return FileExists(fmt.Sprintf(zone, 0))
}

//Read a measurement using this reader strategy
func (r *Sysfs) Read() (Measurement, error) {
	before, err := r.measure()
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	after, err := r.measure()
	if err != nil {
		panic(err)
	}

	delta := after.DeltaSum(before)

	return delta, nil
}

func (r *Sysfs) measure() (Measurement, error) {
	measurement := Measurement{}

	for _, cpu := range Cpus {
		for pkg, _ := range cpu.Packages {
			_, err := ReadStringFromFile(fmt.Sprintf(zoneName, pkg))
			if err != nil {
				return nil, err
			}

			res, err := ReadUintFromFile(fmt.Sprintf(zoneEnergy, pkg))
			if err != nil {
				return nil, err
			}

			result := float64(res) / 1000000.0

			energyPack := make(map[int]Energy)
			energy := Energy{
				Pkg: result,
			}

			measurement[pkg] = energyPack

			for domain, _ := range raplDomains {
				name, err := ReadStringFromFile(fmt.Sprintf(subZoneName, pkg, pkg, domain))
				if err != nil && !errors.Is(err, os.ErrNotExist) {
					return nil, err
				} else if errors.Is(err, os.ErrNotExist) {
					continue
				}

				res, err := ReadUintFromFile(fmt.Sprintf(subZoneEnergy, pkg, pkg, domain))
				if err != nil && !errors.Is(err, os.ErrNotExist) {
					return nil, err
				} else if errors.Is(err, os.ErrNotExist) {
					continue
				}

				result := float64(res) / 1000000.0

				switch name {
				case "core":
					energy.PP0 = result
				case "uncore":
					energy.PP1 = result
				case "dram":
					energy.DRAM = result
				}

				measurement[pkg][0] = energy
			}
		}
	}

	return measurement, nil
}
