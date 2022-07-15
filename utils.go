package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

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

// toFloat {{{
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
