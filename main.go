package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type MemData struct {
	Buffers      float64
	Cached       float64
	MemAvailable float64
	MemFree      float64
	MemTotal     float64
	Shmem        float64
	SReclaimable float64

	SwapFree  float64
	SwapTotal float64
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

var (
	Version   = ""
	CommitSHA = ""

	longFormat  = flag.Bool("numbers", false, "print exact numbers")
	dispVersion = flag.Bool("version", false, "display version")
	dispKey     = flag.Bool("key", false, "display color key")
)

func getTerminalWidth() int {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}

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
			res.Shmem = value
		case "SReclaimable":
			res.SReclaimable = value

		case "SwapTotal":
			res.SwapTotal = value
		case "SwapFree":
			res.SwapFree = value
		}
	}
	return res
}

func parseLine(raw string) (key string, value float64) {
	// remove "kB" at the end and remove spaces
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	// split by the semicolon
	keyValue := strings.Split(text, ":")
	return keyValue[0], toFloat(keyValue[1])
}

func toFloat(raw string) float64 {
	if raw == "" {
		return 0
	}
	res, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		panic(err)
	}
	return res
}

func calcDrawData(data MemData, barWidth float64) DrawData {
	res := DrawData{}

	res.Free = uint((data.MemFree / data.MemTotal) * barWidth)
	res.Cache = uint((data.Cached / data.MemTotal) * barWidth)
	res.Shared = uint(((data.Shmem + data.SReclaimable) / data.MemTotal) * barWidth)
	res.Buffers = uint((data.Buffers / data.MemTotal) * barWidth)
	res.Used = uint(barWidth) - res.Free - res.Cache - res.Shared - res.Buffers

	res.SwapFree = uint((data.SwapFree / data.SwapTotal) * barWidth)
	res.SwapUsed = uint(barWidth) - res.SwapFree

	return res
}

func toMbStr(value float64) string {
	return fmt.Sprint(uint(value/1024), " M")
}

func drawCharts(data DrawData, chartWidth int) {
	fmt.Println(" ╭" + strings.Repeat("─", chartWidth-2) + "╮")
	fmt.Print(" │ Mem: ")
	fmt.Print("\033[33m" + strings.Repeat("█", int(data.Used)) + "\033[0m")
	fmt.Print("\033[35m" + strings.Repeat("█", int(data.Buffers)) + "\033[0m")
	fmt.Print("\033[36m" + strings.Repeat("█", int(data.Shared)) + "\033[0m")
	fmt.Print("\033[34m" + strings.Repeat("█", int(data.Cache)) + "\033[0m")
	fmt.Print(strings.Repeat(" ", int(data.Free)))
	fmt.Print(" │ \n")
	fmt.Println(" ├" + strings.Repeat("─", chartWidth-2) + "┤ ")

	fmt.Print(" │ Swp: ")
	fmt.Print("\033[31m" + strings.Repeat("█", int(data.SwapUsed)) + "\033[0m")
	fmt.Print("\033[0m" + strings.Repeat(" ", int(data.SwapFree)) + "\033[0m")
	fmt.Print(" │ \n")
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

func printKey(chartWidth int) {
	fmt.Println(" ╭" + strings.Repeat("─", chartWidth-2) + "╮")
	fmt.Println(" │ \033[33mUsed\033[0m, \033[35mBuffers\033[0m, \033[36mShared\033[0m, \033[34mCache\033[0m" + strings.Repeat(" ", chartWidth-32) + " │")
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

func printValues(data MemData, chartWidth int) {
	if chartWidth > 40 {
		chartWidth = 40
	}
	labels := [10]string{
		"Total",
		"Used",
		"Free",
		"Shared",
		"Buffers",
		"Cached",
		"Available",
		"Swap Total",
		"Swap Used",
		"Swap Free",
	}
	values := [10]string{
		toMbStr(data.MemTotal),
		toMbStr(data.MemTotal - data.MemFree - data.Buffers - data.Cached),
		toMbStr(data.MemFree),
		toMbStr(data.Shmem),
		toMbStr(data.Buffers),
		toMbStr(data.Cached + data.SReclaimable),
		toMbStr(data.MemAvailable),
		toMbStr(data.SwapTotal),
		toMbStr(data.SwapTotal - data.SwapFree),
		toMbStr(data.SwapFree),
	}
	fmt.Println(" ╭─────────────┬" + strings.Repeat("─", chartWidth-16) + "╮")
	for i := 0; i < len(labels); i++ {
		if i == len(labels)-3 {
			fmt.Println(" ├─────────────┼" + strings.Repeat("─", chartWidth-16) + "┤")
		}
		spacesNumberLabel := 10 - len(labels[i])
		spacesNumberValue := chartWidth - 19 - len(values[i])
		fmt.Println(" │",
			labels[i],
			strings.Repeat(" ", spacesNumberLabel),
			"│",
			strings.Repeat(" ", spacesNumberValue),
			values[i],
			"│")
	}
	fmt.Println(" ╰─────────────┴" + strings.Repeat("─", chartWidth-16) + "╯")
}

func main() {
	flag.Parse()

	if *dispVersion {
		if Version == "" {
			fmt.Print("frei [built from source]")
		} else {
			fmt.Print("frei v" + Version)
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
	drawData := calcDrawData(data, float64(barWidth))

	drawCharts(drawData, chartWidth)
	if *dispKey {
		printKey(chartWidth)
	}
	if *longFormat {
		printValues(data, chartWidth)
	}
}
