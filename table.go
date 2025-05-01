package main

import (
	"fmt"
	"strings"
)

// printTable() {{{
func printTable(data MemData, human bool) {
	chartWidth := getTerminalWidth() - 2
	if chartWidth > 40 {
		chartWidth = 40
	}

	labels := [10]string{
		"Total",
		"Used",
		"Shared",
		"Buffers",
		"Cache",
		"Available",
		"Free",
		"Swap Total",
		"Swap Used",
		"Swap Free",
	}
	values := [10]string{
		toHumanStr(data.MemTotal, human),
		toHumanStr(data.MemUsed, human),
		toHumanStr(data.MemShared, human),
		toHumanStr(data.MemBuffers, human),
		toHumanStr(data.MemCached, human),
		toHumanStr(data.MemAvailable, human),
		toHumanStr(data.MemFree, human),
		toHumanStr(data.SwapTotal, human),
		toHumanStr(data.SwapUsed, human),
		toHumanStr(data.SwapFree, human),
	}
	colors := [10]string{
		"\033[1m",
		"\033[32m",
		"\033[35m",
		"\033[34m",
		"\033[33m",
		"",
		"",
		"",
		"",
		"",
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
		fmt.Printf(" │ %s%s\033[0m%s  │ %s %s │\n", colors[i], labels[i], strings.Repeat(" ", spacesNumberLabel), strings.Repeat(" ", spacesNumberValue), values[i])
	}

	// tail
	fmt.Println(" ╰─────────────┴" + strings.Repeat("─", chartWidth-16) + "╯")
}

// }}}
