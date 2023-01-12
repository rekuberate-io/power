package readers

import (
	"fmt"
	"syscall"
)

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
	var fds []int

	for _, cpu := range cpus {
		for _, core := range cpu.Cores {
			fd, err := r.open(core)
			if err == nil {
				fds = append(fds, fd)
			}
		}
	}

	defer func([]int) {
		for _, fd := range fds {
			err := r.close(fd)
			if err != nil {

			}
		}
	}(fds)

	return nil, nil
}

func (r *MsrReader) open(core Core) (int, error) {
	path := fmt.Sprintf(msrPath, core.Id)

	fd, err := syscall.Open(path, syscall.O_RDONLY, 777)
	if err != nil {
		return -1, err
	}

	return fd, nil
}

func (r *MsrReader) close(fd int) error {
	if fd > 0 {
		err := syscall.Close(fd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MsrReader) read(fd int) (uint64, error) {
	return 0, nil
}
