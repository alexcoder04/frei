package main

import (
	"fmt"
	"strings"
)

// calcDrawData() {{{
func calcDrawData(data MemData, barWidth int) DrawData {
	res := DrawData{}

	res.Used = uint((data.MemUsed / data.MemTotal) * float64(barWidth))
	res.Buffers = uint((data.MemBuffers / data.MemTotal) * float64(barWidth))
	res.Shared = uint((data.MemShared / data.MemTotal) * float64(barWidth))
	res.Cache = uint((data.MemCached / data.MemTotal) * float64(barWidth))
	res.Free = uint(barWidth) - res.Used - res.Buffers - res.Shared - res.Cache

	res.SwapFree = uint((data.SwapFree / data.SwapTotal) * float64(barWidth))
	res.SwapUsed = uint(barWidth) - res.SwapFree

	return res
}

// }}}

// drawBar() {{{
// draw a number if it fits in, otherwise empty bar
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

	fmt.Printf("\033[42;30m%s\033[0m", drawBar(data.Used, stats.MemUsed))
	fmt.Printf("\033[45;30m%s\033[0m", drawBar(data.Shared, stats.MemShared))
	fmt.Printf("\033[44;30m%s\033[0m", drawBar(data.Buffers, stats.MemBuffers))
	fmt.Printf("\033[43;30m%s\033[0m", drawBar(data.Cache, stats.MemCached))

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
	fmt.Println(" │ \033[32mUsed\033[0m, \033[35mShared\033[0m, \033[34mBuffers\033[0m, \033[33mCache\033[0m" + strings.Repeat(" ", chartWidth-32) + " │")
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

// }}}

// printNumbers() {{{
func printNumbers(data MemData, chartWidth int, human bool) {
	if chartWidth > 40 {
		chartWidth = 40
	}

	labels := [9]string{
		"Total",
		"Used",
		"Shared",
		"Buffers",
		"Cache",
		"Available",
		"Swap Total",
		"Swap Used",
		"Swap Free",
	}
	values := [9]string{
		toHumanStr(data.MemTotal, human),
		toHumanStr(data.MemUsed, human),
		toHumanStr(data.MemShared, human),
		toHumanStr(data.MemBuffers, human),
		toHumanStr(data.MemCached, human),
		toHumanStr(data.MemAvailable, human),
		toHumanStr(data.SwapTotal, human),
		toHumanStr(data.SwapUsed, human),
		toHumanStr(data.SwapFree, human),
	}

	// head
	fmt.Println(" ╭─────────────┬" + strings.Repeat("─", chartWidth-16) + "╮")
	// body
	for i := 0; i < len(labels); i++ {
		// separator before swap
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
	// tail
	fmt.Println(" ╰─────────────┴" + strings.Repeat("─", chartWidth-16) + "╯")
}

// }}}
