package output

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tislib/data-handler/pkg/model"
	"github.com/tislib/data-handler/pkg/service/annotations"
	"github.com/tislib/data-handler/pkg/types"
	"github.com/tislib/data-handler/pkg/util"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
)

type consoleWriter struct {
	writer   io.Writer
	describe bool
}

func (c consoleWriter) IsBinary() bool {
	return false
}

func (c consoleWriter) DescribeResource(resource *model.Resource) {
	const padding = 3
	w := tabwriter.NewWriter(c.writer, 0, 0, padding, ' ', 0)

	c.out(w, "Name: \t\t %s", resource.Name)
	c.out(w, "Namespace: \t\t %s", resource.Namespace)
	c.out(w, "Version: \t\t %d", resource.Version)
	c.out(w, "")

	c.out(w, "Source Config:")
	c.out(w, "  DataSource: \t\t %s", resource.SourceConfig.DataSource)
	c.out(w, "  Catalog: \t\t %s", resource.SourceConfig.Catalog)
	c.out(w, "  Entity: \t\t %s", resource.SourceConfig.Entity)
	c.out(w, "")

	if resource.AuditData != nil {
		c.out(w, "AuditData:")
		c.out(w, "  Created By: \t\t %s", resource.AuditData.CreatedBy)
		c.out(w, "  Created On: \t\t %s", resource.AuditData.CreatedOn.AsTime().String())
		c.out(w, "  Updated By: \t\t %s", resource.AuditData.UpdatedBy)
		c.out(w, "  Updated On: \t\t %s", resource.AuditData.UpdatedOn.AsTime().String())
		c.out(w, "")
	}

	if len(resource.Annotations) > 0 {
		c.out(w, "Annotations:")
		for key, value := range resource.Annotations {
			c.out(w, fmt.Sprintf("%s:\t%s", key, value))
		}
		c.out(w, "")
	}

	c.out(w, "Properties:")

	var data [][]string

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "Mapping", "Type", "Required", "Unique", "Primary", "Length", "Annotations"})
	c.configureTable(table)

	for _, item := range resource.Properties {

		typeStr := strings.ToLower(item.Type.String())[5:]

		data = append(data, []string{
			item.Name,
			item.Mapping,
			typeStr,
			strconv.FormatBool(item.Required),
			strconv.FormatBool(item.Unique),
			strconv.FormatBool(item.Primary),
			strconv.Itoa(int(item.Length)),
			annotations.ToString(item),
		})
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	_ = w.Flush()

	if len(resource.Indexes) > 0 {
		table = tablewriter.NewWriter(w)

		c.out(c.writer, "")
		c.out(c.writer, "Indexes:")

		data = [][]string{}
		table.SetHeader([]string{"IndexType", "Unique", "Properties", "Annotations"})
		c.configureTable(table)

		for _, item := range resource.Indexes {
			data = append(data, []string{
				item.IndexType.String(),
				strconv.FormatBool(item.Unique),
				strings.Join(util.ArrayMapToString(item.Properties, func(t *model.ResourceIndexProperty) string {
					return t.Name
				}), ", "),
				annotations.ToString(item),
			})
		}

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	}

	c.out(w, "")
	_ = w.Flush()
}

func (c consoleWriter) out(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format+"\n", a...)
}

func (c consoleWriter) configureTable(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
}

func (c consoleWriter) WriteResources(resources []*model.Resource) {
	if c.describe {
		for _, resource := range resources {
			c.DescribeResource(resource)
		}
	} else {
		c.ShowResourceTable(resources)
	}
}

func (c consoleWriter) ShowResourceTable(resources []*model.Resource) {
	var data [][]string

	table := tablewriter.NewWriter(c.writer)
	table.SetHeader([]string{"Id", "Name", "Namespace", "DataSource", "Catalog", "Entity", "Version"})
	c.configureTable(table)

	for _, item := range resources {
		data = append(data, []string{
			item.Id,
			item.Name,
			item.Namespace,
			item.SourceConfig.DataSource,
			item.SourceConfig.Catalog,
			item.SourceConfig.Entity,
			strconv.Itoa(int(item.Version)),
		})
	}

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func (c consoleWriter) WriteRecords(resource *model.Resource, total uint32, recordsChan chan *model.Record) {
	table := tablewriter.NewWriter(c.writer)
	columns := []string{"Id"}

	if !annotations.IsEnabled(resource, annotations.DisableVersion) {
		columns = append(columns, "version")
	}

	for _, prop := range resource.Properties {
		columns = append(columns, prop.Name)
	}

	table.SetHeader(columns)
	c.configureTable(table)

	var i = 0
	for item := range recordsChan {
		row := []string{
			item.Id,
		}

		if !annotations.IsEnabled(resource, annotations.DisableVersion) {
			row = append(row, strconv.Itoa(int(item.Version)))
		}

		for _, prop := range resource.Properties {
			typeHandler := types.ByResourcePropertyType(prop.Type)
			packedVal := item.Properties[prop.Name]

			if packedVal == nil {
				row = append(row, "Null")
			} else {
				value, err := typeHandler.UnPack(packedVal)

				check(err)
				valStr := typeHandler.String(value)

				row = append(row, valStr)
			}
		}
		i++

		table.Append(row)

		if i%1000 == 0 {
			table.Render()
		}
	}

	table.Render()
}
