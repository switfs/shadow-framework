package utils

import (
	"os"

	. "github.com/switfs/shadow-framework/logger"
	"github.com/switfs/shadow-framework/utils/ipip"
)

func IpLocInit(datpath string) {
	err := ipip.Load(datpath)
	if err != nil {
		Log.Errorf("load ip data failed, program will exit, err = %v", err)
		os.Exit(-1)
	}
}

func IpLocStr(ip string) string {
	ipinfo, err := ipip.FindLocation(ip)
	if err != nil || ipinfo == nil {
		Log.Warnf("Cannot recognize ip = %v", ip)
		return ip
	}

	loc := "[" + ipinfo.ISO2 + "]" + ipinfo.Country + ipinfo.Province + ipinfo.City + ipinfo.Isp
	return loc
}

func IpLocStr2(ip string) string {
	ipinfo, err := ipip.FindLocation(ip)
	if err != nil || ipinfo == nil {
		Log.Warnf("Cannot recognize ip = %v", ip)
		return ip
	}

	loc := ipinfo.Country + ipinfo.Province + ipinfo.City
	return loc
}

func IpGetCode(ip string) string {
	ipinfo, err := ipip.FindLocation(ip)
	if err != nil || ipinfo == nil {
		Log.Warnf("Cannot recognize ip = %v, err = %v", ip, err)
		return ip
	}

	return ipinfo.ISO2
}

func IpGetInfo(ip string) *ipip.Location {
	ipinfo, err := ipip.FindLocation(ip)
	if err != nil || ipinfo == nil {
		Log.Warnf("Cannot recognize ip = %v", ip)
		return nil
	}

	return ipinfo
}
