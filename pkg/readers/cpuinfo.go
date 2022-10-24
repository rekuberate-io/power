package readers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	cpuInfoPath           = "/proc/cpuinfo"
	physicalPackageIdPath = "/sys/devices/system/cpu/cpu%d/topology/physical_package_id"
)

type CpuInfo struct {
	Processor int
	Vendor    string
	Family    string
	Model     string
}

func (c *CpuInfo) String() string {
	return fmt.Sprintf("{ Processor: %d, Vendor: %s, Family: %s, Model: %s }", c.Processor, c.Vendor, c.Family, c.Model)
}

func DetectCpu() ([]*CpuInfo, error) {
	// cpuInfoCollection := make([]*CpuInfo, 0)
	cpuInfoCollection := []*CpuInfo{}

	var cpuInfo *CpuInfo

	file, err := os.Open(cpuInfoPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "processor") {
			if cpuInfo != nil {
				cpuInfoCollection = append(cpuInfoCollection, cpuInfo)
			}
			processor, _ := strconv.Atoi(text[12:])
			cpuInfo = &CpuInfo{Processor: processor}
		}

		if strings.HasPrefix(text, "vendor_id") {
			vendor_id := text[12:]
			cpuInfo.Vendor = vendor_id
		}

		if strings.HasPrefix(text, "cpu family") {
			cpu_family := text[12:]
			cpuInfo.Family = cpu_family
		}

		if strings.HasPrefix(text, "model		:") {
			model := text[9:]
			cpuInfo.Model = model
		}
	}

	if cpuInfo != nil {
		cpuInfoCollection = append(cpuInfoCollection, cpuInfo)
	}

	file.Close()
	return cpuInfoCollection, nil
}

func DetectPackages() (int, int, int, error) {
	packages := map[int]bool{}
	var totalCores, totalPackages int = 0, 0

	cpuInfoCollection, err := DetectCpu()
	if err != nil {
		return -1, -1 , len(packages), err
	}

	for cpuInfoIdx := 0; cpuInfoIdx < len(cpuInfoCollection); cpuInfoIdx++ {
		physicalPackageIdPath := fmt.Sprintf(physicalPackageIdPath, cpuInfoIdx)
		contents, err := os.ReadFile(physicalPackageIdPath)
		if err != nil {
			break
		}
		packageId, err := strconv.Atoi(strings.TrimSpace(string(contents)))
		if err == nil {
			if _, exists := packages[packageId]; !exists {
				totalPackages++
				packages[packageId] = true
			}
		}

		totalCores++
	}

	return totalCores, totalPackages, len(packages), err
}
