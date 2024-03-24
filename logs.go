package main

import (
	"fmt"

	"github.com/fatih/color"
)

func logSuccess(msg string) {
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	fmt.Println(green(fmt.Sprintf("✅ %s", msg)))
}

func logInfo(msg string) {
	blue := color.New(color.Bold, color.FgCyan).SprintFunc()
	fmt.Println(blue(fmt.Sprintf("🌱 %s", msg)))
}

func logBlocking(msg string) {
	magenta := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	fmt.Println(magenta(fmt.Sprintf("⌛️ %s", msg)))
}
func logBlockingDot() {
	magenta := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	fmt.Print(magenta("."))
}

func logError(msg string) {
	red := color.New(color.Bold, color.FgRed).SprintFunc()
	fmt.Println(red(fmt.Sprintf("❌ %s", msg)))
}

func logWarning(msg string) {
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	fmt.Println(yellow(fmt.Sprintf("⚠️  %s", msg)))
}
