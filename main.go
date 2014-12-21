package main

import (
	"flag"
	"fmt"
	"github.com/Banrai/PiScan/scanner"

	"github.com/Unknwon/macaron"
	//"log"
)

var Printcode string

func main() {
	var (
		device string
	)

	flag.StringVar(&device, "device", scanner.SCANNER_DEVICE, fmt.Sprintf("The '/dev/input/event' device associated with your scanner (defaults to '%s')", scanner.SCANNER_DEVICE))

	processScanFn := func(barcode string) {
		fmt.Println("newcode:" + barcode)
		Printcode = barcode
	}

	errorFn := func(e error) {
		fmt.Println(e)
	}
	fmt.Println("capturing barcode scanner")
	go scanner.ScanForever(device, processScanFn, errorFn)

	fmt.Println("web server running")
	Printcode = "test"
	m := macaron.New()
	m.Get("/", func() (string barcode) {
		barcode = Printcode
		Printcode = ""
	})

	m.Run()

}
