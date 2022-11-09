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
	var cpuInfoCollection []*CpuInfo

	var cpuInfo *CpuInfo

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
			if cpuInfo != nil {
				cpuInfoCollection = append(cpuInfoCollection, cpuInfo)
			}
			processor, _ := strconv.Atoi(text[12:])
			cpuInfo = &CpuInfo{Processor: processor}
		}

		if strings.HasPrefix(text, "vendor_id") {
			vendorId := text[12:]
			cpuInfo.Vendor = vendorId
		}

		if strings.HasPrefix(text, "cpu family") {
			cpuFamily := text[12:]
			cpuInfo.Family = cpuFamily
		}

		if strings.HasPrefix(text, "model		:") {
			model := text[9:]
			cpuInfo.Model = model
		}
	}

	if cpuInfo != nil {
		cpuInfoCollection = append(cpuInfoCollection, cpuInfo)
	}

	return cpuInfoCollection, nil
}

func DetectPackages() (totalCores int, packages map[int64]bool, err error) {
	packages = map[int64]bool{}
	totalCores = 0

	cpuInfoCollection, err := DetectCpu()
	if err != nil {
		return -1, nil, err
	}

	for cpuInfoIdx := 0; cpuInfoIdx < len(cpuInfoCollection); cpuInfoIdx++ {
		physicalPackageIdPath := fmt.Sprintf(physicalPackageIdPath, cpuInfoIdx)
		packageId, err := ReadIntFromFile(physicalPackageIdPath)
		if err == nil {
			if _, exists := packages[packageId]; !exists {
				// totalPackages++
				packages[packageId] = true
			}
		}

		totalCores++
	}

	return totalCores, packages, err
}
