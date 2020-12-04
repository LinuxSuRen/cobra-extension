package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"reflect"
	"strings"
)

// OutputOption represent the format of output
type OutputOption struct {
	Format string

	Columns        string
	WithoutHeaders bool
	Filter         []string

	Writer        io.Writer
	CellRenderMap map[string]RenderCell
}

// RenderCell render a specific cell in a table
type RenderCell = func(string) string

// FormatOutput is the interface of format output
type FormatOutput interface {
	Output(obj interface{}, format string) (data []byte, err error)
}

const (
	// JSONOutputFormat is the format of json
	JSONOutputFormat string = "json"
	// YAMLOutputFormat is the format of yaml
	YAMLOutputFormat string = "yaml"
	// TableOutputFormat is the format of table
	TableOutputFormat string = "table"
)

// Output print the object into byte array
// Deprecated see also OutputV2
func (o *OutputOption) Output(obj interface{}) (data []byte, err error) {
	switch o.Format {
	case JSONOutputFormat:
		return json.MarshalIndent(obj, "", "  ")
	case YAMLOutputFormat:
		return yaml.Marshal(obj)
	}

	return nil, fmt.Errorf("not support format %s", o.Format)
}

// OutputV2 print the data line by line
func (o *OutputOption) OutputV2(obj interface{}) (err error) {
	if o.Writer == nil {
		err = fmt.Errorf("no writer found")
		return
	}

	if len(o.Columns) == 0 {
		err = fmt.Errorf("no columns found")
		return
	}

	//cmd.logger.Debug("start to output", zap.Any("filter", o.Filter))
	obj = o.ListFilter(obj)

	var data []byte
	switch o.Format {
	case JSONOutputFormat:
		data, err = json.MarshalIndent(obj, "", "  ")
	case YAMLOutputFormat:
		data, err = yaml.Marshal(obj)
	case TableOutputFormat, "":
		table := CreateTableWithHeader(o.Writer, o.WithoutHeaders)
		table.AddHeader(strings.Split(o.Columns, ",")...)
		items := reflect.ValueOf(obj)
		for i := 0; i < items.Len(); i++ {
			table.AddRow(o.GetLine(items.Index(i))...)
		}
		table.Render()
	default:
		err = fmt.Errorf("not support format %s", o.Format)
	}

	if err == nil && len(data) > 0 {
		_, err = o.Writer.Write(data)
	}
	return
}

// ListFilter filter the data list by fields
func (o *OutputOption) ListFilter(obj interface{}) interface{} {
	if len(o.Filter) == 0 {
		return obj
	}

	elemType := reflect.TypeOf(obj).Elem()
	elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 10)
	items := reflect.ValueOf(obj)
	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)
		if o.Match(item) {
			elemSlice = reflect.Append(elemSlice, item)
		}
	}
	return elemSlice.Interface()
}

// Match filter an item
func (o *OutputOption) Match(item reflect.Value) bool {
	for _, f := range o.Filter {
		arr := strings.Split(f, "=")
		if len(arr) < 2 {
			continue
		}

		key := arr[0]
		val := arr[1]

		if !strings.Contains(ReflectFieldValueAsString(item, key), val) {
			return false
		}
	}
	return true
}

// GetLine returns the line of a table
func (o *OutputOption) GetLine(obj reflect.Value) []string {
	columns := strings.Split(o.Columns, ",")
	values := make([]string, 0)

	if o.CellRenderMap == nil {
		o.CellRenderMap = make(map[string]RenderCell, 0)
	}

	for _, col := range columns {
		cell := ReflectFieldValueAsString(obj, col)
		if renderCell, ok := o.CellRenderMap[col]; ok && renderCell != nil {
			cell = renderCell(cell)
		}

		values = append(values, cell)
	}
	return values
}

// SetFlag set flag of output format
// Deprecated, see also SetFlagWithHeaders
func (o *OutputOption) SetFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Format, "output", "o", TableOutputFormat,
		"Format the output, supported formats: table, json, yaml")
	cmd.Flags().BoolVarP(&o.WithoutHeaders, "no-headers", "", false,
		`When using the default output format, don't print headers (default print headers)`)
	cmd.Flags().StringArrayVarP(&o.Filter, "filter", "", []string{},
		"Filter for the list by fields")
}

// SetFlagWithHeaders set the flags of output
func (o *OutputOption) SetFlagWithHeaders(cmd *cobra.Command, headers string) {
	o.SetFlag(cmd)
	cmd.Flags().StringVarP(&o.Columns, "columns", "", headers,
		"The columns of table")
}
