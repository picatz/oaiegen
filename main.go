package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/picatz/oaiegen/pkg/eval"
)

// Note: currently only supports the "evals.elsuite.basic.match:Match"
//       evaluation type, and only one at a time.

// CLI flags.
var (
	evalfile string
	workdir  string
)

func init() {
	// Set the flags.
	flag.StringVar(&evalfile, "file", "", "HCL file containing the evaluation data")
	flag.StringVar(&workdir, "workdir", "", "the working directory to write files to")

	// Parse the flags.
	flag.Parse()
}

func main() {
	// Ensure all flags are set.
	if evalfile == "" || workdir == "" {
		fmt.Println("missing required flags")
		flag.Usage()
		os.Exit(1)
	}

	// Ensure the working directory exists.
	if _, err := os.Stat(workdir); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read the evaluation data.
	matches, err := eval.ReadHCL(evalfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write the evaluation data to a file.
	err = matches.WriteFile(filepath.Join(workdir, "samples.jsonl"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
