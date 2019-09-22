package printer

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// PrintData todo
func PrintData(out io.Writer, headers []string, data [][]string) {
	table := tablewriter.NewWriter(out)
	table.SetHeader(headers)
	//table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	//table.SetCenterSeparator("")
	//table.SetColumnSeparator("")
	//table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	//table.SetHeaderLine(false)
	table.AppendBulk(data)
	table.Render()
}
