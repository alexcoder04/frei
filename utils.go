package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
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

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// }}}

// getTerminalWidth() {{{
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
