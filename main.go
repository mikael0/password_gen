package main

import (
	"flag"
	"fmt"
	"os"
	"password_gen/m/v2/pkg/optimal"
)

func main() {
	dictPath := flag.String("dict", "", "Path to the dictionary file")
	mode := flag.String("mode", "fast", "Mode of operation (fast or optimal)")

	flag.Parse()

	if *dictPath == "" {
		fmt.Println("Error: --dict is required")
		flag.Usage()
		os.Exit(1)
	}

	if *mode != "fast" && *mode != "optimal" {
		fmt.Println("Error: --mode must be 'fast' or 'optimal'")
		flag.Usage()
		os.Exit(1)
	}

	switch *mode {
	case "optimal":
		optimal.Find(*dictPath, 20, 24, 4)
	}
}
