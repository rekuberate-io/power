package readers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"math"
	"syscall"
	"unsafe"
)

const (
	msrPath = "/dev/cpu/%d/msr"
)

var (
	raplUnits        int64
	pkgEnergyStatus  int64
	pp0EnergyStatus  int64
	pp1EnergyStatus  int64
	dramEnergyStatus int64
	psysEnergyStatus int64
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
	pkgUnits, err := r.initUnits()
	if err != nil {
		panic(err)
	}

	klog.V(10).Infof("PkgUnits: %+v\n", pkgUnits)

	for _, cpu := range cpus {
		var energy1 = Energy{}

		for _, core := range cpu.Cores {
			var fds []int
			defer func([]int) {
				for _, fd := range fds {
					err := r.close(fd)
					if err != nil {
						klog.Errorln("closing fd failed")
					}
				}
			}(fds)

			fd, err := r.open(core)
			if err == nil {
				fds = append(fds, fd)
			}

			for _, fd := range fds {
				energy1.Pkg += r.readEnergy(fd, pkgEnergyStatus, pkgUnits[core.Package][core.Id].CpuEnergy, cpu.ByteOrder)
				energy1.PP0 += r.readEnergy(fd, pp0EnergyStatus, pkgUnits[core.Package][core.Id].CpuEnergy, cpu.ByteOrder)
				energy1.PP1 += r.readEnergy(fd, pp1EnergyStatus, pkgUnits[core.Package][core.Id].CpuEnergy, cpu.ByteOrder)
				energy1.DRAM += r.readEnergy(fd, dramEnergyStatus, pkgUnits[core.Package][core.Id].DramEnergy, cpu.ByteOrder)
				energy1.PSys += r.readEnergy(fd, psysEnergyStatus, pkgUnits[core.Package][core.Id].CpuEnergy, cpu.ByteOrder)
			}
		}

		//var energy3 = Energy{}
		//energy3.Pkg = energy2.Pkg - energy1.Pkg
		//energy3.PP0 = energy2.PP0 - energy1.PP0
		//energy3.PP1 = energy2.PP1 - energy1.PP1
		//energy3.DRAM = energy2.DRAM - energy1.DRAM
		//energy3.PSys = energy2.PSys - energy1.PSys

		klog.Infof("Energy: %+v\n", energy1)
	}

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

func (r *MsrReader) readEnergy(fd int, offset int64, unit float64, order binary.ByteOrder) uint64 {
	result, err := r.read(fd, offset, order)
	if err != nil {
		klog.Errorln("reading offset:%d failed", offset, err)
		return 0
	}

	return uint64(unit * float64(result))
}

// /dev/cpu/CPUNUM/msr provides an interface to read and write the
// model-specific registers (MSRs) of an x86 CPU.  CPUNUM is the
// number of the CPU to access as listed in /proc/cpuinfo.
//
// The register access is done by opening the file and seeking to
// the MSR number as offset in the file, and then reading or writing
// in chunks of 8 bytes.  An I/O transfer of more than 8 bytes means
// multiple reads or writes of the same register.
//
// This file is protected so that it can be read and written only by
// the user root, or members of the group root.
// https://man7.org/linux/man-pages/man4/msr.4.html
func (r *MsrReader) read(fd int, offset int64, order binary.ByteOrder) (uint64, error) {
	var chunkType uint64
	chunkSize := int(unsafe.Sizeof(chunkType))

	buffer := make([]byte, chunkSize)
	bytes, err := syscall.Pread(fd, buffer, offset)
	if err != nil {
		return 0, err
	}

	if bytes != chunkSize {
		return 0, errors.New(fmt.Sprintf("wrong reading: %d", bytes))
	}

	result := order.Uint64(buffer)
	return result, nil
}

func (r *MsrReader) initPerVendor(cpu Cpu) {
	switch cpu.Vendor {
	case AMD:
		raplUnits = MSR_AMD_RAPL_POWER_UNIT
		pkgEnergyStatus = MSR_AMD_PKG_ENERGY_STATUS
		pp0EnergyStatus = MSR_AMD_PP0_ENERGY_STATUS
	case Intel:
		raplUnits = MSR_INTEL_RAPL_POWER_UNIT
		pkgEnergyStatus = MSR_INTEL_PKG_ENERGY_STATUS
		pp0EnergyStatus = MSR_INTEL_PP0_ENERGY_STATUS
		pp1EnergyStatus = MSR_PP1_ENERGY_STATUS
		dramEnergyStatus = MSR_DRAM_ENERGY_STATUS
		psysEnergyStatus = MSR_PLATFORM_ENERGY_STATUS
	}
}

func (r *MsrReader) initUnits() (map[int64]map[int]Units, error) {
	pkgUnits := make(map[int64]map[int]Units)

	for _, cpu := range cpus {
		r.initPerVendor(*cpu)

		for _, core := range cpu.Cores {
			var fds []int
			defer func([]int) {
				for _, fd := range fds {
					err := r.close(fd)
					if err != nil {
						klog.Errorln("closing fd failed")
					}
				}
			}(fds)

			fd, err := r.open(core)
			if err == nil {
				fds = append(fds, fd)
			}

			for _, fd := range fds {
				result, err := r.read(fd, raplUnits, cpu.ByteOrder)
				if err != nil {
					return nil, err
				}

				var units = Units{
					Power:      math.Pow(0.5, float64(result&0xf)),
					Time:       math.Pow(0.5, float64((result>>16)&0xf)),
					CpuEnergy:  math.Pow(0.5, float64((result>>8)&0x1f)),
					DramEnergy: math.Pow(0.5, float64((result>>8)&0x1f)),
				}

				if _, exists := pkgUnits[core.Package]; !exists {
					coreUnits := make(map[int]Units)
					coreUnits[core.Id] = units
					pkgUnits[core.Package] = coreUnits
				} else {
					pkgUnits[core.Package][core.Id] = units
				}
			}
		}
	}

	return pkgUnits, nil
}
