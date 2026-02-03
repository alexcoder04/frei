package main

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
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

// }}}

// getTerminalWidth() {{{
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80 // default fallback
	}
	return width
}

// }}}

// toHumanStr() {{{
func toHumanStr(value float64, human bool) string {
	if !human {
		return fmt.Sprint(uint(value/1024), " M")
	}
	units := []string{"K", "M", "G", "T", "P", "E", "Z", "Y"}
	for _, unit := range units {
		if value < 1024 {
			return fmt.Sprintf("%.2f %s", value, unit)
		}
		value = value / 1024
	}
	return "Too much"
}

// }}}

// toFloat() {{{
func toFloat(raw string) float64 {
	if raw == "" {
		return 0
	}
	res, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return res
}

// }}}
