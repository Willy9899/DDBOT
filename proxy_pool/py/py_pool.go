package py

import (
	"errors"
	"fmt"
	"github.com/Sora233/Sora233-MiraiGo/proxy_pool"
	"github.com/asmcos/requests"
)

// https://github.com/jhao104/proxy_pool
type ProxyPool struct {
	Host string
}

type GetResponse struct {
	CheckCount int    `json:"check_count"`
	FailCount  int    `json:"fail_count"`
	LastStatus int    `json:"last_status"`
	LastTime   string `json:"last_time"`
	Proxy      string `json:"proxy"`
	Region     string `json:"region"`
	Source     string `json:"source"`
	Type       string `json:"type"`
}

type PyProxy struct {
	proxy string
}

func (p *PyProxy) ProxyString() string {
	return p.proxy
}

func (pool *ProxyPool) Get() (proxy_pool.IProxy, error) {
	if pool == nil {
		return nil, errors.New("<nil>")
	}
	resp, err := requests.Get(fmt.Sprintf("%v/get", pool.Host))
	if err != nil {
		return nil, err
	}
	getResp := new(GetResponse)
	err = resp.Json(getResp)
	if err != nil {
		return nil, err
	}
	if getResp.Proxy == "" {
		return nil, errors.New("empty proxy")
	}
	proxy := getResp.Proxy
	return &PyProxy{proxy}, nil
}

type DeleteResponse struct {
	Code int `json:"code"`
	Src  int `json:"src"`
}

func (pool *ProxyPool) Delete(proxy proxy_pool.IProxy) bool {
	if pool == nil {
		return false
	}
	resp, err := requests.Get(fmt.Sprintf("%v/delete?proxy=%v", pool.Host, proxy.ProxyString()))
	if err != nil {
		return false
	}
	deleteResp := new(DeleteResponse)
	err = resp.Json(deleteResp)
	if err != nil {
		return false
	}
	return deleteResp.Src == 1
}

func NewPYProxyPool(host string) (*ProxyPool, error) {
	pool := &ProxyPool{Host: host}
	resp, err := requests.Get(pool.Host)
	if err != nil {
		return nil, err
	}
	resp.Content()
	if resp.R.StatusCode != 200 {
		return nil, errors.New(resp.R.Status)
	}
	return pool, nil
}