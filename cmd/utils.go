package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func constructUrl(base string, params url.Values) (uri string) {
	p := params.Encode()
	uri = base + "?" + p
	return
}

// doRequest doing request
func doRequest(method string, url string, m interface{}, headers map[string]string, reqBody []byte) (err error, code int, respBody []byte) {

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

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
