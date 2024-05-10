package main

import (
	"fmt"
	"os"

	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
	"github.com/fatih/color"
)

const Version string = "0.1.3"
const BuildDate string = "2024-05-10"

func cliCommandDisplayHelp(args []string) {
	displayVersion := utils.StringInSlice("-v", args[1:]) || utils.StringInSlice("--version", args[1:])

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
