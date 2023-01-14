package readers

import (
	"encoding/binary"
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"math"
	"syscall"
	"time"
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

	units map[int64]map[int]Units
)

//MsrReader is collecting RAPL results on Linux by using raw-access to the underlying MSRs under /dev/cpu/%d/msr. This requires root.
type MsrReader struct {
}

//Available checks if this RAPL reading strategy is available on this machine
func (r *MsrReader) Available() bool {
	return FileExists(fmt.Sprintf(msrPath, 0))
}

//Read a measurement using this reader strategy
func (r *MsrReader) Read() (Measurement, error) {
	pkgUnits, err := r.initUnits()
	if err != nil {
		panic(err)
	}

	units = pkgUnits

	klog.V(10).Infof("PkgUnits: %+v\n", units)

	before, err := r.measure()
	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	after, err := r.measure()
	if err != nil {
		panic(err)
	}

	delta := after.DeltaSum(before)

	return delta, nil
}

func (r *MsrReader) measure() (Measurement, error) {
	measurement := Measurement{}

	for _, cpu := range Cpus {
		for _, core := range cpu.Cores {
			byteOrder := cpu.ByteOrder
			var energy = Energy{}
			var fd int
			defer func(int) {
				err := r.close(fd)
				if err != nil {
					klog.Errorf("closing fd for core %d failed", core.Id)
				}
			}(fd)

			fd, err := r.open(core)
			if err != nil {
				klog.Errorf("opening fd for core %d failed", core.Id)
				return nil, err
			}

			cpuEnergyUnit := units[core.Package][core.Id].CpuEnergy
			dramEnergyUnit := units[core.Package][core.Id].DramEnergy

			energy.Pkg = r.readEnergy(fd, pkgEnergyStatus, cpuEnergyUnit, byteOrder)
			energy.PP0 = r.readEnergy(fd, pp0EnergyStatus, cpuEnergyUnit, byteOrder)
			energy.PP1 = r.readEnergy(fd, pp1EnergyStatus, cpuEnergyUnit, byteOrder)
			energy.DRAM = r.readEnergy(fd, dramEnergyStatus, dramEnergyUnit, byteOrder)
			energy.PSys = r.readEnergy(fd, psysEnergyStatus, cpuEnergyUnit, byteOrder)

			if _, exists := measurement[core.Package]; !exists {
				coreEnergy := make(map[int]Energy)
				coreEnergy[core.Id] = energy
				measurement[core.Package] = coreEnergy
			} else {
				measurement[core.Package][core.Id] = energy
			}
		}
	}
	return measurement, nil
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

func (r *MsrReader) readEnergy(fd int, offset int64, unit float64, order binary.ByteOrder) float64 {
	result, err := r.read(fd, offset, order)
	if err != nil {
		klog.Errorf("reading offset: %d failed, %s", offset, err)
		return 0
	}

	return unit * float64(result)
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

	for _, cpu := range Cpus {
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
