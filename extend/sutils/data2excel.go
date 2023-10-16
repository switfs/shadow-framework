package sutils

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"io"

	. "github.com/switfs/shadow-framework/extend/global"

	"github.com/switfs/shadow-framework/extend/sutils/xlsx"
	. "github.com/switfs/shadow-framework/logger"
)

var ExcelService TExcelService

type TExcelService struct {
	sheet     string
	CnNameMap map[string]string
	isRunning bool
	stop      bool
	mutex     sync.Mutex
}

func (this *TExcelService) CheckRunning() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.isRunning {
		return true
	}
	return false
}

func (this *TExcelService) setRunStatus(isRunning bool) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if isRunning && this.isRunning {
		return false
	} else {
		this.isRunning = isRunning
	}
	return true
}

func (this *TExcelService) Stop() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.isRunning {
		this.stop = true
	} else {
		this.stop = false
	}
}

// Call InitCnNameMap first
func (this *TExcelService) GetExecl(w io.Writer, o interface{}, colNames []string) error {
	if this.CheckRunning() {
		return errors.New("有其他导出请求还未完成,请稍后在试")
	}
	if !this.setRunStatus(true) {
		return errors.New("有其他导出请求还未完成,请稍后在试")
	}
	defer this.setRunStatus(false)

	Log.Info("Start ...")

	sh := this.CreateXlsx(colNames)
	ww := xlsx.NewWorkbookWriter(w)
	defer ww.Close()
	sw, err := ww.NewSheetWriter(sh)
	if err != nil {
		return err
	}
	header, err := this.CreateHeader(sh, sw, colNames)
	if err != nil {
		return err
	}
	err = this.InertSliceData2Xlsx(sh, sw, o, header)
	if err != nil {
		return err
	}

	// return sh.SaveToWriter(buf)
	return nil
}

// Call InitCnNameMap first
func (this *TExcelService) CreateFetchExecl(w io.Writer, colNames []string) (*xlsx.Sheet, *xlsx.WorkbookWriter, *xlsx.SheetWriter, error) {
	if this.CheckRunning() {
		return nil, nil, nil, errors.New("有其他导出请求还未完成,请稍后在试")
	}
	if !this.setRunStatus(true) {
		return nil, nil, nil, errors.New("有其他导出请求还未完成,请稍后在试")
	}
	Log.Info("Start ...")

	sh := this.CreateXlsx(colNames)
	ww := xlsx.NewWorkbookWriter(w)
	sw, err := ww.NewSheetWriter(sh)
	return sh, ww, sw, err
}

func (this *TExcelService) EndFetchExecl(ww *xlsx.WorkbookWriter) {
	ww.Close()
}

func (this *TExcelService) CloseFetchExecl() {
	this.setRunStatus(false)
}

func (this *TExcelService) InitCnNameMap(initData map[string]string) error {
	this.CnNameMap = initData
	if this.CnNameMap == nil {
		this.CnNameMap = map[string]string{}
	}
	return nil
}

func (this *TExcelService) CreateXlsx(colNames []string) *xlsx.Sheet {
	colSize := len(colNames)
	cols := make([]xlsx.Column, colSize)
	for i, cname := range colNames {
		colDefault := xlsx.Column{
			Name:  cname,
			Width: 10,
		}
		cols[i] = colDefault
	}
	sh := xlsx.NewSheetWithColumns(cols)
	return &sh
}

func (this *TExcelService) CreateHeader(sh *xlsx.Sheet, sw *xlsx.SheetWriter, colNames []string) (map[string]int, error) {
	if len(colNames) == 0 {
		return nil, errors.New("colName is null")
	}
	return this.setHeader(sh, sw, colNames)
}

func (this *TExcelService) InertSliceData2Xlsx(
	sh *xlsx.Sheet, sw *xlsx.SheetWriter,
	o interface{}, header map[string]int) error {
	k := reflect.TypeOf(o).Kind()
	v := reflect.ValueOf(o)
	if k == reflect.Ptr {
		v = v.Elem()
		k = v.Type().Kind()
	}
	// oldTime := time.Now()
	// Log.Info("start time: %v, dataCnt = %v", oldTime, v.Len())
	if k == reflect.Slice || k == reflect.Array {
		vLen := v.Len()
		for row := 0; row < vLen; row++ {
			if this.stop {
				this.stop = false
				break
			}
			this.setRow(sh, sw, v.Index(row), header)
			// if row%10000 == 0 {
			// 	Log.Info("vLen = %v, row = %v", vLen, row)
			// }
		}
	} else {
		return errors.New("data is not slice")
	}
	// nowTime := time.Now()
	// Log.Info("end Time: %v, diff: %v", nowTime, nowTime.Sub(oldTime))
	return nil
}

func (this *TExcelService) setHeader(
	sh *xlsx.Sheet, sw *xlsx.SheetWriter,
	colNames []string) (map[string]int, error) {
	header := map[string]int{}
	r := sh.NewRow()
	for i, str := range colNames {
		header[str] = i
		cnStr, ok := this.CnNameMap[str]
		if ok {
			str = cnStr
		}
		r.Cells[i] = xlsx.Cell{
			Type:  xlsx.CellTypeInlineString,
			Value: str,
		}
	}
	// sh.AppendRow(r)
	// return header, nil
	return header, sw.WriteRows([]xlsx.Row{r})
}

func (this *TExcelService) setRow(
	sh *xlsx.Sheet, sw *xlsx.SheetWriter,
	value reflect.Value, header map[string]int) error {
	r := sh.NewRow()
	kind := value.Type().Kind()
	// fmt.Printf("vd:%v, kind:%v\n", value, kind)
	if kind == reflect.Struct {
		t := value.Type()
		for i := 0; i < value.NumField(); i++ { //NumField取出这个接口所有的字段数量
			f := t.Field(i)                    //取得结构体的第i个字段
			data := value.Field(i).Interface() //取得字段的值
			dataType := f.Type
			fieldName := f.Name
			// Log.Info("%s: %v = %v\n", fieldName, dataType, data) //第i个字段的名称,类型,值
			this.setCellData(&r, header, fieldName, dataType, data)
		}
	} else if kind == reflect.Ptr {
		val := value.Elem()
		t := val.Type()
		for i := 0; i < val.NumField(); i++ { //NumField取出这个接口所有的字段数量
			f := val.Field(i)     //取得结构体的第i个字段
			data := f.Interface() //取得字段的值
			dataType := f.Type()
			fieldName := t.Field(i).Name
			// fmt.Printf("%s: %v = %v\n", fieldName, dataType, data) //第i个字段的名称,类型,值
			this.setCellData(&r, header, fieldName, dataType, data)
		}
	}
	// sh.AppendRow(r)
	// return nil
	return sw.WriteRows([]xlsx.Row{r})
}

func (this *TExcelService) setCellData(row *xlsx.Row,
	header map[string]int,
	fieldName string, dataType reflect.Type, data interface{},
) {
	col, ok := header[fieldName]
	// Log.Debug("col: %v,%v\n", col, ok)
	if ok {
		strData := ""
		if dataType.String() == "time.Time" {
			dt := data.(time.Time)
			if this.isNullTime(dt) {
			} else {
				strData = TimeStr(dt)
			}
		} else {
			str := fmt.Sprintf("%v", data)
			cnStr, ok := this.CnNameMap[str]
			if ok {
				str = cnStr
			}
			strData = str
		}
		row.Cells[col] = xlsx.Cell{
			Type:  xlsx.CellTypeInlineString,
			Value: strData,
		}
	}
}

func (this *TExcelService) getCellNumber(col, row int64) string {
	str := ""
	for i := int(col); i > 0; i /= 27 {
		h := i % 27
		if h == 0 {
			h += 1
		}
		str = string(int('@')+h) + str
	}
	return fmt.Sprintf("%s%d", str, row)
}

func (this *TExcelService) isNullTime(t time.Time) bool {
	year := t.Year()
	if year == 1 {
		return true
	} else {
		return false
	}
}
