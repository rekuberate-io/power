package readers

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/klog/v2"
)

type RaplReaderStrategy int

const (
	firstAvailable RaplReaderStrategy = iota
	sysfs                             // Reading the files under /sys/class/sysfs/intel-rapl/intel-rapl:0 using the sysfs interface. This requires no special permissions, and was introduced in Linux 3.13
	perf_event                        // Using the perf_event interface with Linux 3.14 or newer. This requires root or a paranoid less than 1 (as do all system wide measurements with -a) sudo perf stat -a -e "power/energy-cores/" /bin/ls Available events can be found via perf list or under/sys/bus/event_source/devices/power/events/
	msr                               // Using raw-access to the underlying MSRs under /dev/msr. This requires root.
)

var (
	Cpus                             map[int]*Cpu
	raplReaderStrategyNotImplemented error = errors.New("rapl reader strategy not implemented yet")
)

func init() {
	klog.V(5).Infoln("initializing rapl readers...")

	var err error
	Cpus, err = DetectPackages()
	if err != nil {
		klog.Errorln(err)
	}

	for _, cpu := range Cpus {
		klog.V(5).Infof("detected %s processor '%s/%s/Fam:%d' on socket %d (packages: %d, cores: %d)", cpu.Vendor.String(), strings.TrimSpace(cpu.Model.Name), cpu.Model.InternalName, cpu.Family, cpu.PhysicalId, len(cpu.Cores), len(cpu.Packages))

		switch cpu.Vendor {
		case NotAvailable:
			panic(errors.New("failed to determine the cpu vendor"))
		case AMD:
			if cpu.Family < AMDMinimumSupportedCpuFamily {
				panic(errors.New(fmt.Sprintf("unsupported cpu family, for amd processors it should be minimum: %d", AMDMinimumSupportedCpuFamily)))
			}
		case Intel:
			if cpu.Family < IntelMinimumSupportedCpuFamily {
				panic(errors.New(fmt.Sprintf("unsupported cpu family, for intel processors it should be minimum: %d", IntelMinimumSupportedCpuFamily)))
			}
		}
	}
}

type RaplReader interface {
	Available() bool
	Read() (Measurement, error)
	measure() (Measurement, error)
}

func NewRaplReader(forceRaplReaderStrategyIfAvailable RaplReaderStrategy) (RaplReader, error) {
	sysfsRaplReader := &Sysfs{}
	perfEventReader := &PerfEventReader{}
	msrReader := &MsrReader{}

	switch forceRaplReaderStrategyIfAvailable {
	case sysfs:
		if sysfsRaplReader.Available() {
			return sysfsRaplReader, nil
		}
	case perf_event:
		if perfEventReader.Available() {
			return perfEventReader, nil
		}
	case msr:
		if msrReader.Available() {
			return msrReader, nil
		}
	case firstAvailable:
		if sysfsRaplReader.Available() {
			return sysfsRaplReader, nil
		} else if perfEventReader.Available() {
			return perfEventReader, nil
		} else if msrReader.Available() {
			return msrReader, nil
		}
	}

	if sysfsRaplReader.Available() {
		return sysfsRaplReader, nil
	} else if perfEventReader.Available() {
		return perfEventReader, nil
	} else if msrReader.Available() {
		return msrReader, nil
	}

	return nil, fmt.Errorf("no available rapl reader strategy")
}
