package sutils

// import (
// 	"comutils/orm"
// 	. "comutils/serverutils"
// 	"time"

// 	. "alpcserver"
// 	// . "comutils/serverutils"
// 	// "time"

// 	// "comutils/orm"

// 	cmap "github.com/orcaman/concurrent-map"
// )

// // const (
// // 	SCHEDULE_TYPE_TIME     = "schedule_time"
// // 	SCHEDULE_TYPE_INTERVAL = "schedule_interval"
// // )

// type TAccessCtrl struct {
// 	UserAccessMap    cmap.ConcurrentMap
// 	IpAccessMap      cmap.ConcurrentMap
// 	MethodMap        cmap.ConcurrentMap
// 	AccessRecordChan chan *TAccessRecord
// }

// type TAccessData struct {
// 	ReqTime time.Time
// }

// var AccessCtrl TAccessCtrl

// func AccessInit() {
// 	AccessCtrl.UserAccessMap = cmap.New()
// 	AccessCtrl.IpAccessMap = cmap.New()
// 	AccessCtrl.MethodMap = cmap.New()
// 	AccessCtrl.AccessRecordChan = make(chan *TAccessRecord, 8192)

// 	LoadAccess()

// 	GoRtMngr.NewLoopGoRoutine("AccessRecordRoutine_Loop", AccessRecordRoutine)
// }

// func LoadAccess() {

// 	var accessArr []TAccessConfig

// 	o := orm.NewOrm()

// 	_, err := o.Raw("select * from t_access_config where enabled=1").QueryRows(&accessArr)
// 	if err != nil && err != orm.ErrNoRows {
// 		Log.Error("load t_access_config config failed, err = %v", err)
// 		return
// 	}

// 	for _, acc := range accessArr {
// 		ac2 := acc
// 		AccessCtrl.MethodMap.Set(ac2.Method, &ac2)
// 	}

// 	var delkeys []string

// 	// 清理已删除的schedule
// 	methmap := AccessCtrl.MethodMap
// 	for item := range methmap.IterBuffered() {
// 		val := item.Val
// 		ac := val.(*TAccessConfig)
// 		found := false

// 		for _, a := range accessArr {
// 			if a.Method == ac.Method {
// 				found = true
// 				break
// 			}
// 		}

// 		if !found { // 如果找不到, 则是被删除或过期的schedule, 需要在缓存里一并删除
// 			delkeys = append(delkeys, ac.Method)
// 		}
// 	}

// 	for _, dk := range delkeys {
// 		AccessCtrl.MethodMap.Remove(dk)
// 	}
// }

// func AsyncAddAccessRecord(req *TMsgReq) {

// 	rec := &TAccessRecord{
// 		Method:   req.Header.Method,
// 		UserName: req.Header.UserName,
// 		Ip:       req.Header.ClientIp,
// 		AddTime:  time.Now(),
// 	}

// 	AccessCtrl.AccessRecordChan <- rec
// }

// func AccessRecordRoutine() {
// 	rec := <-AccessCtrl.AccessRecordChan

// 	method := rec.Method
// 	_, ok := AccessCtrl.MethodMap.Get(method)
// 	if !ok { // 找不到则无需保存检查
// 		return
// 	}

// 	o := orm.NewOrm()
// 	_, err := o.Insert(rec)
// 	if err != nil {
// 		Log.Warning("insert access record failed, rec = %v", FormatStruct(rec))
// 	}
// }

// func CleanAccessRecord() {

// 	o := orm.NewOrm()
// 	_, err := o.Raw("delete from t_access_record where add_time < date(now())").Exec()
// 	if err != nil {
// 		Log.Warning("delete t_access_record failed")
// 	}
// }

// func makeUserAccessKey(method string, username string) string {
// 	return method + "###" + username
// }

// func makeIpAccessKey(method string, ip string) string {
// 	return method + "###" + ip
// }

// func AccessCheckFreq(ar *TAccessConfig, method string, username string, ip string) TErr {

// 	freq := ar.Frequence
// 	amode := ar.AccessMode
// 	if amode == ACCESS_MODE_USER {
// 		key := makeUserAccessKey(method, username)
// 		obj, ok := AccessCtrl.UserAccessMap.Get(key)
// 		if ok {
// 			oad := obj.(*TAccessData)
// 			timenow := time.Now()
// 			dur := timenow.Sub(oad.ReqTime)
// 			freqdur := time.Duration(freq) * time.Second
// 			if dur < freqdur {
// 				return MakeErr(EC_USER_ACCESS_TOO_FAST)
// 			} else { // 刷新访问时间
// 				oad.ReqTime = timenow
// 				AccessCtrl.UserAccessMap.Set(key, oad)
// 			}
// 		} else {
// 			ad := &TAccessData{
// 				ReqTime: time.Now(),
// 			}
// 			AccessCtrl.UserAccessMap.Set(key, ad)
// 		}
// 	} else if amode == ACCESS_MODE_IP {
// 		key := makeIpAccessKey(method, ip)
// 		obj, ok := AccessCtrl.IpAccessMap.Get(key)
// 		if ok {
// 			oad := obj.(*TAccessData)
// 			timenow := time.Now()
// 			dur := timenow.Sub(oad.ReqTime)
// 			freqdur := time.Duration(freq) * time.Second
// 			if dur < freqdur {
// 				return MakeErr(EC_IP_ACCESS_TOO_FAST)
// 			} else { // 刷新访问时间
// 				oad.ReqTime = timenow
// 				AccessCtrl.UserAccessMap.Set(key, oad)
// 			}
// 		} else {
// 			ad := &TAccessData{
// 				ReqTime: time.Now(),
// 			}
// 			AccessCtrl.IpAccessMap.Set(key, ad)
// 		}
// 	} else {
// 		Log.Warning("Invalid access config, amode = %v", amode)
// 		return nil
// 	}

// 	return nil
// }

// func AccessCheckCycle(ar *TAccessConfig, method string, username string, ip string) TErr {

// 	acnt := ar.AccessCount
// 	acycle := ar.Cycle
// 	amode := ar.AccessMode
// 	if amode == ACCESS_MODE_USER {
// 		var cnt int64 = 0
// 		o := orm.NewOrm()
// 		err := o.Raw(`select count(0) as cnt
// 		from t_access_record
// 		where user_name=? and method=? and add_time > date_sub(now(), interval ? second)`,
// 			username, method, acycle).QueryRow(&cnt)
// 		if err != nil {
// 			Log.Warning("count t_access_record for user failed, username=%v, method=%v", username, method)
// 			return MakeErr(EC_DB_ERROR, "访问失败")
// 		}

// 		if cnt >= acnt {
// 			return MakeErr(EC_USER_ACCESS_TOO_FAST)
// 		}

// 	} else if amode == ACCESS_MODE_IP {
// 		var cnt int64 = 0
// 		o := orm.NewOrm()
// 		err := o.Raw(`select count(0) as cnt
// 		from t_access_record
// 		where ip=? and method=? and add_time > date_sub(now(), interval ? second)`,
// 			ip, method, acycle).QueryRow(&cnt)
// 		if err != nil {
// 			Log.Warning("count t_access_record for ip failed, ip=%v, method=%v", ip, method)
// 			return MakeErr(EC_DB_ERROR, "访问失败")
// 		}

// 		if cnt >= acnt {
// 			return MakeErr(EC_IP_ACCESS_TOO_FAST)
// 		}
// 	} else {
// 		Log.Warning("Invalid access config, amode = %v", amode)
// 		return nil
// 	}

// 	return nil
// }

// func AccessCheck(method string, username string, ip string) TErr {

// 	arobj, ok := AccessCtrl.MethodMap.Get(method)
// 	if !ok { // 找不到则无需检查
// 		return nil
// 	}

// 	ar := arobj.(*TAccessConfig)
// 	atype := ar.AccessType

// 	if atype == ACCESS_TYPE_FREQ {
// 		return AccessCheckFreq(ar, method, username, ip)
// 	} else if atype == ACCESS_TYPE_CYCLE {
// 		return AccessCheckCycle(ar, method, username, ip)
// 	} else {
// 		Log.Warning("Invalid access config, atype = %v", atype)
// 		return nil
// 	}

// 	return nil
// }
