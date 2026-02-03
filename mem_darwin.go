//go:build darwin

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

// GetMemInfo reads memory data from macOS system calls and vm_stat
func GetMemInfo() (MemData, error) {
	res := MemData{}

	// Get total memory via sysctl
	memSize, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return res, fmt.Errorf("failed to get hw.memsize: %w", err)
	}
	// Convert bytes to KB (to match Linux format)
	res.MemTotal = float64(memSize) / 1024

	// Get page size
	pageSize, err := unix.SysctlUint32("vm.pagesize")
	if err != nil {
		return res, fmt.Errorf("failed to get vm.pagesize: %w", err)
	}

	// Parse vm_stat output for detailed memory info
	vmStats, err := getVMStats()
	if err != nil {
		return res, fmt.Errorf("failed to get vm_stat: %w", err)
	}

	// Convert page counts to KB
	pageSizeKB := float64(pageSize) / 1024

	pagesFree := vmStats["Pages free"] * pageSizeKB
	pagesActive := vmStats["Pages active"] * pageSizeKB
	pagesInactive := vmStats["Pages inactive"] * pageSizeKB
	pagesSpeculative := vmStats["Pages speculative"] * pageSizeKB
	pagesWired := vmStats["Pages wired down"] * pageSizeKB
	pagesCompressed := vmStats["Pages occupied by compressor"] * pageSizeKB
	pagesPurgeable := vmStats["Pages purgeable"] * pageSizeKB

	res.MemFree = pagesFree
	res.MemUsed = pagesActive + pagesWired + pagesCompressed
	res.MemCached = pagesInactive + pagesSpeculative + pagesPurgeable
	res.MemBuffers = 0 // No direct equivalent on macOS
	res.MemAvailable = pagesFree + pagesInactive + pagesSpeculative + pagesPurgeable
	res.MemShared = 0 // Would require additional parsing

	// Get swap info via sysctl
	swapUsage, err := unix.SysctlRaw("vm.swapusage")
	if err == nil && len(swapUsage) >= 24 {
		// struct xsw_usage { uint64_t xsu_total, xsu_avail, xsu_used; }
		swapTotal := *(*uint64)(unsafe.Pointer(&swapUsage[0]))
		swapAvail := *(*uint64)(unsafe.Pointer(&swapUsage[8]))
		swapUsed := *(*uint64)(unsafe.Pointer(&swapUsage[16]))

		res.SwapTotal = float64(swapTotal) / 1024
		res.SwapFree = float64(swapAvail) / 1024
		res.SwapUsed = float64(swapUsed) / 1024
	}

	if res.MemTotal <= 0 {
		fmt.Println("Error: your total installed memory appears to be 0 or less.")
		os.Exit(1)
	}

	return res, nil
}

// getVMStats parses the output of vm_stat command
func getVMStats() (map[string]float64, error) {
	stats := make(map[string]float64)

	cmd := exec.Command("vm_stat")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// Skip header line
		if strings.HasPrefix(line, "Mach Virtual Memory Statistics") {
			continue
		}
		// Parse lines like: "Pages free:                             1234."
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])
		valueStr = strings.TrimSuffix(valueStr, ".")

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			continue
		}
		stats[key] = value
	}

	return stats, nil
}
