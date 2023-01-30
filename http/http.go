package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

func (h Client) PostJson(url string, v any, data map[string]any) (any, error) {
	params := h.Reader(data)
	req, err := h.Client.Post(url, "application/json", params)
	if err != nil {
		return nil, err
	}
	resultContent, err := io.ReadAll(req.Body)
	result, err := h.Unmarshal(resultContent, v)
	if err != nil {
		return nil, err
	}
	return result, nil
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
