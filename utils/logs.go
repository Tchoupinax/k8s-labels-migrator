package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func LogSuccess(msg string) {
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	fmt.Println(green(fmt.Sprintf("✅ %s", msg)))
}

func LogInfo(msg string) {
	blue := color.New(color.Bold, color.FgCyan).SprintFunc()
	fmt.Println(blue(fmt.Sprintf("🌱 %s", msg)))
}

func LogBlocking(msg string) {
	magenta := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	fmt.Println(magenta(fmt.Sprintf("⌛️ %s", msg)))
}

func LogBlockingDot() {
	magenta := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	fmt.Print(magenta("."))
}

func LogError(msg string) {
	red := color.New(color.Bold, color.FgRed).SprintFunc()
	fmt.Println(red(fmt.Sprintf("❌ %s", msg)))
}

func LogWarning(msg string) {
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	fmt.Println(yellow(fmt.Sprintf("⚠️  %s", msg)))
}
