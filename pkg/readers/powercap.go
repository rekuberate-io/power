package readers

import (
	"fmt"
)

type PowercapReader struct {
}

func (r *PowercapReader) Available() bool {
	return fileExists(fmt.Sprintf(packageNamePath, 0))
}
