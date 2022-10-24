package readers

import "fmt"

type RaplReaderStrategy int

const (
	First_Available RaplReaderStrategy = iota
	Intel_Rapl // Reading the files under /sys/class/powercap/intel-rapl/intel-rapl:0 using the powercap interface. This requires no special permissions, and was introduced in Linux 3.13
	Perf_Event // Using the perf_event interface with Linux 3.14 or newer. This requires root or a paranoid less than 1 (as do all system wide measurements with -a) sudo perf stat -a -e "power/energy-cores/" /bin/ls Available events can be found via perf list or under/sys/bus/event_source/devices/power/events/
	Raw_Access // Using raw-access to the underlying MSRs under /dev/msr. This requires root.
)

var raplDomains [5]string = [5]string{"energy-cores", "energy-gpu", "energy-pkg", "energy-ram", "energy-psys"}

type RaplReader interface {
	Available() bool
	Read() (map[string]uint64, error)
}

func NewRaplReader(forceRaplReaderStrategyIfAvailable RaplReaderStrategy) (RaplReader, error) {
	intelRaplReader := &IntelRapl{}
	perfEventReader := &PerfEventReader{}
	msrReader := &MsrReader{}

	switch forceRaplReaderStrategyIfAvailable {
	case Intel_Rapl:
		if intelRaplReader.Available() {
			return intelRaplReader, nil
		}
	case Perf_Event:
		if perfEventReader.Available() {
			return perfEventReader, nil
		}
	case Raw_Access:
		if msrReader.Available() {
			return msrReader, nil
		}
	case First_Available:
		if intelRaplReader.Available() {
			return intelRaplReader, nil
		} else if perfEventReader.Available() {
			return perfEventReader, nil
		} else if msrReader.Available() {
			return msrReader, nil
		}
	}

	if intelRaplReader.Available() {
		return intelRaplReader, nil
	} else if perfEventReader.Available() {
		return perfEventReader, nil
	} else if msrReader.Available() {
		return msrReader, nil
	}

	return nil, fmt.Errorf("no available power reader strategy")
}
