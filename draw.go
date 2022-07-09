package main

import (
	"fmt"
	"strings"
)

// calcDrawData() {{{
func calcDrawData(data MemData, barWidth int) DrawData {
	res := DrawData{}

	res.Free = uint((data.MemFree / data.MemTotal) * float64(barWidth))
	res.Cache = uint((data.Cached / data.MemTotal) * float64(barWidth))
	res.Shared = uint(((data.Shared) / data.MemTotal) * float64(barWidth))
	res.Buffers = uint((data.Buffers / data.MemTotal) * float64(barWidth))

	// this is the only value that can become negative due to rounding errors
	if res.Free+res.Cache+res.Shared+res.Buffers >= uint(barWidth) {
		res.Used = 0
	} else {
		res.Used = uint(barWidth) - res.Free - res.Cache - res.Shared - res.Buffers
	}

	res.Free = res.Free - ((res.Free + res.Cache + res.Shared + res.Buffers + res.Used) - uint(barWidth))

	res.SwapFree = uint((data.SwapFree / data.SwapTotal) * float64(barWidth))
	res.SwapUsed = uint(barWidth) - res.SwapFree

	return res
}

// }}}

// drawBar() {{{
func drawBar(width uint, number float64) string {
	human := toHumanStr(number, true)
	if len(human)+2 > int(width) {
		return strings.Repeat(" ", int(width))
	}
	return fmt.Sprintf(" %s%s", human, strings.Repeat(" ", int(width)-1-len(human)))
}

// }}}

// printCharts() {{{
func printCharts(data DrawData, chartWidth int, stats MemData) {
	// head
	fmt.Println(" ╭" + strings.Repeat("─", chartWidth-2) + "╮")

	// memory
	fmt.Print(" │ Mem: ")

	fmt.Printf("\033[43m%s\033[0m", drawBar(data.Used, stats.Used))
	fmt.Printf("\033[45m%s\033[0m", drawBar(data.Buffers, stats.Buffers))
	fmt.Printf("\033[46m%s\033[0m", drawBar(data.Shared, stats.Shared))
	fmt.Printf("\033[44m%s\033[0m", drawBar(data.Cache, stats.Cached))

	fmt.Print(strings.Repeat(" ", int(data.Free)))
	fmt.Print(" │ \n")

	// delimeter
	fmt.Println(" ├" + strings.Repeat("─", chartWidth-2) + "┤ ")

	// swap
	fmt.Print(" │ Swp: ")
	fmt.Printf("\033[41m%s\033[0m", drawBar(data.SwapUsed, stats.SwapUsed))
	fmt.Printf("\033[0m%s\033[0m", strings.Repeat(" ", int(data.SwapFree)))
	fmt.Print(" │ \n")

	// tail
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

// }}}

// printKey() {{{
func printKey(chartWidth int) {
	fmt.Println(" ╭" + strings.Repeat("─", chartWidth-2) + "╮")
	fmt.Println(" │ \033[33mUsed\033[0m, \033[35mBuffers\033[0m, \033[36mShared\033[0m, \033[34mCache\033[0m" + strings.Repeat(" ", chartWidth-32) + " │")
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

// }}}

// printNumbers() {{{
func printNumbers(data MemData, chartWidth int, human bool) {
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
		toHumanStr(data.MemTotal, human),
		toHumanStr(data.Used, human),
		toHumanStr(data.MemFree, human),
		toHumanStr(data.Shared, human),
		toHumanStr(data.Buffers, human),
		toHumanStr(data.Cached, human),
		toHumanStr(data.MemAvailable, human),
		toHumanStr(data.SwapTotal, human),
		toHumanStr(data.SwapUsed, human),
		toHumanStr(data.SwapFree, human),
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

// }}}
