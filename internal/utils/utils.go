package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"musictags-joiner/pkgs/logger"
	"net/http"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigFastest
)

//SendResponse send response
func SendResponse(logger *logger.Logger, w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Errorf("couldn't send data to connection: %v", err)
	}
}

func ConstructUrl(base string, params url.Values) (uri string) {
	p := params.Encode()
	uri = base + "?" + p
	return
}

// doRequest doing request
func DoRequest(method string, url string, m interface{}, headers map[string]string, reqBody []byte) (err error, code int, respBody []byte) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("Can't get request; %v", err), 1, respBody
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err, 0, respBody
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Can't read response; %v", err), resp.StatusCode, respBody
	}
	defer resp.Body.Close()

	if m != nil && resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(respBody, m)
		if err != nil {
			return fmt.Errorf("Can't Unmarshal response; %v", err), resp.StatusCode, respBody
		}
	}
	return nil, resp.StatusCode, respBody
}

func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
