package readers

import "fmt"

const (
	msrPath = "/dev/cpu/%d/msr"
)

type MsrReader struct {
}

func (r *MsrReader) Available() bool {
	return FileExists(fmt.Sprintf(msrPath, 0))
}

func (r *MsrReader) Read() (map[string]uint64, error) {
	return nil, nil
}
