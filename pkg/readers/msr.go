package readers

import "fmt"

const (
	msrPath = "/dev/cpu/%d/msr"
)

//MsrReader is collecting RAPL results on Linux by using raw-access to the underlying MSRs under /dev/msr. This requires root.
type MsrReader struct {
}

//Available checks if this RAPL reading strategy is available on this machine
func (r *MsrReader) Available() bool {
	return FileExists(fmt.Sprintf(msrPath, 0))
}

//Read a measurement using this reader strategy
func (r *MsrReader) Read() (map[string]uint64, error) {
	return nil, nil
}
