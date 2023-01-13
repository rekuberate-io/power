package main

import (
	"flag"
	"fmt"
	"github.com/rekuberate-io/power/pkg/readers"
	"os"
	"strings"

	"k8s.io/klog/v2"
)

var (
	strategy = flag.Int("strategy", 1, "rapl reader strategy")
)

func main() {
	defer exit()

	raplReader, err := readers.NewRaplReader(readers.RaplReaderStrategy(*strategy))
	if err != nil {
		klog.Fatalln(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		klog.Fatalln(err)
	}

	klog.V(5).Infoln(fmt.Sprintf(
		"starting rapl measuring session on %s { reader: %T, duration: %dsec, interval: %dsec }",
		hostname,
		raplReader,
	))

	for _, cpu := range readers.Cpus {
		fmt.Printf(
			"%s '%s/%s/Fam:%d' on socket %d (packages: %d, cores: %d)",
			cpu.Vendor.String(),
			strings.TrimSpace(cpu.Model.Name),
			cpu.Model.InternalName,
			cpu.Family,
			cpu.PhysicalId,
			len(cpu.Packages),
			len(cpu.Cores),
		)
	}

	fmt.Println()
	fmt.Println()

	measurement, err := raplReader.Read()
	if err != nil {
		klog.Errorln(err)
	}

	for pkgId, cores := range measurement {
		fmt.Printf("Package: %d\n", pkgId)
		for _, core := range cores {
			fmt.Printf("\t\tPackage energy: %v J\n", core.Pkg)
			fmt.Printf("\t\tPowerPlane0 (cores): %v J\n", core.PP0)
			fmt.Printf("\t\tPowerPlane1 (on-core GPU if avail): %v J\n", core.PP1)
			fmt.Printf("\t\tDRAM: %v J\n", core.DRAM)
			fmt.Printf("\t\tPSYS: %v J\n", core.PSys)
		}
	}
}

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func exit() {
	klog.V(5).Infoln("exiting rapl measuring session")
}
