package main

import (
	"fmt"
	"strings"
)

func printPlainTextOutput(data MemData, human bool) {
	labels := []string{
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
	values := []string{
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

	for i := 0; i < len(labels); i++ {
		fmt.Printf("%s:%s%s\n", labels[i], strings.Repeat(" ", (11-len(labels[i]))+(8-len(values[i]))), values[i])
	}
}
