package sutils

// import (
// 	. "alpcserver"
// 	. "comutils/serverutils"
// 	"errors"

// 	"io"

// 	"github.com/jlaffaye/ftp"
// )

// var FtpService TFtpService

// type TFtpService struct {
// }

// func (this *TFtpService) GetFtpConfig() *TFtpApiConfig {
// 	var config TFtpApiConfig
// 	iConf := GetInternalConfByType(INTERNAL_CONFIG_FTPAPI, true)
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

// func (this *TFtpService) ftpConnect(config *TFtpApiConfig) (*ftp.ServerConn, error) {
// 	ftp, err := ftp.Connect(config.Url)
// 	if err != nil {
// 		Log.Warning("addr:%v, err:%v", config.Url, err)
// 		return nil, errors.New("FTP文件服务器连线失败")
// 	}
// 	err = ftp.Login(config.UserName, config.Password)
// 	if err != nil {
// 		Log.Warning("user:%v, pw:%v, err:%v", config.UserName, config.Password, err)
// 		return nil, errors.New("FTP文件服务器登入失败")
// 	}
// 	return ftp, nil
// }

// func (this *TFtpService) ftpDisconnect(ftp *ftp.ServerConn) {
// 	ftp.Logout()
// 	ftp.Quit()
// }

// func (this *TFtpService) UploadFile(saveName string, file io.Reader) (string, *TErr {
// 	config := this.GetFtpConfig()
// 	if config == nil {
// 		return "", MakeErr(EC_SERVER_INTERNAL_ERROR, "无法获取FtpApi配置")
// 	}
// 	ftp, err := this.ftpConnect(config)
// 	if err != nil {
// 		return "", MakeErr(EC_SERVER_INTERNAL_ERROR, err)
// 	}
// 	defer this.ftpDisconnect(ftp)

// 	if len(config.RemotePath) > 0 {
// 		ftp.MakeDir(config.RemotePath)
// 	}
// 	ftp.ChangeDir(config.RemotePath)
// 	dir, _ := ftp.CurrentDir()
// 	Log.Info("ftp.CurrentDir:%v", dir)

// 	err = ftp.Stor(saveName, file)
// 	if err != nil {
// 		return "", MakeErr(EC_SERVER_INTERNAL_ERROR, "文档上传失败")
// 	}
// 	filePath := config.RemotePath + "/" + saveName
// 	return filePath, SUCCESS
// }
