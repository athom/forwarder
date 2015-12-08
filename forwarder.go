package forwarder

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func Forward(w http.ResponseWriter, r *http.Request, targetURL string) (err error) {
	request, err := http.NewRequest(r.Method, targetURL, r.Body)
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

	str := r.Header.Get(`Accept-Encoding`)
	rd := res.Body
	if strings.Contains(str, `gzip`) {
		rd, err = gzip.NewReader(res.Body)
		if err != nil { // 请求不要求gzip编码
			return
		}
	}
	_, err = io.Copy(w, rd)
	if err != nil {
		return
	}

	w.WriteHeader(res.StatusCode)
	return
}
