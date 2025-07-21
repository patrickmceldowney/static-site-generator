package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/patrickmceldowney/static-site-generator/internal/builder"
)

func main() {
	build := flag.Bool("build", false, "Build the static site")
	watch := flag.Bool("watch", false, "Watch files and auto-rebuild on changes")
	flag.Parse()

	inputDir := "content"
	outputDir := "output"
	templateDir := "templates"

	if *build {
		err := builder.Build(inputDir, outputDir, templateDir)

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
