package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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

var (
	Version   = ""
	CommitSHA = ""

	dispHuman   = flag.Bool("h", false, "display human-readable numbers (implies -numbers)")
	dispKey     = flag.Bool("key", false, "display color key")
	dispVersion = flag.Bool("version", false, "display version and exit")
	longFormat  = flag.Bool("numbers", false, "print numbers in addition to the chart")
)

// collect and parse data {{{
func readMemoryStats() MemData {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	res := MemData{}
	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		switch key {
		case "Buffers":
			res.Buffers = value
		case "Cached":
			res.Cached = value
		case "MemAvailable":
			res.MemAvailable = value
		case "MemFree":
			res.MemFree = value
		case "MemTotal":
			res.MemTotal = value
		case "Shmem":
			res.Shared += value
		case "SReclaimable":
			res.Cached += value

		case "SwapTotal":
			res.SwapTotal = value
		case "SwapFree":
			res.SwapFree = value
		}
	}

	res.Used = res.MemTotal - res.MemFree - res.Buffers - res.Cached
	res.SwapUsed = res.SwapTotal - res.SwapFree

	return res
}

func parseLine(raw string) (key string, value float64) {
	// remove "kB" at the end and remove spaces
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	// split by the semicolon
	keyValue := strings.Split(text, ":")
	return keyValue[0], toFloat(keyValue[1])
}

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

	data := readMemoryStats()

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
