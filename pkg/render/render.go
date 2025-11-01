package render

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// Render a formatted table with header and data to stdout using blueprint styling with row separators and text wrapping.
func RenderTable(w io.Writer, header []string, data ...[]string) error {
	table := tablewriter.NewTable(w,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal},
				ColMaxWidths: tw.CellWidth{Global: 25},
			},
		}),
	)
	table.Header(header)
	for _, row := range data {
		table.Append(row)
	}
	err := table.Render()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorFailedToRenderTable, err)
	}
	return nil
}
