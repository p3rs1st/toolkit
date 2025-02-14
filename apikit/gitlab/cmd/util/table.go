package util

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func RenderTable[E any](cmd *cobra.Command, headers table.Row, datas []E, dataFunc func(E) []interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(cmd.OutOrStdout())
	t.AppendHeader(headers)
	for _, data := range datas {
		t.AppendRow(dataFunc(data))
	}
	t.Render()
}
