package main

import (
	"encoding/json"
	"fmt"
)

func printJsonOutput(data MemData, human bool) {
	dataMap := map[string]string{
		"Total":      toHumanStr(data.MemTotal, human),
		"Used":       toHumanStr(data.MemUsed, human),
		"Shared":     toHumanStr(data.MemShared, human),
		"Buffers":    toHumanStr(data.MemBuffers, human),
		"Cache":      toHumanStr(data.MemCached, human),
		"Available":  toHumanStr(data.MemAvailable, human),
		"Free":       toHumanStr(data.MemFree, human),
		"Swap Total": toHumanStr(data.SwapTotal, human),
		"Swap Used":  toHumanStr(data.SwapUsed, human),
		"Swap Free":  toHumanStr(data.SwapFree, human),
	}

	dataJson, err := json.MarshalIndent(dataMap, "", " ")
	if err != nil {
		panic("Error marshaling json")
	}

	fmt.Println(string(dataJson))
}
