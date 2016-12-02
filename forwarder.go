package forwarder

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var Debug = false
var UseGZip = true

func Forward(w http.ResponseWriter, r *http.Request, targetURL string) (err error) {
	var buf io.Reader
	buf = r.Body
	if Debug {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		log.Println(`curl -d '` + string(b) + `' ` + targetURL)
		buf = bytes.NewBuffer(b)
	}

	request, err := http.NewRequest(r.Method, targetURL, buf)
	if err != nil {
		return
	}
	request.Header = r.Header

	httpClient := &http.Client{}
	res, err := httpClient.Do(request)
	if err != nil {
		return
	}

	defer res.Body.Close()

	rd := res.Body
	if UseGZip {
		str := r.Header.Get(`Accept-Encoding`)
		if strings.Contains(str, `gzip`) {
			rd, err = gzip.NewReader(res.Body)
			if err != nil { // 请求不要求gzip编码
				return
			}
		}
	}
	for k, _ := range res.Header {
		v := res.Header.Get(k)
		w.Header().Add(k, v)
	}

	_, err = io.Copy(w, rd)
	if err != nil {
		return
	}

	w.WriteHeader(res.StatusCode)
	return
}
