package sutils

// import (
// 	. "alpcserver"
// 	. "comutils/serverutils"
// )

// const (
// 	NEED_COMPRESS_LEN = 2048
// )

// func CommMsgDecrypt(key string, encryptedMsg string) string {

// 	b64Msg := FastDecrypt(key, encryptedMsg)
// 	b64str, err := Base64Decode([]byte(b64Msg))
// 	if err != nil {
// 		Log.Debug("decrypt base64 failed, b64Msg=%v", b64Msg)
// 		return ""
// 	}

// 	return string(b64str)
// }

// func CommMsgEncrypt(key string, plainMsg string) string {
// 	b64msg := Base64Encode([]byte(plainMsg))
// 	encMsg := FastEncrypt(key, string(b64msg))

// 	return encMsg
// }

// func CompressMsg(msg string) string {
// 	compressed, err := LzmaEncode([]byte(msg))
// 	if err != nil {
// 		Log.Warning("compress failed, using original msg, err = %v", err)
// 		return msg
// 	}

// 	return compressed
// }

// func DecompressMsg(compressedMsg string) string {

// 	msg, err := LzmaDecode((compressedMsg))
// 	if err != nil {
// 		Log.Debug("uncompressed msg failed , compressedMsg = %v", compressedMsg)
// 		return ""
// 	}

// 	return msg
// }

// func DecompressGameReq(req *TMsgReq, msg string) (*TMsgReq, string) {

// 	if req.Header.Compress == 1 {
// 		datastr := req.Data.(string)
// 		decomDataStr := DecompressMsg(datastr)

// 		var d interface{}
// 		err := UnserializeFromJson(decomDataStr, &d)
// 		if err != nil {
// 			Log.Warning("unserialize decompressDataStr failed, err = %v", err)
// 			return nil, ""
// 		}

// 		dreq := *req
// 		dreq.Data = d
// 		newmsg := SerializeToJson(dreq)

// 		return &dreq, newmsg
// 	} else {
// 		return req, msg
// 	}
// }

// func CompressGameResp(clntmsg string) string {

// 	enable := GetInt64Config(CFG_GS_ENABLE_COMPRESS)
// 	if enable != 1 {
// 		return clntmsg
// 	}

// 	if len(clntmsg) < NEED_COMPRESS_LEN || IsMetaData(clntmsg) {
// 		return clntmsg
// 	}

// 	var clntresp TMsgClientResp
// 	err := UnserializeFromJson(clntmsg, &clntresp)
// 	if err != nil {
// 		Log.Warning("parse json failed, err = %v", err)
// 		return clntmsg
// 	}

// 	// ä½¿ç”¨åŽ‹ç¼©
// 	clntresp.Header.Compress = 1
// 	data := clntresp.Data
// 	datastr := SerializeToJson(data)
// 	compressed, err := LzmaEncode([]byte(datastr))
// 	if err != nil {
// 		Log.Warning("compress failed, using original msg, err = %v", err)
// 		return clntmsg
// 	}

// 	newrsp := TMsgClientResp{
// 		Header: clntresp.Header,
// 		Data:   compressed,
// 	}

// 	newmsg := SerializeToJson(newrsp)

// 	return newmsg
// }
