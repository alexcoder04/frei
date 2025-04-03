package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// parseLine() {{{
func parseLine(raw string) (string, float64) {
	// remove "kB" at the end and remove spaces
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	// split by the semicolon
	keyValue := strings.Split(text, ":")
	return keyValue[0], toFloat(keyValue[1])
}

// }}}

// GetMemInfo() {{{
// read memory data from /proc/meminfo and convert it into categories that make sense
func GetMemInfo() (MemData, error) {
	res := MemData{}

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return res, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// we use memory calculation mechanism from htop/procps
	// https://github.com/htop-dev/htop/blob/3.4.0/linux/LinuxMachine.c#L200
	// https://gitlab.com/procps-ng/procps/-/blob/v4.0.5/library/meminfo.c

	var sreclaimable float64
	var cached float64
	var shmem float64
	var free float64

	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		switch key {
		case "MemTotal":
			res.MemTotal = value

		case "Buffers":
			res.MemBuffers = value
		case "MemAvailable":
			res.MemAvailable = value

		case "Cached":
			cached = value
		case "SReclaimable":
			sreclaimable = value
		case "Shmem":
			shmem = value
		case "MemFree":
			free = value

		case "SwapTotal":
			res.SwapTotal = value
		case "SwapFree":
			res.SwapFree = value
		}
	}

	res.MemShared = shmem
	res.MemCached = cached + sreclaimable - shmem

	usedDiff := free + cached + sreclaimable + res.MemBuffers
	if res.MemTotal >= usedDiff {
		res.MemUsed = res.MemTotal - usedDiff
	} else {
		res.MemUsed = res.MemTotal - free
	}

	res.SwapUsed = res.SwapTotal - res.SwapFree

	// just to be on the safe side
	if res.MemTotal <= 0 {
		fmt.Println("Error: your total installed memory appears to be 0 or less.")
		os.Exit(1)
	}

	return res, nil
}

// }}}
