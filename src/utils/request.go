package utils

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Request[T any] struct {
	c       *resty.Client
	rawResp *resty.Response
}

func NewRequest[T any]() *Request[T] {
	return &Request[T]{
		c:       resty.New(),
		rawResp: nil,
	}
}
func (r *Request[T]) Get(url string, params map[string]string) (T, error) {
	var result T
	userAgent := JsonDB.Get("user_agent").(string)
	resp, err := r.c.R().
		SetHeader("User-Agent", userAgent).
		SetQueryParams(params).
		Get(url)
	if err != nil {
		return result, err
	}
	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("请求失败  %d", resp.StatusCode())
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return result, err
	}
	r.rawResp = resp
	// fmt.Println("[Debug]", string(resp.Body()))
	return result, nil
}
func (r *Request[T]) Post(url string, body any) (T, error) {
	var result T
	userAgent := JsonDB.Get("user_agent").(string)
	resp, err := r.c.R().
		SetHeader("User-Agent", userAgent).
		SetBody(body).
		Post(url)
	if err != nil {
		return result, err
	}
	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("请求失败  %d", resp.StatusCode())
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return result, err
	}
	r.rawResp = resp
	return result, nil
}
func (r *Request[T]) SetCookie(cookie string) *Request[T] {
	r.c.SetHeader("Cookie", cookie)
	return r
}
func (r *Request[T]) GetRawResp() *resty.Response {
	return r.rawResp
}
func (r *Request[T]) GetRawRespBody() []byte {
	return r.rawResp.Body()
}
