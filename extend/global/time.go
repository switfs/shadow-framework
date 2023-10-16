package global

import (
	"time"
)

const (
	TIME_FORMAT_WITH_MS         = "2006-01-02 15:04:05.000"
	TIME_FORMAT                 = "2006-01-02 15:04:05"
	TIME_FORMAT_WITH_T          = "2006-01-02T15:04:05"
	TIME_FORMAT1                = "02-01-2006 15:04:00" // "dd-MM-yyyy HH:mm:ss" 秒数为 00
	TIME_FORMAT2                = "02-01-2006 15:04:05" // "dd-MM-yyyy HH:mm:ss"
	TIME_FORMAT3                = "01/02/2006 15:04:05" // MM/DD/YYYY HH:MM:SS
	TIME_FORMAT_COMPACT         = "20060102150405"
	TIME_FORMAT_COMPACT1        = "060102150405"
	TIME_FORMAT_WITH_MS_COMPACT = "20060102150405.000"
	DATE_FORMAT                 = "2006-01-02"
	DATE_FORMAT1                = "2006/01/02"
	DATE_FORMAT_COMPACT         = "20060102"
	DATE_FORMAT_COMPACT1        = "20060102 15:04:05"
	MONTH_FORMAT                = "2006-01"
)

const (
	TIME_LOC_ASIA_SHANGHAI    = "Asia/Shanghai"    //+0800
	TIME_LOC_ASIA_TAIPEI      = "Asia/Taipei"      //+0800
	TIME_LOC_AMERICA_NEW_YORK = "America/New_York" //-0400(夏) -0500(冬）
	TIME_LOC_UTC              = "UTC"
)

// Millisecond timestrap
func GetTimeOfMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// Nanosecond timestrap
func GetTimeOfNs() int64 {
	return time.Now().UnixNano()
}

// Second timestrap
func GetTimeOfS() int64 {
	return time.Now().Unix()
}

// string time to time.Time in location
func GetTimeFromFormat(layout string, timeStr string, location string) (time.Time, error) {
	var timestamp time.Time
	local, err := time.LoadLocation(location)
	if nil != err {
		timestamp, err = time.Parse(layout, timeStr)
	} else {
		timestamp, err = time.ParseInLocation(layout, timeStr, local)
	}

	return timestamp, err
}

// time.Time to string time in location
func GetTimeStrFromFormat(layout string, timeClass time.Time, location string) string {
	local, err := time.LoadLocation(location)
	if nil != err {
		return timeClass.Format(layout)
	}
	return timeClass.In(local).Format(layout)
}

// change time zone
func GetTimeByChangeTimeZone(layout string, timeStr string, srcTimeZone string, destTimeZone string) (time.Time, error) {
	var timestamp time.Time

	localTimeStr, err := time.LoadLocation(srcTimeZone)
	if err != nil {
		return timestamp, err
	}
	timestamp, err = time.ParseInLocation(layout, timeStr, localTimeStr)
	if err != nil {
		return timestamp, err
	}

	local, err := time.LoadLocation(destTimeZone)
	if err != nil {
		return timestamp, err
	}

	return timestamp.In(local), nil
}
