package readers

const(
	perfEventPath = "/sys/bus/event_source/devices/power/type"
)

type PerfEventReader struct {
	
}

func (r *PerfEventReader) Available() bool {
	return fileExists(perfEventPath)
}