package main

import (
	"fmt"
	"os"
	"strings"
)

// calcDrawData() {{{
func calcDrawData(data MemData, barWidth int) DrawData {
	res := DrawData{}

	// this is a straightforward conversion of kb sizes into terminal block widths
	res.Used = int((data.MemUsed / data.MemTotal) * float64(barWidth))
	res.Buffers = int((data.MemBuffers / data.MemTotal) * float64(barWidth))
	res.Shared = int((data.MemShared / data.MemTotal) * float64(barWidth))
	res.Cache = int((data.MemCached / data.MemTotal) * float64(barWidth))
	res.Free = barWidth - res.Used - res.Buffers - res.Shared - res.Cache

	res.SwapFree = int((data.SwapFree / data.SwapTotal) * float64(barWidth))
	res.SwapUsed = barWidth - res.SwapFree

	// we want to be safe from weird edge cases (only god knows)
	if res.Free < 0 {
		if res.Used+res.Shared+res.Buffers > barWidth {
			fmt.Println("Something very weird just happened.")
			fmt.Println("You appear to be using way more memory that is installed on your system.")
			os.Exit(1)
		}
		res.Free = 0
		res.Cache = barWidth - res.Used - res.Shared - res.Buffers
	}

	return res
}

// }}}

// drawBar() {{{
// draw a number if it fits in, otherwise empty bar
func drawBar(width int, number float64) string {
	human := toHumanStr(number, true)
	if len(human)+2 > int(width) {
		return strings.Repeat(" ", width)
	}
	return fmt.Sprintf(" %s%s", human, strings.Repeat(" ", width-1-len(human)))
}

// }}}

// printCharts() {{{
func printCharts(stats MemData) {
	chartWidth := getTerminalWidth() - 2
	barWidth := chartWidth - 4 - 5
	data := calcDrawData(stats, barWidth)

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
	if stats.SwapTotal > 0 {
		fmt.Printf("\033[41m%s\033[0m", drawBar(data.SwapUsed, stats.SwapUsed))
		fmt.Printf("\033[0m%s\033[0m", strings.Repeat(" ", int(data.SwapFree)))
	} else {
		fmt.Print(strings.Repeat(" ", barWidth))
	}
	fmt.Print(" │ \n")

	// tail
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

// }}}
