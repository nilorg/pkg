package binding

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/nilorg/pkg/gin/model"
	"github.com/nilorg/sdk/convert"
)

var (
	// JQueryDataTables gin binding jquery datatables model
	JQueryDataTables = JQueryDataTablesBinding{}
)

type JQueryDataTablesBinding struct{}

func (JQueryDataTablesBinding) Name() string {
	return "jqueryDataTables"
}

var ErrBindNotIsJQueryDataTablesParameters = errors.New("this struts not is JQueryDataTablesParameters")

func (JQueryDataTablesBinding) Bind(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	switch obj.(type) {
	case *model.JQueryDataTablesParameters:
		return deepCopy(obj, dataTablesParametersParseMap(values))
	default:
		return ErrBindNotIsJQueryDataTablesParameters
	}
}

// deepCopy 深度copy
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// DataTablesParametersParseMap 解析get参数
func dataTablesParametersParseMap(values url.Values) *model.JQueryDataTablesParameters {
	// 先把报文分类
	baseValues := make(map[string]string)
	columnValues := make(map[string]string)
	orderValues := make(map[string]string)
	for key, value := range values {
		if strings.Index(key, "columns") != -1 {
			columnValues[key] = value[0]
		} else if strings.Index(key, "order") != -1 {
			orderValues[key] = value[0]
		} else {
			baseValues[key] = value[0]
		}
	}

	// 处理列
	columns := make([]model.JQueryDataTablesColumns, 0)
	columnAttrLen := 6
	for i := 1; i <= len(columnValues)/columnAttrLen; i++ {
		for j := 1; j <= columnAttrLen; j++ {
			index := j - 1
			name := columnValues[fmt.Sprintf("columns[%d][name]", index)]
			data := columnValues[fmt.Sprintf("columns[%d][data]", index)]
			searchable := columnValues[fmt.Sprintf("columns[%d][searchable]", index)]
			orderable := columnValues[fmt.Sprintf("columns[%d][orderable]", index)]
			searchValue := columnValues[fmt.Sprintf("columns[%d][search][value]", index)]
			searchRegex := columnValues[fmt.Sprintf("columns[%d][search][regex]", index)]
			columns = append(columns, model.JQueryDataTablesColumns{
				Data:       data,
				Name:       name,
				Searchable: searchable == "true",
				Orderable:  orderable == "true",
				Search: model.JQueryDataTablesSearch{
					Value: searchValue,
					Regex: searchRegex == "true",
				},
			})
		}
	}
	// 排序
	orders := make([]model.JQueryDataTablesOrder, 0)
	orderAttrLen := 2
	for i := 1; i <= len(orderValues)/orderAttrLen; i++ {
		for j := 1; j <= orderAttrLen; j++ {
			index := j - 1
			columnIndex := convert.ToInt(columnValues[fmt.Sprintf("order[%d][column]", index)])
			dir := columnValues[fmt.Sprintf("order[%d][dir]", index)]
			orders = append(orders, model.JQueryDataTablesOrder{
				Column: columnIndex,
				Dir:    dir,
			})
		}
	}
	return &model.JQueryDataTablesParameters{
		Columns: columns,
		Order:   orders,
		Draw:    convert.ToInt(baseValues["draw"]),
		Start:   convert.ToInt64(baseValues["start"]),
		Length:  convert.ToInt64(baseValues["length"]),
	}
}
