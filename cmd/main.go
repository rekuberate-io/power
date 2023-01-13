package main

import (
	"flag"
	"fmt"
	"github.com/rekuberate-io/power/pkg/readers"
	"os"

	"k8s.io/klog/v2"
)

var (
	duration = flag.Uint("duration", 30, "duration in seconds")
	interval = flag.Uint("interval", 1, "interval in seconds")
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
		*duration,
		*interval,
	))

	//endAt := time.Now().UTC().Add(time.Duration(*duration) * time.Second)
	//
	//for endAt.After(time.Now().UTC()) {
	//	klog.V(10).Infof("measuring...")
	//	measurements, err := raplReader.Read()
	//	if err != nil {
	//		klog.Errorln(err)
	//	}
	//
	//	for k, v := range measurements {
	//		klog.Infof("%-10s %30v\n", k, v)
	//	}
	//
	//	time.Sleep(time.Duration(*interval) * time.Second)
	//}

	_, err = raplReader.Read()
	if err != nil {
		klog.Errorln(err)
	}
}

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func exit() {
	//exitCode := 10
	klog.V(5).Infoln("exiting rapl measuring session")
	//klog.FlushAndExit(klog.ExitFlushTimeout, exitCode)
}
