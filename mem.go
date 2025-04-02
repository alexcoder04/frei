package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

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

func parseLine(raw string) (string, float64) {
	// remove "kB" at the end and remove spaces
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	// split by the semicolon
	keyValue := strings.Split(text, ":")
	return keyValue[0], toFloat(keyValue[1])
}

func GetMemInfo() (MemData, error) {
	res := MemData{}

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return res, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

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

	return res, nil
}
