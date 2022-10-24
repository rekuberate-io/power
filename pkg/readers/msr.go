package readers

import "fmt"

const (
	msrPath = "/dev/cpu/%d/msr"
)

type MsrReader struct {
}

func (r *MsrReader) Available() bool {
	return fileExists(fmt.Sprintf(msrPath, 0))
}
