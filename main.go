package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexcoder04/meme"
)

// types {{{
type MemData struct {
	Buffers      float64
	Cached       float64
	MemAvailable float64
	MemFree      float64
	MemTotal     float64
	Shared       float64
	Used         float64

	SwapFree  float64
	SwapTotal float64
	SwapUsed  float64
}

type DrawData struct {
	Buffers uint
	Cache   uint
	Free    uint
	Shared  uint
	Used    uint

	SwapFree uint
	SwapUsed uint
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
	Version   = ""
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
		if Version == "" {
			fmt.Print("frei [built from source]")
		} else {
			fmt.Print("frei " + Version)
		}
		if CommitSHA != "" {
			fmt.Print(" (commit " + CommitSHA + ")")
		}
		fmt.Print("\n")
		os.Exit(0)
	}

	data, err := meme.GetMemInfo()
	if err != nil {
		panic("Cannot get memory info")
	}

	chartWidth := getTerminalWidth() - 2
	barWidth := chartWidth - 4 - 5
	drawData := calcDrawData(data, barWidth)

	printCharts(drawData, chartWidth, data)
	if *dispKey {
		printKey(chartWidth)
	}
	if *longFormat || *dispHuman {
		printNumbers(data, chartWidth, *dispHuman)
	}
}
