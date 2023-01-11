package readers

import "fmt"

const (
	msrPath = "/dev/cpu/%d/msr"
)

//MsrReader is collecting RAPL results on Linux by using raw-access to the underlying MSRs under /dev/cpu/%d/msr. This requires root.
type MsrReader struct {
}

//Available checks if this RAPL reading strategy is available on this machine
func (r *MsrReader) Available() bool {
	return FileExists(fmt.Sprintf(msrPath, 0))
}

//Read a measurement using this reader strategy
func (r *MsrReader) Read() (map[string]uint64, error) {
	for pkg, processors := range packages {
		for index, processor := range processors {
			fmt.Printf("Core: %d (Pkg: %d)", processor.Id, pkg)

			if index != len(processors)-1 {
				fmt.Print(", ")
			}
		}
	}

	fmt.Println()
	return nil, nil
}
