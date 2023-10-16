package sutils

import (
	"io"
	"net/http"
	"strings"

	. "github.com/switfs/shadow-framework/logger"
)

func allowCrossRequest(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                                                //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Set-Cookie, Content-Encoding, hex") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
}

func contentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json") //返回数据格式是json
}

func contentTypeHtml(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func SetHeaderAllowCrossRequest(w http.ResponseWriter) {
	allowCrossRequest(w)
	contentTypeJSON(w)
}

func NoList(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	Log.Info("NotFoundHandler Entering ... ...")

	if r.URL.Path == "/" {
		http.Redirect(w, r, "http://www.google.com", http.StatusFound)
		return
	}

	html := `Not Found 404`
	io.WriteString(w, html)
}
