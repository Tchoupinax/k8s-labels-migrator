package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const Version string = "0.1.0"
const BuildDate string = "2024-03-24"

func cliCommandDisplayHelp(args []string) {
	displayVersion := stringInSlice("-v", args[1:]) || stringInSlice("--version", args[1:])

	if displayVersion {
		bold := color.New(color.Bold).SprintFunc()
		italic := color.New(color.Italic).SprintFunc()

		fmt.Println()
		fmt.Println(bold("⚡️ Kubernetes labels migrator"))
		fmt.Println()
		fmt.Println("build date: ", bold(BuildDate))
		fmt.Println("version:         ", bold(Version))
		fmt.Println()
		fmt.Println(italic("Need help?"))
		fmt.Println(italic("https://github.com/Tchoupinax/k8s-labels-migrator/issues"))
		os.Exit(0)
	}
}
