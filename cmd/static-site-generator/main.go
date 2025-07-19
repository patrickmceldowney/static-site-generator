package main

import (
	"fmt"
	"flag"
	"os"

	"gostatic/internal/builder"
)

function main() {
	build := flag.Bool("build", false, "Build the static site")
	watch := flag.Bool("watch", false, "Watch files and auto-rebuild on changes")
	flag.Parse()

	if *build {
		err := builder.Build()

		if err != nil {
			fmt.Println("Build failed: ", err)
			os.Exit(1)
		}

		fmt.Println("Site build successfully.")
	}

	if *watch {
		fmt.Println("Watch mode not implemented yet")
	}
}
