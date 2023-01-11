package readers

import (
	"bufio"
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

type Processor struct {
	Id       int
	Vendor   ProcessorVendor
	Family   int
	Model    ProcessorModel
	Cores    int
	Siblings int
}

type ProcessorModel struct {
	Id           int
	Name         string
	InternalName string
}

func (c *Processor) String() string {
	return fmt.Sprintf("{ Id: %d, Vendor: %s, Family: %d, Model: %s }", c.Id, c.Vendor.String(), c.Family, c.Model.InternalName)
}

func parseCpuInfo() ([]*Processor, error) {
	var processors []*Processor
	var processor *Processor
	const parseAt int = 12

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

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "processor") {
			id, _ := strconv.Atoi(text[parseAt:])
			processor = &Processor{Id: id, Model: ProcessorModel{}}
			processors = append(processors, processor)
		}

		if strings.HasPrefix(text, "vendor_id") {
			vendorId := text[parseAt:]

			switch vendorId {
			case "GenuineIntel":
				processor.Vendor = Intel
			case "AuthenticAMD":
				processor.Vendor = AMD
			default:
				processor.Vendor = NotAvailable
			}
		}

		if strings.HasPrefix(text, "cpu family") {
			family, _ := strconv.Atoi(text[parseAt+1:])
			processor.Family = family
		}

		if strings.HasPrefix(text, "model		:") {
			model, _ := strconv.Atoi(text[parseAt-3:])
			processor.Model.Id = model

			if cpuModel, exists := CpuModels[model]; exists {
				processor.Model.InternalName = cpuModel
			} else {
				processor.Model.InternalName = CpuModels[0]
			}
		}

		if strings.HasPrefix(text, "model name") {
			modelName := text[parseAt:]
			processor.Model.Name = modelName
		}

		if strings.HasPrefix(text, "cpu cores") {
			cpuCores, _ := strconv.Atoi(text[parseAt:])
			processor.Cores = cpuCores
		}

		if strings.HasPrefix(text, "siblings") {
			siblings, _ := strconv.Atoi(text[parseAt-1:])
			processor.Siblings = siblings
		}
	}

	//for _, processor := range processors {
	//	fmt.Println(processor)
	//}

	return processors, nil
}

func DetectPackages() (packages map[int64][]*Processor, err error) {
	packages = map[int64][]*Processor{}

	processors, err := parseCpuInfo()
	if err != nil {
		return nil, err
	}

	for _, processor := range processors {
		physicalPackageIdPath := fmt.Sprintf(physicalPackageIdPath, processor.Id)
		packageId, err := ReadIntFromFile(physicalPackageIdPath)
		if err == nil {
			if _, exists := packages[packageId]; !exists {
				var cores []*Processor
				packages[packageId] = append(cores, processor)
			} else {
				cores := packages[packageId]
				packages[packageId] = append(cores, processor)
			}
		}
	}

	return packages, err
}

type ProcessorVendor int

const (
	NotAvailable ProcessorVendor = iota
	Intel
	AMD
)

func (pv ProcessorVendor) String() string {
	var values []string = []string{"NotAvailable", "Intel", "AMD"}
	vendor := values[pv]

	return vendor
}
