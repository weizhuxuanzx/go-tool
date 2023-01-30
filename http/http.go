package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	httpurl "net/url"
)

type Json map[string]any

func (j Json) Read() io.Reader {
	s, _ := json.Marshal(j)
	return bytes.NewReader(s)
}

type Client struct {
	Client *http.Client
}

func NewHttpClient() *Client {
	return &Client{
		Client: &http.Client{},
	}
}

func (h Client) Post(url string, v any, data map[string]any) (any, error) {
	params := h.Reader(data)
	req, err := h.Client.Post(url, "application/json", params)
	if err != nil {
		return nil, err
	}
	if req.StatusCode != 200 {
		return nil, err
	}
	resultContent, err := io.ReadAll(req.Body)
	result, err := h.Unmarshal(resultContent, v)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h Client) Get(url string, v any, data map[string]string) (any, error) {
	params := httpurl.Values{}        //参数集合
	reqUrl, err := httpurl.Parse(url) //请求地址
	if err != nil {
		return nil, err
	}
	for key, val := range data {
		params.Set(key, val)
	}
	reqUrl.RawQuery = params.Encode()      //组合url
	resp, err := http.Get(reqUrl.String()) //发起get请求
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body) //解析请求信息
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, v) //转换为map
	if err != nil {
		return nil, err
	}
	return v, err
}

func (h Client) Reader(data map[string]any) io.Reader {
	sendJson, _ := json.Marshal(data)
	params := bytes.NewReader(sendJson)
	return params
}

func (h Client) Unmarshal(content []byte, v any) (any, error) {
	err := json.Unmarshal(content, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
