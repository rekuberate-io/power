package main

import (
	"flag"
	"fmt"
	"time"

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

	klog.Infoln(fmt.Sprintf("starting rapl measuring session { duration: %dsec, interval: %dsec }", *duration, *interval))

	raplReader, err := readers.NewRaplReader(readers.Intel_Rapl)
	if err != nil {
		klog.Errorln(err)
	}

	endAt := time.Now().UTC().Add(time.Duration(*duration) * time.Second)

	for endAt.After(time.Now().UTC()) {
		klog.Infof("measuring...")
		measurements, err := raplReader.Read()
		if err != nil {
			klog.Errorln(err)
		}

		for k, v := range measurements {
			klog.Infof("%-10s %30v\n", k, v)
		}

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func exit() {
	exitCode := 10
	klog.Infoln("exiting rapl measuring session")
	klog.FlushAndExit(klog.ExitFlushTimeout, exitCode)
}
