package global

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	uuid "github.com/eahydra/gouuid"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	. "github.com/switfs/shadow-framework/logger"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	extra.RegisterFuzzyDecoders()
}

func ASSERT(exp bool, info ...string) { // 接受一个字符串参数
	if !exp {
		infostr := ""
		if len(info) > 0 {
			infostr = info[0]
		}
		Log.Errorf("ASSERT FAILED!\ninfo=[%v]\nstack = [%v]\n", infostr, string(debug.Stack()))
		panic("ASSERT FAILED")
	}
}

func SerializeToJson(st interface{}) string {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(st)

	return buf.String()
}

// instead of original unserialize function
func UnserializeFromJson(jsonstr string, st interface{}) error {
	d := json.NewDecoder(strings.NewReader(jsonstr))
	d.UseNumber()
	return d.Decode(st)
}

func UnserializeJson(jsonstr string, st interface{}) {
	d := json.NewDecoder(strings.NewReader(jsonstr))
	d.UseNumber()
	err := d.Decode(st)
	if err != nil {
		Log.Error("parse json failed, err = %v", err)
		panic(err)
	}
}

func UnserializeKeepNum(jsonstr string, st interface{}) error {
	d := json.NewDecoder(strings.NewReader(jsonstr))
	d.UseNumber()
	return d.Decode(st)
}

func UseMaxCpu() {
	// multiple cups using
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func FormatJsonStr(instr string) string {

	var out bytes.Buffer
	json2.Indent(&out, []byte(instr), "", "  ")

	return "\n" + out.String() + "\n"
}

func FormatStruct(inst interface{}) string {
	instr := SerializeToJson(inst)
	return FormatJsonStr(instr)
}

func GetProgName() string {
	fullPath, _ := exec.LookPath(os.Args[0])
	fname := filepath.Base(fullPath)

	return fname
}

func EncodeURI(data string) string {
	return url.QueryEscape(data)
}

func DecodeURI(data string) (string, error) {

	sdata, err := url.QueryUnescape(data)
	if err != nil {
		Log.Error("url.QueryUnescape err = %v", err)
		return "", err
	}

	return sdata, nil
}

func IsNullTime(t time.Time) bool {
	year := t.Year()
	if year == 1 {
		return true
	} else {
		return false
	}
}

func NowStr() string {
	timenow := time.Now().Format(TIME_FORMAT)
	return timenow
}

func NowStr2() string {
	timenow := time.Now().Format(TIME_FORMAT_COMPACT)
	return timenow
}

func TimeStr(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func TimeStr2(t time.Time) string {
	return t.Format(TIME_FORMAT_COMPACT)
}

func IsValidTimeStr(tstr string) bool {
	_, err := time.ParseInLocation(TIME_FORMAT, tstr, time.Local)
	if err != nil {
		return false
	}

	return true
}

func IsValidDateStr(tstr string) bool {
	_, err := time.ParseInLocation(DATE_FORMAT, tstr, time.Local)
	if err != nil {
		return false
	}

	return true
}

func DateStr(t time.Time) string {
	return t.Format(DATE_FORMAT)
}

func NowWithMs() string {
	timenow := time.Now().Format(TIME_FORMAT_WITH_MS_COMPACT)
	return timenow
}

func NowMs() string {
	timenow := time.Now().Format(TIME_FORMAT_WITH_MS)
	return timenow
}

func NowDateStr() string {
	timenow := time.Now().Format(DATE_FORMAT)
	return timenow
}

func NowDateStr2() string {
	timenow := time.Now().Format(DATE_FORMAT_COMPACT)
	return timenow
}

func String2Time(s string) (*time.Time, error) {

	loc, err := time.LoadLocation("Local")
	if err != nil {
		Log.Error("load location failed, err = %v", err)
		return nil, err
	}

	ltime, err := time.ParseInLocation(TIME_FORMAT, s, loc)
	if err != nil {
		Log.Error("parse in location failed, err = %v", err)
		return nil, err
	}

	return &ltime, nil
}

func Time2Date(t time.Time) time.Time {
	timedate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return timedate
}

func DayStartTimeAndNow() (*time.Time, *time.Time, string) {
	timenow := time.Now()

	loc, err := time.LoadLocation("Local")
	if err != nil {
		Log.Error("load location failed, err = %v", err)
		return nil, nil, ""
	}

	sstime := fmt.Sprintf("%04d%02d%02d", timenow.Year(), timenow.Month(), timenow.Day())
	starttime := sstime + "000000"

	ltime, err := time.ParseInLocation(TIME_FORMAT_COMPACT, starttime, loc)
	if err != nil {
		Log.Error("parse in location failed, err = %v", err)
		return nil, nil, ""
	}

	return &ltime, &timenow, sstime
}

func CalcTimeSecond(stime string) (*time.Time, int64, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, 0, err
	}

	ltime, err := time.ParseInLocation(TIME_FORMAT_COMPACT, stime, loc)
	if err != nil {
		return nil, 0, err
	}

	return &ltime, ltime.Unix(), nil
}

func IsHour(s string) bool {
	hourarr := map[string]string{
		"0":  "0",
		"1":  "1",
		"2":  "2",
		"3":  "3",
		"4":  "4",
		"5":  "5",
		"6":  "6",
		"7":  "7",
		"8":  "8",
		"9":  "9",
		"10": "10",
		"11": "11",
		"12": "12",
		"13": "13",
		"14": "14",
		"15": "15",
		"16": "16",
		"17": "17",
		"18": "18",
		"19": "19",
		"20": "20",
		"21": "21",
		"22": "22",
		"23": "23",
	}
	if _, ok := hourarr[s]; ok {
		return true
	}
	return false
}

func IsHour2(s string) bool {
	hourarr := map[string]string{
		"00": "00",
		"01": "01",
		"02": "02",
		"03": "03",
		"04": "04",
		"05": "05",
		"06": "06",
		"07": "07",
		"08": "08",
		"09": "09",
		"10": "10",
		"11": "11",
		"12": "12",
		"13": "13",
		"14": "14",
		"15": "15",
		"16": "16",
		"17": "17",
		"18": "18",
		"19": "19",
		"20": "20",
		"21": "21",
		"22": "22",
		"23": "23",
	}
	if _, ok := hourarr[s]; ok {
		return true
	}
	return false
}

func IsMinute(s string) bool {
	minutearr := map[string]string{
		"0":  "0",
		"1":  "1",
		"2":  "2",
		"3":  "3",
		"4":  "4",
		"5":  "5",
		"6":  "6",
		"7":  "7",
		"8":  "8",
		"9":  "9",
		"10": "10",
		"11": "11",
		"12": "12",
		"13": "13",
		"14": "14",
		"15": "15",
		"16": "16",
		"17": "17",
		"18": "18",
		"19": "19",
		"20": "20",
		"21": "21",
		"22": "22",
		"23": "23",
		"24": "24",
		"25": "25",
		"26": "26",
		"27": "27",
		"28": "28",
		"29": "29",
		"30": "30",
		"31": "31",
		"32": "32",
		"33": "33",
		"34": "34",
		"35": "35",
		"36": "36",
		"37": "37",
		"38": "38",
		"39": "39",
		"40": "40",
		"41": "41",
		"42": "42",
		"43": "43",
		"44": "44",
		"45": "45",
		"46": "46",
		"47": "47",
		"48": "48",
		"49": "49",
		"50": "50",
		"51": "51",
		"52": "52",
		"53": "53",
		"54": "54",
		"55": "55",
		"56": "56",
		"57": "57",
		"58": "58",
		"59": "59",
	}
	if _, ok := minutearr[s]; ok {
		return true
	}
	return false
}

func IsMinute2(s string) bool {
	minutearr := map[string]string{
		"00": "00",
		"01": "01",
		"02": "02",
		"03": "03",
		"04": "04",
		"05": "05",
		"06": "06",
		"07": "07",
		"08": "08",
		"09": "09",
		"10": "10",
		"11": "11",
		"12": "12",
		"13": "13",
		"14": "14",
		"15": "15",
		"16": "16",
		"17": "17",
		"18": "18",
		"19": "19",
		"20": "20",
		"21": "21",
		"22": "22",
		"23": "23",
		"24": "24",
		"25": "25",
		"26": "26",
		"27": "27",
		"28": "28",
		"29": "29",
		"30": "30",
		"31": "31",
		"32": "32",
		"33": "33",
		"34": "34",
		"35": "35",
		"36": "36",
		"37": "37",
		"38": "38",
		"39": "39",
		"40": "40",
		"41": "41",
		"42": "42",
		"43": "43",
		"44": "44",
		"45": "45",
		"46": "46",
		"47": "47",
		"48": "48",
		"49": "49",
		"50": "50",
		"51": "51",
		"52": "52",
		"53": "53",
		"54": "54",
		"55": "55",
		"56": "56",
		"57": "57",
		"58": "58",
		"59": "59",
	}
	if _, ok := minutearr[s]; ok {
		return true
	}
	return false
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

// s 中是否以 prefix 开始
func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

// s 中是否以 suffix 结尾
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func RemoveAllChar(s string, chs ...byte) string {
	var debug bytes.Buffer
	lenOfStr := len(s)

	for i := 0; i < lenOfStr; i++ {
		exist := false
		for j := range chs {
			if s[i] == chs[j] {
				exist = true
				break
			}
		}
		if !exist {
			debug.WriteByte(s[i])
		}
	}
	return debug.String()
}

func CatchException() {
	if err := recover(); err != nil {
		fullPath, _ := exec.LookPath(os.Args[0])
		fname := filepath.Base(fullPath)

		datestr := NowDateStr()
		outstr := fmt.Sprintf("\n======\n[%v] err=%v, stack=%v\n======\n", time.Now(), err, string(debug.Stack()))
		filename := "./log/panic_" + fname + datestr + ".log"
		f, err2 := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		ASSERT(err2 == nil)
		defer f.Close()
		f.WriteString(outstr)

		Log.Errorf("err = %v ", err)
	}
}

func CatchExceptionWithName(name string) {
	if err := recover(); err != nil {
		fullPath, _ := exec.LookPath(os.Args[0])
		fname := filepath.Base(fullPath)

		datestr := NowDateStr()
		outstr := fmt.Sprintf("\n======\n[%v] err=%v, name = %v, stack=%v\n======\n", time.Now(), err, name, string(debug.Stack()))
		filename := "./log/panic_" + fname + datestr + ".log"
		f, err2 := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		ASSERT(err2 == nil)
		defer f.Close()
		f.WriteString(outstr)

		// ioutil.WriteFile(filename, []byte(outstr), 0666) //写入文件(字节数组)

		Log.Errorf("err = %v ", err)
	}
}

func CatchExceptionWithHandler(f func(interface{}), para interface{}) {
	if err := recover(); err != nil {
		Log.Errorf("err = %v ", err)
		f(para)
	}
}

func GenUUID() string {
	return uuid.NewUUIDByRandom().String()
}
