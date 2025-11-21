package util

import (
	"fmt"
	"io"
	"net/http"
	libURL "net/url"
	"strings"
	"time"
)

type ReqBody struct {
	// ContentType - describe the content-type of
	ContentType string
	Value       string
}

func BuildBaseHttpRequest(
	url string,
	method string,
	headers [][2]string, // [][key=xxx, value=xxx]
	queryParams [][2]string,
	body ReqBody,
) (*http.Request, error) {
	// #1. parse url
	urlObj, err := libURL.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("解析URL时出现异常 - %s", err.Error())
	}
	// #2. normalize method
	var methodMatch bool = false
	if method == "" {
		method = "GET"
	}

	for _, method := range []string{"get", "post", "put", "delete", "head", "options"} {
		if method == strings.ToLower(method) {
			methodMatch = true
			break
		}
	}
	if !methodMatch {
		return nil, fmt.Errorf("不支持的HTTP请求方法 - %s", method)
	}

	// #3. if queryParams content has params, then override original urlObj's rawQuery
	if len(queryParams) > 0 {
		qValues := libURL.Values{}
		for _, qpair := range queryParams {
			k := qpair[0]
			v := qpair[1]

			if qValues.Get(k) == "" {
				qValues.Set(k, v)
			} else {
				qValues.Add(k, v)
			}
		}
		urlObj.RawQuery = qValues.Encode()
	}

	finalHeaders := http.Header{}
	// #4. set body
	if body.ContentType == "" {
		// set default content-type header
		finalHeaders.Set("Content-Type", "text/plain")
	} else {
		finalHeaders.Set("Content-Type", body.ContentType)
	}
	bodyReadCloser := io.NopCloser(strings.NewReader(body.Value))

	// #5. set headers (may override content-type)
	for _, headerPair := range headers {
		k := headerPair[0]
		v := headerPair[1]
		if finalHeaders.Get(k) == "" {
			finalHeaders.Set(k, v)
		} else {
			finalHeaders.Add(k, v)
		}
	}

	return &http.Request{
		Method:     method,
		URL:        urlObj,
		Header:     finalHeaders,
		Body:       bodyReadCloser,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}, nil
}

func SendHttpRequest(req *http.Request, allowRedicrect bool, timeout int) (*http.Response, []byte, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !allowRedicrect {
				return fmt.Errorf("此请求不允许自动重定向")
			}
			return nil
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, []byte{}, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, err
	}

	return resp, content, nil
}
