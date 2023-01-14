package readers

const (
	perfEventPowerPath            = "/sys/bus/event_source/devices/power/type"
	perfEventPowerEventsPath      = "/sys/bus/event_source/devices/power/events/%s"
	perfEventPowerEventsScalePath = "/sys/bus/event_source/devices/power/events/%s.scale"
	perfEventPowerEventsUnitPath  = "/sys/bus/event_source/devices/power/events/%s.unit"
)

type PerfEventAttr struct {
	event string
	scale float64
	unit  string
}

type PerfEventReader struct {
}

func (r *PerfEventReader) Available() bool {
	return FileExists(perfEventPowerPath)
}

func (r *PerfEventReader) Read() (Measurement, error) {
	return nil, raplReaderStrategyNotImplemented
}

func (r *PerfEventReader) measure() (Measurement, error) {
	return nil, raplReaderStrategyNotImplemented
}
