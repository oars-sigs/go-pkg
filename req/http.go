package req

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var DefaultHTTPClient = &http.Client{
	Timeout: 30 * time.Minute,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func ReqJSON(method, urlStr string, body interface{}, respBody interface{}, headers map[string]string) error {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	data, err := ReqBtye(method, urlStr, bodyData, headers)
	if err != nil {
		return err
	}
	if respBody == nil {
		return nil
	}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return err
	}
	return nil
}

func ReqBtye(method, urlStr string, body []byte, headers map[string]string) ([]byte, error) {
	resp, err := Req(method, urlStr, bytes.NewReader(body), headers)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Req(method, urlStr string, body io.Reader, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("code: %d", resp.StatusCode)
	}
	return resp.Body, nil
}
