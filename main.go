package main

import (
	"flag"
	"fmt"
	"github.com/Banrai/PiScan/scanner"
	"github.com/Unknwon/macaron"  
	"log"
)

func main() {
	var (
		device string
	)

	flag.StringVar(&device, "device", scanner.SCANNER_DEVICE, fmt.Sprintf("The '/dev/input/event' device associated with your scanner (defaults to '%s')", scanner.SCANNER_DEVICE))

	processScanFn := func(barcode string) {

		println(barcode)

	}

	errorFn := func(e error) {
		log.Fatal(e)
	}

	scanner.ScanForever(device, processScanFn, errorFn)

	m := macaron.New()

	m.Get("/", func() string {
		println("hello world")
		return "hello world" // HTTP 200 : "hello world"
	})

	m.Run()
}
