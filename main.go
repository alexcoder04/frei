package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// args {{{
var (
	Version   = "[built from source]"
	CommitSHA = ""

	format      = flag.String("format", "chart", "output format (chart/table/charttable/plain/json)")
	dispHuman   = flag.Bool("h", false, "display human-readable numbers")
	dispVersion = flag.Bool("version", false, "display version and exit")
)

// }}}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [options]\n\nOptions:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if *dispVersion {
		fmt.Printf("frei %s", Version)
		if CommitSHA != "" {
			fmt.Printf(" (commit %s)", CommitSHA)
		}
		fmt.Print("\n")
		os.Exit(0)
	}

	data, err := GetMemInfo()
	if err != nil {
		panic("Cannot get memory info")
	}

	switch *format {
	case "chart":
		printCharts(data)
	case "table":
		printTable(data, *dispHuman)
	case "charttable":
		printCharts(data)
		printTable(data, *dispHuman)
	case "plain":
		printPlainTextOutput(data, *dispHuman)
	case "json":
		printJsonOutput(data, *dispHuman)
	default:
		panic("Invalid format specified. Check -help for usage")
	}
}
