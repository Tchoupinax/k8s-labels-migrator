package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func displaySummary(
	namespace string,
	deploymentName string,
	labelToChangeKey string,
	labelToChangeValue string,
	goalOfOperationIsToRemoveLabel bool,
) {
	fmt.Println()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Parameter", "Value"})
	t.AppendRows([]table.Row{{"Deployment name", deploymentName}})
	t.AppendRows([]table.Row{{"Namespace", namespace}})
	if goalOfOperationIsToRemoveLabel {
		t.AppendRows([]table.Row{{"Label", labelToChangeKey}})
	} else {
		t.AppendRows([]table.Row{{"Label", fmt.Sprintf("%s=%s", labelToChangeKey, labelToChangeValue)}})
	}
	t.AppendRows([]table.Row{{"Will the label be removed?", goalOfOperationIsToRemoveLabel}})
	t.SetStyle(table.StyleColoredBright)
	t.Render()
	fmt.Println()
}
