package main

import (
	"fmt"
	"strings"
)

func toHumanStr(value float64, human bool) string {
	if human {
		units := []string{"K", "M", "G", "T", "P", "E", "Z", "Y"}
		for _, unit := range units {
			if value < 1024 {
				return fmt.Sprint(uint(value), " ", unit)
			}
			value = value / 1024
		}
	}
	return fmt.Sprint(uint(value/1024), " M")
}

func calcDrawData(data MemData, barWidth int) DrawData {
	res := DrawData{}

	res.Free = uint((data.MemFree / data.MemTotal) * float64(barWidth))
	res.Cache = uint((data.Cached / data.MemTotal) * float64(barWidth))
	res.Shared = uint(((data.Shmem + data.SReclaimable) / data.MemTotal) * float64(barWidth))
	res.Buffers = uint((data.Buffers / data.MemTotal) * float64(barWidth))
	// this is the only value that can become negative due to rounding errors
	if res.Free+res.Cache+res.Shared+res.Buffers >= uint(barWidth) {
		res.Used = 0
	} else {
		res.Used = uint(barWidth) - res.Free - res.Cache - res.Shared - res.Buffers
	}

	res.SwapFree = uint((data.SwapFree / data.SwapTotal) * float64(barWidth))
	res.SwapUsed = uint(barWidth) - res.SwapFree

	res.Free = res.Free - ((res.Free + res.Cache + res.Shared + res.Buffers + res.Used) - uint(barWidth))

	return res
}

func printCharts(data DrawData, chartWidth int) {
	// head
	fmt.Println(" ╭" + strings.Repeat("─", chartWidth-2) + "╮")
	// memory
	fmt.Print(" │ Mem: ")
	fmt.Print("\033[33m" + strings.Repeat("█", int(data.Used)) + "\033[0m")
	fmt.Print("\033[35m" + strings.Repeat("█", int(data.Buffers)) + "\033[0m")
	fmt.Print("\033[36m" + strings.Repeat("█", int(data.Shared)) + "\033[0m")
	fmt.Print("\033[34m" + strings.Repeat("█", int(data.Cache)) + "\033[0m")
	fmt.Print(strings.Repeat(" ", int(data.Free)))
	fmt.Print(" │ \n")
	// delimeter
	fmt.Println(" ├" + strings.Repeat("─", chartWidth-2) + "┤ ")
	// swap
	fmt.Print(" │ Swp: ")
	fmt.Print("\033[31m" + strings.Repeat("█", int(data.SwapUsed)) + "\033[0m")
	fmt.Print("\033[0m" + strings.Repeat(" ", int(data.SwapFree)) + "\033[0m")
	fmt.Print(" │ \n")
	// tail
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

func printKey(chartWidth int) {
	fmt.Println(" ╭" + strings.Repeat("─", chartWidth-2) + "╮")
	fmt.Println(" │ \033[33mUsed\033[0m, \033[35mBuffers\033[0m, \033[36mShared\033[0m, \033[34mCache\033[0m" + strings.Repeat(" ", chartWidth-32) + " │")
	fmt.Println(" ╰" + strings.Repeat("─", chartWidth-2) + "╯")
}

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
		toHumanStr(data.MemTotal-data.MemFree-data.Buffers-data.Cached, human),
		toHumanStr(data.MemFree, human),
		toHumanStr(data.Shmem, human),
		toHumanStr(data.Buffers, human),
		toHumanStr(data.Cached+data.SReclaimable, human),
		toHumanStr(data.MemAvailable, human),
		toHumanStr(data.SwapTotal, human),
		toHumanStr(data.SwapTotal-data.SwapFree, human),
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
