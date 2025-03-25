package util

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func RenderTable[E any](cmd *cobra.Command, headers table.Row, datas []E, dataFunc func(E) []interface{}) {
	writer := table.NewWriter()
	writer.SetOutputMirror(cmd.OutOrStdout())
	writer.AppendHeader(headers)

	for _, data := range datas {
		writer.AppendRow(dataFunc(data))
	}

	writer.Render()
}
