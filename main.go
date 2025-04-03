package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// types {{{
type MemData struct {
	MemTotal     float64
	MemUsed      float64
	MemShared    float64
	MemBuffers   float64
	MemCached    float64
	MemAvailable float64
	MemFree      float64

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

	dispHuman   = flag.Bool("h", false, "display human-readable numbers (implies -table)")
	dispTable   = flag.Bool("table", false, "print table with numbers in addition to the chart")
	dispVersion = flag.Bool("version", false, "display version and exit")
)

// }}}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options]\n\nOptions:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
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

	if *dispTable || *dispHuman {
		printTable(data, chartWidth, *dispHuman)
	}
}
