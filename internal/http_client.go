/*
Copyright © 2022
Author Bhakiyaraj Kalimuthu
Email bhakiya.kalimuthu@gmail.com
*/

package internal

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type HttpClient interface {
	GenerateLoad(body interface{}) error
}
type httpClient struct {
	log    *zap.Logger
	client *http.Client
	url    string
}

func NewHttpClient(log *zap.Logger, url string) HttpClient {
	return &httpClient{
		log: log,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		url: url,
	}
}

func (h *httpClient) GenerateLoad(body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		h.log.Error("failed to marshal body", zap.Error(err))
		return err
	}
	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewBuffer(b))
	if err != nil {
		h.log.Error("failed to create request", zap.Error(err))
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", getIP())
	_, err = h.client.Do(req)
	if err != nil {
		h.log.Error("http request failed", zap.Error(err))
	}
	return err
}


