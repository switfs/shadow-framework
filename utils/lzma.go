package utils

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/switfs/shadow-framework/extra/lzma"
)

func LZMADecode(mstr string) (string, error) {

	m, err := hex.DecodeString(mstr)
	if err != nil {
		return "", err
	}

	b := new(bytes.Buffer)
	in := bytes.NewBuffer([]byte(m))
	r := lzma.NewReader(in)
	defer r.Close()
	b.Reset()
	_, err = io.Copy(b, r)
	if err != nil {
		return "", err
	}
	//    if err == nil { // if err != nil, there is little chance that data is decoded correctly, if at all
	s := b.String()
	return s, nil
	//    }
}

func LZMAEncode(m []byte) (string, error) {
	b := new(bytes.Buffer)
	pr, pw := io.Pipe()
	defer pr.Close()
	in := bytes.NewBuffer(m)
	size := int64(len(m))
	var oerr error = nil
	go func() {
		defer pw.Close()
		w := lzma.NewWriterSizeLevel(pw, size, 1)
		defer w.Close()
		_, err := io.Copy(w, in)
		if err != nil {
			oerr = err
		}
	}()
	b.Reset()
	if oerr != nil {
		return "", oerr
	}

	_, err := io.Copy(b, pr)
	if err != nil {
		return "", err
	}
	res := b.Bytes()

	compressedStr := hex.EncodeToString(res)
	return compressedStr, nil
}
