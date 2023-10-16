package sutils

import (
	"io/ioutil"
	"strings"
)

func HandleJSONComments(origContent string) string {
	strArr := strings.Split(origContent, "\n")
	var prodStr string = ""
	for _, str := range strArr {
		str = strings.TrimSpace(str)
		// del remark
		strLen := len(str)
		quotes := false
		remark := false
		var idx int = 0
		for ; idx < strLen; idx++ {
			// 如果已经出现双引号，则跳過检查到出现下个双引号
			if str[idx] == '"' {
				quotes = !quotes
			} else if !quotes {
				if str[idx] == '/' && idx+1 < strLen {
					if str[idx+1] == '/' {
						remark = true
						break
					}
				}
			}
		}
		if remark {
			str = string((str[:idx]))
		}
		prodStr += str
	}
	return prodStr
}

type JSONFile struct {
	filepath     string
	origContent  []byte
	prodContent  []byte
	sOrigContent string
	sProdContent string
}

func InitJSONFile(filepath string) (*JSONFile, error) {
	jsonfile := &JSONFile{filepath: filepath}
	cont, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	jsonfile.origContent = cont[:]
	jsonfile.sOrigContent = string(jsonfile.origContent)
	delSpaceAndRemark(jsonfile)
	return jsonfile, nil
}

func delSpaceAndRemark(jsonfile *JSONFile) {
	strArr := strings.Split(jsonfile.sOrigContent, "\n")
	var prodStr string = ""
	for _, str := range strArr {
		str = strings.TrimSpace(str)
		// del remark
		strLen := len(str)
		quotes := false
		remark := false
		var idx int = 0
		for ; idx < strLen; idx++ {
			// 如果已经出现双引号，则跳過检查到出现下个双引号
			if str[idx] == '"' {
				quotes = !quotes
			} else if !quotes {
				if str[idx] == '/' && idx+1 < strLen {
					if str[idx+1] == '/' {
						remark = true
						break
					}
				}
			}
		}
		if remark {
			str = string((str[:idx]))
		}
		prodStr += str
	}
	jsonfile.sProdContent = prodStr
	jsonfile.prodContent = []byte(jsonfile.sProdContent)
}

func (this *JSONFile) GetFPath() string {
	return this.filepath
}

func (this *JSONFile) GetOrigContent() []byte {
	return this.origContent
}

func (this *JSONFile) GetProdContent() []byte {
	return this.prodContent
}

func (this *JSONFile) GetSOrigContent() string {
	return this.sOrigContent
}

func (this *JSONFile) GetSProdContent() string {
	return this.sProdContent
}
