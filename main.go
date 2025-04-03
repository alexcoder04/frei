package main

import (
	"flag"
	"fmt"
	"os"
)

// types {{{
type MemData struct {
	MemTotal     float64
	MemUsed      float64
	MemShared    float64
	MemBuffers   float64
	MemCached    float64
	MemAvailable float64

	SwapFree  float64
	SwapUsed  float64
	SwapTotal float64
}

type DrawData struct {
	Buffers int
	Cache   int
	Free    int
	Shared  int
	Used    int

	SwapFree int
	SwapUsed int
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// }}}

// args {{{
var (
	Version   = "[built from source]"
	CommitSHA = ""

	dispHuman   = flag.Bool("h", false, "display human-readable numbers (implies -numbers)")
	dispKey     = flag.Bool("key", false, "display color key")
	dispVersion = flag.Bool("version", false, "display version and exit")
	longFormat  = flag.Bool("numbers", false, "print numbers in addition to the chart")
)

// }}}

func main() {
	flag.Parse()

	if *dispVersion {
		fmt.Printf("frei %s", Version)
		if CommitSHA != "" {
			fmt.Printf(" (commit %s)", CommitSHA)
		}
		fmt.Print("\n")
		os.Exit(0)
	}

	data, err := GetMemInfo()
	if err != nil {
		panic("Cannot get memory info")
	}

	chartWidth := getTerminalWidth() - 2
	barWidth := chartWidth - 4 - 5
	drawData := calcDrawData(data, barWidth)

	printCharts(drawData, chartWidth, barWidth, data)

	if *dispKey {
		printKey(chartWidth)
	}

	if *longFormat || *dispHuman {
		printNumbers(data, chartWidth, *dispHuman)
	}
}
