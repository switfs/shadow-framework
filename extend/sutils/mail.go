package sutils

// import (
// 	. "alpcserver"
// 	. "comutils/serverutils"
// 	"net/smtp"
// 	"strings"
// )

// const (
// 	MAIL_SERV_COMPANY_GMAIL = "GMAIL"
// )

// var MailService TMailService

// func SendMailForSmtp(user, password, host, to, subject, body, mailtype string) error {
// 	hp := strings.Split(host, ":")
// 	auth := smtp.PlainAuth("", user, password, hp[0])
// 	var content_type string
// 	if mailtype == "html" {
// 		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
// 	} else {
// 		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
// 	}

// 	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
// 	send_to := strings.Split(to, ";")
// 	err := smtp.SendMail(host, auth, user, send_to, msg)
// 	return err
// }

// type TMailService struct {
// }

// func (this *TMailService) GetMailConfig() *TMailApiConfig {
// 	var config TMailApiConfig
// 	iConf := GetInternalConfByType(INTERNAL_CONFIG_MAILAPI, true)
// 	if iConf == nil || iConf.Config == "" {
// 		Log.Error("Get mail config error.")
// 		return nil
// 	}

// 	deConf, err := AesStrDecrypt(iConf.Config)
// 	if err != nil {
// 		Log.Error("mail config decrypt error.")
// 		return nil
// 	}

// 	err = UnserializeFromJson(deConf, &config)
// 	if err != nil {
// 		Log.Error("mail config UnserializeFromJson error.")
// 		return nil
// 	}

// 	return &config
// }

// func (this *TMailService) CallMailService(req *TUserApplyVerificationRecord) TErr {
// 	config := this.GetMailConfig()
// 	if config == nil {
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "无法获取MailApi配置")
// 	}

// 	var ec TErr
// 	switch config.ServCompany {
// 	case MAIL_SERV_COMPANY_GMAIL:
// 		ec = this.callGmailMailService(req, config)
// 	default:
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "不支持该邮件服务商")
// 	}
// 	return ec
// }

// func (this *TMailService) callGmailMailService(req *TUserApplyVerificationRecord, config *TMailApiConfig) TErr {

// 	err := SendMailForSmtp(config.UserName, config.Password, config.Url, req.Destination,
// 		req.Title, req.Detail, "html")

// 	req.ErrCode = "0"
// 	req.ErrMsg = "SUCCESS"
// 	req.SendMsgId = ""

// 	if err != nil {
// 		req.Status = APPLY_VERIFICATION_STATUS_FAIL
// 		Log.Error("Send mail msg error...:%v", err)
// 		req.ErrCode = "-1"
// 		req.ErrMsg = (err.Error())[0:120]
// 		return MakeErr(EC_SERVER_INTERNAL_ERROR, "邮件发送失败,请联系客服")
// 	} else {
// 		req.Status = APPLY_VERIFICATION_STATUS_USABLE
// 	}

// 	return nil

// }
