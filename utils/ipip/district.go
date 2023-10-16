package ipip

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var DistrictEmptyError = fmt.Errorf("error[%s]", "not found.")

type DistrictDB struct {
	LastTime       time.Time
	ipIndex        [65536]uint32
	all, dataIndex []byte
	offset         int
	sync.RWMutex
}

type District struct {
	Country   string
	Province  string
	City      string
	District  string
	Code      string
	Radius    string
	Longitude string
	Latitude  string
}

func NewDistrictDB() *DistrictDB {

	return &DistrictDB{}
}

func (db *DistrictDB) Load(fileName string) error {
	db.Lock()
	defer db.Unlock()
	info, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	db.LastTime = info.ModTime()

	db.all, err = ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	db.offset = int(binary.BigEndian.Uint32(db.all[0:4]))

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			k := i*256 + j
			db.ipIndex[k] = binary.LittleEndian.Uint32(db.all[(k+1)*4 : (k+1)*4+4])
		}
	}

	db.dataIndex = make([]byte, db.offset-4)
	copy(db.dataIndex, db.all[4:db.offset-4])

	return nil
}

func (db *DistrictDB) _find(s string) ([]byte, error) {
	if len(db.all) == 0 {
		return nil, fmt.Errorf("error[%s]", "load ipdb file failed..")
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("error[%s]", "ip format error")
	}
	ip = ip.To4()
	var ipInt = binary.BigEndian.Uint32(ip)
	var prefix = int(ip[0])*256 + int(ip[1])
	var start = int(ipIndex[prefix])
	var maxValue = db.offset - 262144 - 4
	var indexOffset = -1
	var indexLength = -1

	for start = start*13 + 262144; start < maxValue; start += 13 {
		tmpInt := binary.BigEndian.Uint32(db.dataIndex[start : start+4])
		if ipInt >= tmpInt {
			if ipInt <= binary.BigEndian.Uint32(db.dataIndex[start+4:start+8]) {
				indexOffset = int(binary.LittleEndian.Uint32(db.dataIndex[start+8 : start+12]))
				indexLength = int(db.dataIndex[start+12])
				break
			}
		}
	}

	if indexOffset == -1 || indexLength == -1 {
		return nil, DistrictEmptyError
	}

	var area = make([]byte, indexLength)
	indexOffset = int(db.offset) + indexOffset - 262144
	copy(area, db.all[indexOffset:indexOffset+indexLength])

	return area, nil
}

func (db *DistrictDB) Find(s string) (District, error) {
	var err error
	var dis District

	db.RLock()
	defer db.RUnlock()

	bs, err := db._find(s)
	if err != nil {
		return dis, err
	}
	loc := strings.Split(string(bs), "\t")

	dis.Country = loc[0]
	dis.Province = loc[1]
	dis.City = loc[2]
	dis.District = loc[3]
	dis.Code = loc[4]
	dis.Radius = loc[5]
	dis.Longitude = loc[6]
	dis.Latitude = loc[7]

	return dis, err
}
