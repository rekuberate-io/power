package readers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"strconv"
	"strings"
)

const (
	cpuInfoPath           = "/proc/cpuinfo"
	physicalPackageIdPath = "/sys/devices/system/cpu/cpu%d/topology/physical_package_id"
)

const (
	IntelMinimumSupportedCpuFamily int = 6
	AMDMinimumSupportedCpuFamily   int = 23
)

var CpuModels map[int]string

func init() {
	CpuModels = map[int]string{}

	CpuModels[0] = "CPU_UNKNOWN_MODEL"
	CpuModels[42] = "CPU_SANDYBRIDGE"
	CpuModels[45] = "CPU_SANDYBRIDGE_EP"
	CpuModels[58] = "CPU_IVYBRIDGE"
	CpuModels[62] = "CPU_IVYBRIDGE_EP"
	CpuModels[60] = "CPU_HASWELL"
	CpuModels[69] = "CPU_HASWELL_ULT"
	CpuModels[70] = "CPU_HASWELL_GT3E"
	CpuModels[63] = "CPU_HASWELL_EP"
	CpuModels[61] = "CPU_BROADWELL"
	CpuModels[71] = "CPU_BROADWELL_GT3E"
	CpuModels[79] = "CPU_BROADWELL_EP"
	CpuModels[78] = "CPU_SKYLAKE"
	CpuModels[94] = "CPU_SKYLAKE_HS"
	CpuModels[85] = "CPU_SKYLAKE_X"
	CpuModels[87] = "CPU_KNIGHTS_LANDING"
	CpuModels[133] = "CPU_KNIGHTS_MILL"
	CpuModels[142] = "CPU_KABYLAKE_MOBILE"
	CpuModels[158] = "CPU_KABYLAKE"
	CpuModels[55] = "CPU_ATOM_SILVERMONT"
	CpuModels[76] = "CPU_ATOM_AIRMONT"
	CpuModels[74] = "CPU_ATOM_MERRIFIELD"
	CpuModels[90] = "CPU_ATOM_MOOREFIELD"
	CpuModels[92] = "CPU_ATOM_GOLDMONT"
	CpuModels[122] = "CPU_ATOM_GEMINI_LAKE"
	CpuModels[95] = "CPU_ATOM_DENVERTON"
}

type Cpu struct {
	PhysicalId int
	Vendor     Vendor
	Model      Model
	Family     int
	Cores      []Core
	Packages   map[int64]bool
	ByteOrder  binary.ByteOrder
}

type Core struct {
	Id      int
	Package int64
}

type Model struct {
	Id           int
	Name         string
	InternalName string
}

func (c *Cpu) String() string {
	return fmt.Sprintf("{ Name: %s, Vendor: %s, Family: %d, Model: %s }", c.Model.Name, c.Vendor.String(), c.Family, c.Model.InternalName)
}

func GetNumberOfSockets() (int, error) {
	cpuSockets := make(map[int]bool)

	file, err := os.Open(cpuInfoPath)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			klog.Errorln(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "physical id") {
			id, _ := strconv.Atoi(text[12:])
			if _, exists := cpuSockets[id]; !exists {
				cpuSockets[id] = true
			}
		}
	}

	return len(cpuSockets), nil
}

func parseCpuInfo() (map[int]*Cpu, error) {
	const parseAt int = 12
	cpuSockets, err := GetNumberOfSockets()
	if err != nil {
		return nil, err
	}

	cpus := make(map[int]*Cpu, cpuSockets)

	file, err := os.Open(cpuInfoPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			klog.Errorln(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var cpu *Cpu = &Cpu{}

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "processor") {
			cpu = &Cpu{}

			id, _ := strconv.Atoi(text[parseAt:])
			core := Core{Id: id, Package: -1}

			cpu.Cores = append(cpu.Cores, core)
		}

		if strings.HasPrefix(text, "vendor_id") {
			vendorId := text[parseAt:]

			switch vendorId {
			case "GenuineIntel":
				cpu.Vendor = Intel
			case "AuthenticAMD":
				cpu.Vendor = AMD
			default:
				cpu.Vendor = NotAvailable
			}
		}

		if strings.HasPrefix(text, "cpu family") {
			family, _ := strconv.Atoi(text[parseAt+1:])
			cpu.Family = family
		}

		if strings.HasPrefix(text, "cpu family") {
			family, _ := strconv.Atoi(text[parseAt+1:])
			cpu.Family = family
		}

		if strings.HasPrefix(text, "model		:") {
			model, _ := strconv.Atoi(text[parseAt-3:])
			cpu.Model.Id = model

			if cpuModel, exists := CpuModels[model]; exists {
				cpu.Model.InternalName = cpuModel
			} else {
				cpu.Model.InternalName = CpuModels[0]
			}
		}

		if strings.HasPrefix(text, "model name") {
			modelName := text[parseAt:]
			cpu.Model.Name = modelName
		}

		if strings.HasPrefix(text, "physical id") {
			id, _ := strconv.Atoi(text[parseAt:])
			cpu.PhysicalId = id

			if _, exists := cpus[id]; !exists {
				endianness, err := GetEndianness()
				if err != nil {
					return nil, err
				}

				cpu.ByteOrder = endianness

				cpus[id] = cpu
			} else {
				cpus[id].Cores = append(cpus[id].Cores, cpu.Cores[0])
			}
		}
	}

	return cpus, nil
}

func DetectPackages() (map[int]*Cpu, error) {
	cpus, err := parseCpuInfo()
	if err != nil {
		return nil, err
	}

	for _, cpu := range cpus {
		cpu.Packages = make(map[int64]bool)

		for coreIdx, core := range cpu.Cores {
			physicalPackageIdPath := fmt.Sprintf(physicalPackageIdPath, core.Id)
			packageId, err := ReadIntFromFile(physicalPackageIdPath)
			if err == nil {
				if _, exists := cpu.Packages[packageId]; !exists {
					cpu.Packages[packageId] = true
				}

				cpu.Cores[coreIdx].Package = packageId
			}
		}
	}

	return cpus, err
}

type Vendor int

const (
	NotAvailable Vendor = iota
	Intel
	AMD
)

func (pv Vendor) String() string {
	var values []string = []string{"NotAvailable", "Intel", "AMD"}
	vendor := values[pv]

	return vendor
}
