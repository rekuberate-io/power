package readers

import "fmt"

const (
	packageNamePath = "/sys/class/powercap/intel-rapl/intel-rapl:%d/"
)

type RaplReaderStrategy int

const (
	FirstAvailable RaplReaderStrategy = iota
	Powercap
	PerfEvent
	MSR
)

type RaplReader interface {
	Available() bool
}

func NewRaplReader(forceRaplReaderStrategyIfAvailable RaplReaderStrategy) (RaplReader, error) {
	powercapReader := &PowercapReader{}
	perfEventReader := &PerfEventReader{}
	msrReader := &MsrReader{}

	switch forceRaplReaderStrategyIfAvailable {
    case Powercap:
        if powercapReader.Available() {
			return powercapReader, nil
		}
    case PerfEvent:
        if perfEventReader.Available() {
			return perfEventReader, nil
		} 
	case MSR:
        if msrReader.Available() {
			return msrReader, nil
		}
	case FirstAvailable:
        if powercapReader.Available() {
			return powercapReader, nil
		} else if perfEventReader.Available() {
			return perfEventReader, nil
		} else if msrReader.Available() {
			return msrReader, nil
		}
    }

	if powercapReader.Available() {
		return powercapReader, nil
	} else if perfEventReader.Available() {
		return perfEventReader, nil
	} else if msrReader.Available() {
		return msrReader, nil
	}

	return nil, fmt.Errorf("no available power reader strategy")
}