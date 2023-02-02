package excel

import (
	"github.com/xuri/excelize/v2"
	"strconv"
)

type ExportServer struct {
	Path         string
	SheetName    string //工作表名称
	Download     bool   //是否需要下载
	HttpDownload bool   //通过请求调用
	Header       Header //表头
	Data         []Data //数据内容
	excel        *excelize.File
	LineTitle    []Line
}
type Header struct {
	Width    float64 //列宽
	IsExport bool    //是否需要表头
	Height   float64 //行高
	Value    string  //表头值
}

type Data map[string]string

type Line struct {
	Title    string
	SubTitle []Line
}

func NewExport(export ExportServer) *ExportServer {
	service := excelize.NewFile()
	return &ExportServer{
		Path:      export.Path,
		Header:    export.Header,
		Data:      export.Data,
		excel:     service,
		SheetName: export.SheetName,
		LineTitle: export.LineTitle,
	}
}

func (e ExportServer) Export() (string, error) {
	index := 1
	sheetIndex, err := e.excel.NewSheet(e.SheetName)
	e.excel.SetActiveSheet(sheetIndex)
	if err != nil {
		return "", err
	}
	if e.Header.IsExport == true {
		index, err = e.SetHeader(index)
		if err != nil {
			return "", err
		}
	}
	index = e.SetLineTitle(index)
	e.SetValue(index)
	e.SaveFile()
	return e.Path, nil
}
func (e ExportServer) SaveFile() {
	err := e.excel.SaveAs(e.Path)
	if err != nil {
		return
	}
}
func (e ExportServer) SetValue(index int) {
	for k, value := range e.Data {
		err := e.excel.SetCellValue(e.SheetName, e.decimalToColumn(k+1)+strconv.Itoa(k+index), value)
		if err != nil {
			return
		}
	}
}
func (e ExportServer) SetLineTitle(index int) int {
	//合并单元格记录器
	start := 1
	subKey := 1
	for _, v := range e.LineTitle {
		//合并单元格
		end := len(v.SubTitle) + start - 1
		startLine := e.decimalToColumn(start)
		endLine := e.decimalToColumn(end)
		err := e.excel.MergeCell(e.SheetName, startLine+strconv.Itoa(index), endLine+strconv.Itoa(index))
		if err != nil {
			return 0
		}
		err = e.excel.SetCellValue(e.SheetName, startLine+strconv.Itoa(index), v.Title)
		if err != nil {
			return 0
		}
		start += len(v.SubTitle)
		//设置子菜单
		if len(v.SubTitle) != 0 {
			for _, subtitle := range v.SubTitle {
				err = e.excel.SetCellValue(e.SheetName, e.decimalToColumn(subKey)+strconv.Itoa(index+1), subtitle.Title)
				if err != nil {
					return 0
				}
				subKey++
			}
		} else {
			err = e.excel.SetCellValue(e.SheetName, e.decimalToColumn(subKey)+strconv.Itoa(index+1), v.Title)
			subKey++
			if err != nil {
				return 0
			}
		}
	}
	index += 2
	return index
}
func (e ExportServer) SetHeader(index int) (int, error) {
	MergeStart := e.decimalToColumn(1)
	MergeEnd := e.decimalToColumn(e.getCount())
	err := e.excel.MergeCell(e.SheetName, MergeStart+strconv.Itoa(index), MergeEnd+strconv.Itoa(index))
	if err != nil {
		return 0, err
	}
	if e.Header.Width != 0 {
		err = e.excel.SetColWidth(e.SheetName, "A", "A", e.Header.Width)
		if err != nil {
			return 0, err
		}
	}
	if e.Header.Height != 0 {
		err = e.excel.SetRowHeight(e.SheetName, index, e.Header.Height)
		if err != nil {
			return 0, err
		}
	}
	//设置单元格值
	err = e.excel.SetCellValue(e.SheetName, "A"+strconv.Itoa(index), e.Header.Value)
	if err != nil {
		return 0, err
	}
	index++
	return index, nil
}
func (e ExportServer) getCount() int {
	count := 0
	for _, l := range e.LineTitle {
		count += len(l.SubTitle)
	}
	return count
}
func (e ExportServer) decimalToColumn(num int) string {
	columns := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := len(columns)
	col := ""
	for num > 0 {
		index := num % length
		char := columns[length-1:]
		num = num/length - 1
		if index != 0 {
			char = columns[index-1 : index]
			num++
		}
		col = char + col
	}
	return col
}
