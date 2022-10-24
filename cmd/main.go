package main

import (
	"flag"
	"fmt"

	"github.com/rekuberate-io/power/pkg/readers"

	"k8s.io/klog/v2"
)

var (
	duration = flag.Uint("duration", 20, "duration in seconds")
	interval = flag.Uint("interval", 1, "interval in seconds")
)

func main() {
	defer exit()

	klog.InitFlags(nil)
	flag.Parse()

	klog.Infoln(fmt.Sprintf("starting power measuring session { duration: %dsec, interval: %dsec }", *duration, *interval))

	// endAt := time.Now().UTC().Add(time.Duration(*duration) * time.Second)

	// for endAt.After(time.Now().UTC()) {
	// 	// klog.Infoln("measuring RAPL...")
	// 	_, err := readers.DetectCpu()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	time.Sleep(time.Duration(*interval) * time.Second)
	// }

	// cpuInfoCollection, _ := readers.DetectCpu()
	// for _, cpuInfo := range cpuInfoCollection {
	// 	fmt.Println(cpuInfo)
	// }

	totalCores, totalPackages, lenPackages, err := readers.DetectPackages()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("tc: %d, tp: %d, lp: %d \n", totalCores, totalPackages, lenPackages)

	raplReader, err := readers.NewRaplReader(readers.MSR)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(raplReader.Available())

}

func exit() {
	exitCode := 10
	klog.Infoln("exiting power measuring session")
	klog.FlushAndExit(klog.ExitFlushTimeout, exitCode)
}
