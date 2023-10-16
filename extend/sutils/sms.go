package sutils

// import (
// 	. "alpcserver"
// 	. "comutils/serverutils"
// 	"comutils/serverutils/gorequest"
// 	"encoding/xml"
// 	"fmt"
// 	"net/url"
// 	"time"
// )

// const (
// 	SMS_SERV_COMPANY_HU_YI = "HU_YI"
// )

// var SmsService TSmsService

// type TSmsReqMsg struct {
// 	Target string `json:"target"`
// 	Detail string `json:"detail"`
// }

// type TIHuYiSubmitResult struct {
// 	XMLName xml.Name `xml:"SubmitResult" json:"XMLName"`
// 	Code    string   `xml:"code" json:"code"`
// 	Msg     string   `xml:"msg" json:"msg"`
// 	Smsid   string   `xml:"smsid" json:"smsid"`
// }

// type TSmsService struct {
// }

// func (this *TSmsService) GetSmsConfig() *TSmsApiConfig {
// 	var config TSmsApiConfig
// 	iConf := GetInternalConfByType(INTERNAL_CONFIG_SMSAPI, true)
// 	if iConf == nil || iConf.Config == "" {
// 		Log.Error("Get sms config error.")
// 		return nil
// 	}

// 	deConf, err := AesStrDecrypt(iConf.Config)
// 	if err != nil {
// 		Log.Error("sms config decrypt error.")
// 		return nil
// 	}

// 	err = UnserializeFromJson(deConf, &config)
// 	if err != nil {
// 		Log.Error("sms config UnserializeFromJson error.")
// 		return nil
// 	}

// 	return &config
// }

// func (this *TSmsService) CallSmsService(req *TUserApplyVerificationRecord) TErr {
// 	config := this.GetSmsConfig()
// 	if config == nil {
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "无法获取SmsApi配置")
// 	}

// 	var ec TErr
// 	switch config.ServCompany {
// 	case SMS_SERV_COMPANY_HU_YI:
// 		ec = this.callHuYiSmsService(req, config)
// 	default:
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "不支持该短信服务商")
// 	}
// 	return ec
// }

// func (this *TSmsService) callHuYiSmsService(req *TUserApplyVerificationRecord, config *TSmsApiConfig) TErr {
// 	urlPath := config.Url

// 	form := url.Values{}
// 	form.Add("account", config.UserName)
// 	form.Add("password", config.Password)
// 	form.Add("mobile", req.Destination)
// 	form.Add("content", req.Detail)
// 	form.Add("time", time.Now().Format("2006-01-02 15:04:05"))
// 	Log.Debug("url:%v, sendMsg:%v", urlPath, form.Encode())
// 	_, body, errs := gorequest.New().
// 		Post(urlPath).
// 		Send(form.Encode()).
// 		End()
// 	if errs != nil {
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "讯息发送失败")
// 	}

// 	Log.Debug("respMsg:%v", body)
// 	var submitResult TIHuYiSubmitResult
// 	err := xml.Unmarshal([]byte(body), &submitResult)
// 	if err != nil {
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "处理服务商回应消息失败")
// 	}

// 	req.ErrCode = submitResult.Code
// 	req.ErrMsg = submitResult.Msg
// 	req.SendMsgId = submitResult.Smsid

// 	if submitResult.Code != "2" {
// 		req.Status = APPLY_VERIFICATION_STATUS_FAIL
// 		emsg := fmt.Sprintf("短信错误代码:%v.", submitResult.Code)
// 		Log.Error("Send sms msg error...:%v", submitResult)
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, emsg)
// 	} else {
// 		req.Status = APPLY_VERIFICATION_STATUS_USABLE
// 	}

// 	return nil
// }
