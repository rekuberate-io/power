package readers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func (r *PerfEventReader) Read() (map[string]uint64, error) {

	t, err := os.ReadFile(perfEventPowerPath)
	if err != nil {
		return nil, err
	}

	fmt.Printf("event.type: %s", string(t))

	perfEvents := r.getPerfEventsMap()
	fmt.Println(len(perfEvents))

	return nil, nil
}

func (r *PerfEventReader) getPerfEventsMap() map[string]PerfEventAttr {
	perfEvents := map[string]PerfEventAttr{}

	for _, raplDomain := range raplDomains {
		if !FileExists(fmt.Sprintf(perfEventPowerEventsPath, raplDomain)) {
			continue
		}

		event, err := os.ReadFile(fmt.Sprintf(perfEventPowerEventsPath, raplDomain))
		if err != nil {
			break
		}

		scale, err := os.ReadFile(fmt.Sprintf(perfEventPowerEventsScalePath, raplDomain))
		if err != nil {
			break
		}

		unit, err := os.ReadFile(fmt.Sprintf(perfEventPowerEventsUnitPath, raplDomain))
		if err != nil {
			break
		}

		fmt.Printf("Event: %s \t Config: %s", raplDomain, event)
		fmt.Printf("Scale: %s %s\n", strings.TrimSpace(string(scale)), strings.TrimSpace(string(unit)))

		scaleAsFloat, err := strconv.ParseFloat(strings.TrimSpace(string(scale)), 64)
		if err != nil {
			scaleAsFloat = -1
		}

		perfEvent := PerfEventAttr{event: string(event), scale: scaleAsFloat, unit: string(unit)}
		perfEvents[raplDomain] = perfEvent
	}

	return perfEvents
}
