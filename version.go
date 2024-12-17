package main

import (
	"fmt"
	"os"

	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
	"github.com/fatih/color"
)

var (
	version   string
	buildDate string
	commit    string
)

func cliCommandDisplayHelp(args []string) {
	displayVersion := utils.StringInSlice("-v", args[1:]) || utils.StringInSlice("--version", args[1:])

	if displayVersion {
		bold := color.New(color.Bold).SprintFunc()
		italic := color.New(color.Italic).SprintFunc()

		fmt.Println()
		fmt.Println(bold("⚡️ Kubernetes labels migrator"))
		fmt.Println()
		fmt.Println("build date: ", bold(buildDate))
		fmt.Println("version:    ", bold(version))
		fmt.Println("commit:     ", bold(commit))
		fmt.Println()
		fmt.Println(italic("Need help?"))
		fmt.Println(italic("https://github.com/Tchoupinax/k8s-labels-migrator/issues"))
		os.Exit(0)
	}
}
